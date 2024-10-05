package handlers_test

import (
	"errors"
	"music/internal/handlers"
	"music/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDeleteSongHandler(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to the database")

	err = db.AutoMigrate(&models.Artist{}, &models.SongDetail{})
	require.NoError(t, err, "Failed to migrate database")

	artist := models.Artist{Name: "Test Artist"}
	require.NoError(t, db.Create(&artist).Error, "Failed to create test artist")

	song := models.SongDetail{
		ArtistID:  artist.ID,
		SongName:  "Test Song",
		GroupName: "Test Artist",
	}
	require.NoError(t, db.Create(&song).Error, "Failed to create test song")

	router := chi.NewRouter()
	router.Delete("/songs/{songName}", handlers.DeleteSongHandler(db))

	req, err := http.NewRequest(http.MethodDelete, "/songs/Test%20Song", nil)
	require.NoError(t, err, "Failed to create request")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	var deletedSong models.SongDetail
	err = db.Where("song_name = ?", "Test Song").First(&deletedSong).Error
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Expected song to be deleted from the database")

	var existingArtist models.Artist
	err = db.Where("name = ?", "Test Artist").First(&existingArtist).Error
	require.NoError(t, err, "Expected artist to still exist in the database")
	assert.Equal(t, "Test Artist", existingArtist.Name)

	req, err = http.NewRequest(http.MethodDelete, "/songs/NonExistentSong", nil)
	require.NoError(t, err, "Failed to create request for non-existent song")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
