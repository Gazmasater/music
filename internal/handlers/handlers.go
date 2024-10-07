package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"music/internal/models"
	"music/internal/utils"
	"music/pkg/logger"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
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

// GetSongsHandler возвращает обработчик HTTP для получения списка песен с поддержкой фильтрации и пагинации.
// @Summary Получить список песен
// @Description Получение списка песен с поддержкой фильтрации и пагинации.
// @Tags songs
// @Param field query string false "Поле для фильтрации (song_name, artist_name, release_date)"
// @Param value query string false "Значение для фильтрации"
// @Param limit query int false "Количество записей на странице"
// @Param page query int false "Номер страницы"
// @Success 200 {object} models.SongsResponse "Успешное получение списка песен"
// @Failure 400 {object} nil "Неверное поле для фильтрации"
// @Failure 500 {object} nil "Ошибка на сервере"
// @Router /songs [get]
func GetSongsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger.Info(ctx, "Handling GetSongs request...")

		// Получаем параметры запроса
		field := r.URL.Query().Get("field")
		value := r.URL.Query().Get("value")

		// Нормализуем поля фильтрации
		normalizedField := utils.NormalizeSongName(field)
		normalizedValue := utils.NormalizeSongName(value)

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

		logger.DebugKV(ctx, "Filter parameters", "field", normalizedField, "value", normalizedValue)
		logger.DebugKV(ctx, "Pagination", "limit", limit, "page", page, "offset", offset)

		// Подготавливаем запрос с фильтрацией и пагинацией
		var songs []models.SongDetail
		query := db.Model(&models.SongDetail{})

		if normalizedField != "" && normalizedValue != "" {
			switch normalizedField {
			case "song_name":
				// Используем ILIKE для точного соответствия, игнорируя регистр
				query = query.Where("song_name ILIKE ?", normalizedValue)
			case "artist_name":
				query = query.Joins("JOIN artists ON artists.id = song_details.artist_id").
					Where("artists.name ILIKE ?", "%"+normalizedValue+"%")
			case "release_date":
				releaseDate, err := time.Parse("2006-01-02", normalizedValue)
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
		response := models.SongsResponse{
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
		var songInput models.SongInput

		// Декодируем запрос с использованием отдельной функции
		if err := utils.DecodeInput(w, r, r.Context(), &songInput, "Decoded song input"); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Нормализация названия песни
		songInput.Song = utils.NormalizeSongName(songInput.Song)
		logger.DebugKV(ctx, "Normalized song input", "song_input", songInput)

		// Нормализация имени исполнителя
		songInput.Group = utils.NormalizeSongName(songInput.Group) // Здесь используем ту же функцию
		logger.DebugKV(ctx, "Normalized artist name", "artist_name", songInput.Group)

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

		// Возвращаем статус 200 Created
		w.WriteHeader(http.StatusOK) // Измените статус на 200 OK
		if err := json.NewEncoder(w).Encode(newSong); err != nil {
			logger.Error(ctx, "Failed to encode new song response", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.Info(ctx, "New song added", newSong)
	}
}

// DeleteSongHandler возвращает обработчик HTTP, который удаляет песню из базы данных по её имени.
// @Summary Удалить песню
// @Router /songs/{songName} [delete]
// @Param songName path string true "Имя песни для удаления"
// @Success 204 {object} nil "Успешное удаление"
// @Failure 404 {object} nil "Песня не найдена"
// @Failure 500 {object} nil "Ошибка при удалении песни"
func DeleteSongHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		songName := chi.URLParam(r, "songName")

		// Используем функцию DecodeURLParameter для декодирования
		decodedSongName, ok := utils.DecodeURLParameter(ctx, songName, w, "Invalid song name")
		if !ok {
			return
		}

		// Нормализуем имя песни с помощью utils
		normalizedSongName := utils.NormalizeSongName(decodedSongName)
		logger.Debug(ctx, "Нормализуем", "normalizedSongName", normalizedSongName)

		var song models.SongDetail
		if err := db.Where("song_name = ?", normalizedSongName).First(&song).Error; err != nil {
			logger.Warn(ctx, "Attempt to delete non-existent song", "songName", normalizedSongName)
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
// @Summary Изменение данных песни
// @Param songName path string true "Имя песни для обновления"
// @Param body body models.SongUpdateResponse true "Обновленные данные песни. Все поля являются необязательными."
// @Success 200 {object} models.SongUpdateResponse "Успешное обновление песни"
// @Failure 400 {object} nil "Некорректный запрос"
// @Failure 404 {object} nil "Песня не найдена"
// @Failure 500 {object} nil "Ошибка при обновлении песни"
// @Description Обновляет данные существующей песни по имени. Поля, которые не переданы, останутся без изменений.
func UpdateSongHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger.Debug(ctx, "Entering UpdateSongHandler")

		// Получаем название песни из URL и декодируем его с помощью DecodeURLParameter
		songName := chi.URLParam(r, "songName")
		decodedSongName, ok := utils.DecodeURLParameter(ctx, songName, w, "Invalid song name")
		if !ok {
			return
		}

		// Нормализуем название песни через utils
		normalizedSongName := utils.NormalizeSongName(decodedSongName)
		logger.Debug(ctx, "Normalized song name from URL param", "songName", normalizedSongName)

		// Проверяем, существует ли песня
		var song models.SongDetail
		if err := db.Where("song_name = ?", normalizedSongName).First(&song).Error; err != nil {
			logger.Warn(ctx, "Attempt to update non-existent song", "songName", normalizedSongName, "error", err)
			http.Error(w, "Song Not Found", http.StatusNotFound)
			return
		}

		// Получаем данные для обновления
		var updatedData models.SongUpdateResponse

		if err := utils.DecodeInput(w, r, ctx, &updatedData, "Decoded updated data"); err != nil {
			return
		}

		// Проверка на наличие полей для обновления
		if updatedData.SongName == "" && updatedData.ArtistName == "" && updatedData.GroupLink == "" && len(updatedData.Text.Verses) == 0 && updatedData.ReleaseDate == "" {
			logger.Warn(ctx, "No fields to update")
			http.Error(w, "Bad Request: No fields to update", http.StatusBadRequest)
			return
		}

		// Парсим строку даты в time.Time
		if updatedData.ReleaseDate != "" {
			parsedDate, err := time.Parse("2006.01.02", updatedData.ReleaseDate)
			if err != nil {
				logger.Error(ctx, "Failed to parse release date", "error", err)
				http.Error(w, "Bad Request: Invalid release date format", http.StatusBadRequest)
				return
			}
			song.ReleaseDate = parsedDate
			logger.Debug(ctx, "Release date updated", "newReleaseDate", song.ReleaseDate)
		}

		// Обновление информации о исполнителе
		if updatedData.ArtistName != "" {
			normalizedArtistName := utils.NormalizeSongName(updatedData.ArtistName) // Нормализуем имя исполнителя
			var artist models.Artist
			if err := db.Where("name = ?", normalizedArtistName).First(&artist).Error; err != nil {
				// Если исполнитель не найден, возвращаем ошибку
				logger.Error(ctx, "Artist not found", "artistName", normalizedArtistName)
				http.Error(w, "Artist Not Found", http.StatusNotFound)
				return
			}
			song.ArtistID = artist.ID
			logger.Debug(ctx, "Artist ID updated", "artistID", artist.ID)
		}

		// Обновление полей песни
		if updatedData.SongName != "" {
			normalizedNewSongName := utils.NormalizeSongName(updatedData.SongName)
			song.SongName = normalizedNewSongName
			logger.Debug(ctx, "Song name updated", "newSongName", normalizedNewSongName)
		}

		if updatedData.GroupLink != "" {
			normalizedGroupLink := utils.NormalizeSongName(updatedData.GroupLink)
			song.SongURL = normalizedGroupLink
			logger.Debug(ctx, "Group link updated", "newGroupLink", normalizedGroupLink)
		}

		if len(updatedData.Text.Verses) > 0 {
			textJSON, err := json.Marshal(updatedData.Text)
			if err != nil {
				logger.Error(ctx, "Failed to marshal updated text", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			song.Text = string(textJSON)
			logger.Debug(ctx, "Song text updated", "newText", updatedData.Text)
		}

		// Логируем текущее состояние песни перед сохранением
		logger.Debug(ctx, "Saving song", "song", song)

		// Сохранение обновленной песни в базу данных
		if err := db.Save(&song).Error; err != nil {
			logger.Error(ctx, "Failed to update song in database", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Формирование ответа с обновленными данными
		response := models.SongUpdateResponse{
			ArtistName:  updatedData.ArtistName,
			SongName:    song.SongName,
			ReleaseDate: song.ReleaseDate.Format("2006.01.02"), // Форматируем обратно в строку
			GroupLink:   song.SongURL,
			Text:        updatedData.Text,
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
// @Summary Получение текста песни с пагинацией по куплетам
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

		// Извлекаем и декодируем название песни из параметров маршрута
		songName := chi.URLParam(r, "songName")
		decodedSongName, ok := utils.DecodeURLParameter(ctx, songName, w, "Invalid song name")
		if !ok {
			return
		}

		// Нормализуем название песни
		normalizedSongName := utils.NormalizeSongName(decodedSongName)
		logger.Debug(ctx, "Request to get song lyrics", "songName", normalizedSongName)

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

		// Поиск песни в базе данных
		var song models.SongDetail
		if err := db.Where("song_name = ?", normalizedSongName).First(&song).Error; err != nil {
			logger.Warn(ctx, "Song not found", "songName", normalizedSongName)
			http.Error(w, "Song not found", http.StatusNotFound)
			return
		}

		// Проверяем текст песни, предполагая, что он уже разделен на куплеты
		logger.Debug(ctx, "Raw song text", "rawText", fmt.Sprintf("%q", song.Text))

		// Преобразуем текст в структуру SongText
		var lyrics models.SongText
		if err := json.Unmarshal([]byte(song.Text), &lyrics); err != nil {
			logger.Error(ctx, "Failed to unmarshal song text", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

		// Формируем ответ с куплетами и пагинацией
		response := models.PaginatedLyricsRespons{
			SongName:    song.SongName,
			VersePage:   versePage,
			VerseLimit:  verseLimit,
			TotalVerses: totalVerses,
			Verses:      paginatedVerses,
		}

		// Отправляем ответ в формате JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error(ctx, "Failed to encode response", "error", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		logger.Info(ctx, "Song lyrics retrieved successfully", "songName", normalizedSongName)
	}
}
