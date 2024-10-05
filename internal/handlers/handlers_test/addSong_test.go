package handlers_test

import (
	"bytes"
	"context"
	"music/internal/models"
	"music/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"music/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAddSongHandler(t *testing.T) {
	// Создаем тестовую базу данных в памяти
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to the database")

	// Миграция для создания таблиц
	err = db.AutoMigrate(&models.Artist{}, &models.SongDetail{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	logger.Info(context.Background(), "Database migrated successfully!")

	// Создаем роутер и добавляем обработчик
	router := chi.NewRouter()
	router.Post("/songs", handlers.AddSongHandler(db))

	// Создаем тестовый JSON для успешного декодирования
	validJSON := []byte(`{
		"song": "Test Song",
		"group": "Test Artist"
	}`)

	// Создаем тестовый запрос
	req, err := http.NewRequest(http.MethodPost, "/songs", bytes.NewBuffer(validJSON))
	require.NoError(t, err, "Failed to create request")

	// Запускаем тестовый сервер
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Проверяем, что песня была добавлена в базу данных
	var addedSong models.SongDetail
	err = db.Where("song_name = ?", "Test Song").First(&addedSong).Error
	require.NoError(t, err, "Expected song to be added to the database")

	// Проверяем, что добавленная песня имеет правильные данные
	assert.Equal(t, "Test Song", addedSong.SongName)
	assert.Equal(t, "Test Artist", addedSong.GroupName)

	// Проверяем, что артист был добавлен
	var artist models.Artist
	err = db.Where("name = ?", "Test Artist").First(&artist).Error
	require.NoError(t, err, "Expected artist to be added to the database")
	assert.Equal(t, "Test Artist", artist.Name)
}
