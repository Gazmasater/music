package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"music.com/internal/models"
	"music.com/pkg/logger"
)

// GetInfoHandler godoc
// @Summary Get API Information
// @Description Returns general information about the API, including title and version.
// @Tags info
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /info [get]
func GetInfoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Debug(ctx, "Entering GetInfoHandler")

	info := map[string]interface{}{
		"title":   "Music info",
		"version": "0.0.1",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		logger.Error(ctx, "Failed to encode response", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Info(ctx, "API info requested")
}

func GetSongsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger.Info(ctx, "Handling GetSongs request...")

		// Получаем параметры запроса
		field := r.URL.Query().Get("field")
		value := r.URL.Query().Get("value")

		// Параметры пагинации
		limitStr := r.URL.Query().Get("limit")
		pageStr := r.URL.Query().Get("page")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10 // Дефолтное количество записей на страницу
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1 // Дефолтная страница
		}

		offset := (page - 1) * limit

		logger.DebugKV(ctx, "Filter parameters", "field", field, "value", value)
		logger.DebugKV(ctx, "Pagination", "limit", limit, "page", page, "offset", offset)

		// Подготавливаем запрос с фильтрацией и пагинацией
		var songs []models.SongDetail
		query := db.Model(&models.SongDetail{})

		if field != "" && value != "" {
			switch field {
			case "song_name":
				// Используем ILIKE для точного соответствия, игнорируя регистр
				query = query.Where("song_name ILIKE ?", value)
			case "artist_name":
				query = query.Joins("JOIN artists ON artists.id = song_details.artist_id").
					Where("artists.name ILIKE ?", "%"+value+"%")
			case "release_date":
				releaseDate, err := time.Parse("2006-01-02", value)
				if err == nil {
					query = query.Where("release_date = ?", releaseDate)
				}
			default:
				logger.Error(ctx, "Invalid filter field")
				http.Error(w, "Invalid filter field", http.StatusBadRequest)
				return
			}
		} else {
			logger.Debug(ctx, "No filtering parameters provided")
		}

		// Пагинация
		query = query.Limit(limit).Offset(offset)

		if err := query.Find(&songs).Error; err != nil {
			logger.Error(ctx, "Error fetching songs")
			http.Error(w, "Error fetching songs", http.StatusInternalServerError)
			return
		}

		logger.DebugKV(ctx, "Fetched songs count", "count", len(songs))

		// Формируем ответ
		response := struct {
			TotalItems int                 `json:"total_items"`
			Page       int                 `json:"page"`
			Limit      int                 `json:"limit"`
			Songs      []models.SongDetail `json:"songs"`
		}{
			TotalItems: len(songs),
			Page:       page,
			Limit:      limit,
			Songs:      songs,
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error(ctx, "Failed to encode response")
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		logger.Info(ctx, "Successfully handled GetSongs request")
	}
}

// AddSongHandler добавляет новую песню в базу данных.
// @Summary Добавить новую песню
// @Description Добавляет новую песню к исполнителю. Если исполнитель не существует, он будет создан.
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.SongInput true "Информация о песне"
// @Success 201 {object} models.SongDetail "Успешно добавлена новая песня"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 409 {string} string "Песня уже существует"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs [post]
func AddSongHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger.Debug(ctx, "Entering AddSongHandler")

		// Структура для получения базовой информации о песне
		var songInput struct {
			Group string `json:"group" binding:"required"` // Имя исполнителя
			Song  string `json:"song" binding:"required"`  // Название песни
		}

		// Декодируем запрос
		if err := json.NewDecoder(r.Body).Decode(&songInput); err != nil {
			logger.Error(ctx, "Failed to decode new song", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		logger.DebugKV(ctx, "Decoded song input", "song_input", songInput)

		// Проверка на существование исполнителя
		var artist models.Artist
		if err := db.Where("name = ?", songInput.Group).First(&artist).Error; err != nil {
			logger.DebugKV(ctx, "Artist not found, creating new artist", "artist_name", songInput.Group)
			// Если исполнитель не существует, создаем нового
			artist = models.Artist{Name: songInput.Group}
			if err := db.Create(&artist).Error; err != nil {
				logger.Error(ctx, "Failed to add new artist to database", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			logger.Info(ctx, "New artist created", artist)
		} else {
			logger.DebugKV(ctx, "Artist found", "artist_id", artist.ID)
		}

		// Проверка на существование песни с таким названием у данного исполнителя
		var existingSong models.SongDetail
		if err := db.Where("song_name = ? AND artist_id = ?", songInput.Song, artist.ID).First(&existingSong).Error; err == nil {
			// Песня уже существует
			logger.Error(ctx, "Song already exists", err)
			http.Error(w, "Song already exists", http.StatusConflict)
			return
		} else {
			logger.Debug(ctx, "No existing song found")
		}

		// Создаем новую песню с минимальной информацией (название и исполнитель)
		newSong := models.SongDetail{
			ArtistID:  artist.ID, // Приведение типа
			SongName:  songInput.Song,
			GroupName: songInput.Group,
		}

		logger.DebugKV(ctx, "Creating new song", "new_song", newSong)

		// Сохраняем новую песню в базе данных
		if err := db.Create(&newSong).Error; err != nil {
			logger.Error(ctx, "Failed to add new song to database", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Возвращаем статус 201 Created
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(newSong); err != nil {
			logger.Error(ctx, "Failed to encode new song response", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.Info(ctx, "New song added", newSong)
	}
}

// DeleteSongHandler возвращает обработчик HTTP, который удаляет песню из базы данных по её имени.
//
// Параметры:
//
//	db *gorm.DB: экземпляр базы данных GORM, используемый для выполнения операций с БД.
//
// Возвращает:
//
//	http.HandlerFunc: функция-обработчик, которая принимает ResponseWriter и запрос,
//	а затем выполняет логику удаления песни.
//
// В случае успешного удаления возвращает статус 204 No Content.
// Если песня не найдена, возвращает статус 404 Not Found.
// В случае ошибки при удалении возвращает статус 500 Internal Server Error.
//
// @Router /songs/{songName} [delete]
// @Param songName path string true "Имя песни для удаления"
// @Success 204 {object} nil "Успешное удаление"
// @Failure 404 {object} nil "Песня не найдена"
// @Failure 500 {object} nil "Ошибка при удалении песни"
func DeleteSongHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		songName := chi.URLParam(r, "songName")

		var song models.SongDetail
		if err := db.Where("song_name = ?", songName).First(&song).Error; err != nil {
			logger.Warn(ctx, "Attempt to delete non-existent song", "songName", songName)
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Удаляем песню
		if err := db.Delete(&song).Error; err != nil {
			logger.Error(ctx, "Failed to delete song from database", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Возвращаем 204 No Content
		w.WriteHeader(http.StatusNoContent)
	}
}

// @Router /songs/{songName} [put]
// @Param songName path string true "Имя песни для обновления"
// @Param body body models.SongUpdateResponse true "Обновленные данные песни. Все поля являются необязательными."
// @Success 200 {object} models.SongUpdateResponse "Успешное обновление песни"
// @Failure 400 {object} nil "Некорректный запрос"
// @Failure 404 {object} nil "Песня не найдена"
// @Failure 500 {object} nil "Ошибка при обновлении песни"
func UpdateSongHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger.Debug(ctx, "Entering UpdateSongHandler")

		songName := chi.URLParam(r, "songName")
		logger.Debug(ctx, "Song name from URL param", "songName", songName)

		var song models.SongDetail
		if err := db.Where("song_name = ?", songName).First(&song).Error; err != nil {
			logger.Warn(ctx, "Attempt to update non-existent song", "songName", songName, "error", err)
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		var updatedData models.SongUpdateResponse
		if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
			logger.Error(ctx, "Failed to decode updated song", "error", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Находим ID исполнителя по имени
		if updatedData.ArtistName != "" {
			var artist models.Artist
			if err := db.Where("name = ?", updatedData.ArtistName).First(&artist).Error; err != nil {
				logger.Error(ctx, "Artist not found", "artistName", updatedData.ArtistName)
				http.Error(w, "Artist Not Found", http.StatusNotFound)
				return
			}
			song.ArtistID = artist.ID
		}

		// Обновляем только необходимые поля
		if updatedData.SongName != "" {
			song.SongName = updatedData.SongName
		}
		if !updatedData.ReleaseDate.IsZero() {
			song.ReleaseDate = updatedData.ReleaseDate
		}
		if updatedData.GroupLink != "" {
			song.GroupName = updatedData.GroupLink
		}
		if len(updatedData.Text) != 0 {
			song.Text = string(updatedData.Text)
		}

		if err := db.Save(&song).Error; err != nil {
			logger.Error(ctx, "Failed to update song in database", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Возвращаем обновленные данные
		// Возвращаем обновленные данные
		response := models.SongUpdateResponse{
			ArtistName:  updatedData.ArtistName,
			SongName:    song.SongName,
			ReleaseDate: song.ReleaseDate,
			GroupLink:   song.GroupName,
			Text:        song.Text,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error(ctx, "Failed to encode updated song response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.Info(ctx, "Song updated successfully", "updatedSong", song)
	}
}

// GetSongLyricsHandler получает текст песни с поддержкой пагинации.
// @Router /songs/{songName}/lyrics [get]
// @Param songName path string true "Имя песни для получения текста"
// @Param verse_page query int false "Номер страницы куплетов" default(1)
// @Param verse_limit query int false "Количество куплетов на странице" default(3)
// @Success 200 {object} models.PaginatedLyricsRespons "Успешное получение текста песни"
// @Failure 400 {object} nil "Некорректный запрос"
// @Failure 404 {object} nil "Песня не найдена"
// @Failure 500 {object} nil "Ошибка при получении текста песни"
func GetSongLyricsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Извлекаем название песни из параметров маршрута
		songName := chi.URLParam(r, "songName")
		logger.Debug(ctx, "Request to get song lyrics", "songName", songName)

		// Параметры пагинации (страница и лимит куплетов на странице)
		versePageStr := r.URL.Query().Get("verse_page")
		verseLimitStr := r.URL.Query().Get("verse_limit")

		versePage, err := strconv.Atoi(versePageStr)
		if err != nil || versePage < 1 {
			logger.Debug(ctx, "Invalid verse page parameter, defaulting to 1", "versePageStr", versePageStr)
			versePage = 1 // Установим дефолтное значение страницы
		}

		verseLimit, err := strconv.Atoi(verseLimitStr)
		if err != nil || verseLimit < 1 {
			logger.Debug(ctx, "Invalid verse limit parameter, defaulting to 3", "verseLimitStr", verseLimitStr)
			verseLimit = 3 // Установим дефолтное количество куплетов на странице
		}

		logger.Debug(ctx, "Pagination params", "versePage", versePage, "verseLimit", verseLimit)

		var song models.SongDetail
		if err := db.Where("song_name = ?", songName).First(&song).Error; err != nil {
			logger.Warn(ctx, "Song not found", "songName", songName)
			http.Error(w, "Song not found", http.StatusNotFound)
			return
		}

		var lyrics models.SongText
		if err := json.Unmarshal([]byte(song.Text), &lyrics); err != nil {
			logger.Error(ctx, "Failed to parse song lyrics", "songName", songName, "error", err)
			http.Error(w, "Failed to parse song lyrics", http.StatusInternalServerError)
			return
		}

		logger.Debug(ctx, "Lyrics retrieved", "totalVerses", len(lyrics.Verses))

		// Пагинация по куплетам
		totalVerses := len(lyrics.Verses)
		start := (versePage - 1) * verseLimit
		end := start + verseLimit

		if start >= totalVerses {
			logger.Warn(ctx, "Page out of range", "versePage", versePage, "totalVerses", totalVerses)
			http.Error(w, "Page out of range", http.StatusBadRequest)
			return
		}

		if end > totalVerses {
			end = totalVerses
		}

		paginatedVerses := lyrics.Verses[start:end]
		logger.Debug(ctx, "Paginated verses", "start", start, "end", end)

		response := models.PaginatedLyricsRespons{
			SongName:    song.SongName,
			VersePage:   versePage,
			VerseLimit:  verseLimit,
			TotalVerses: totalVerses,
			Verses:      paginatedVerses,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error(ctx, "Failed to encode response", "error", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		logger.Info(ctx, "Song lyrics retrieved successfully", "songName", songName)
	}
}
