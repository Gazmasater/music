package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"music.com/internal/handlers"
)

func NewRouter(db *gorm.DB) http.Handler {
	r := chi.NewRouter()

	r.Get("/info", handlers.GetInfoHandler)
	r.Get("/songs", handlers.GetSongsHandler(db))
	r.Post("/songs", handlers.AddSongHandler(db))
	r.Delete("/songs/{songName}", handlers.DeleteSongHandler(db))
	r.Put("/songs/{songName}", handlers.UpdateSongHandler(db))
	r.Get("/songs/{songName}/lyrics", handlers.GetSongLyricsHandler(db))

	return r
}
