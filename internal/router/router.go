package router

import (
	"net/http"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger" // Импортируйте для Swagger
	"gorm.io/gorm"
	_ "music.com/docs" // Импортируйте сгенерированные файлы Swagger
	"music.com/internal/handlers"
)

func NewRouter(db *gorm.DB) http.Handler {
	r := chi.NewRouter()

	// Роуты для API
	r.Get("/info", handlers.GetInfoHandler)
	r.Get("/songs", handlers.GetSongsHandler(db))
	r.Post("/songs", handlers.AddSongHandler(db))
	r.Delete("/songs/{songName}", handlers.DeleteSongHandler(db))
	r.Put("/songs/{songName}", handlers.UpdateSongHandler(db))
	r.Get("/songs/{songName}/lyrics", handlers.GetSongLyricsHandler(db))

	// Роут для Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler) // Доступ к Swagger документации

	return r
}
