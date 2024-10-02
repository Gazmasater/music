package models

import (
	"encoding/json"
	"time"
)

// Artist представляет исполнителя
type Artist struct {
	ID        uint      `json:"id" gorm:"primaryKey"`             // Уникальный идентификатор исполнителя
	Name      string    `json:"name" gorm:"unique"`               // Имя исполнителя
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // Дата создания записи
}

type SongText struct {
	Verses []string `json:"verses"` // Срез для хранения куплетов
}

type SongDetail struct {
	ID          int             `json:"id"`
	ArtistID    uint            `json:"artist_id"`
	GroupName   string          `json:"group_name"`
	SongName    string          `json:"song_name"`
	ReleaseDate time.Time       `json:"release_date"`
	Text        json.RawMessage `json:"text"` // Используйте RawMessage для хранения JSON
	CreatedAt   time.Time       `json:"created_at"`
}

// Song представляет минимальную информацию о песне для создания
type Song struct {
	ID        uint      `json:"id" gorm:"primaryKey"`             // Уникальный идентификатор песни
	ArtistID  uint      `json:"artist_id" gorm:"not null"`        // Ссылка на ID исполнителя
	SongName  string    `json:"song_name" gorm:"not null"`        // Название песни
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // Дата создания записи
}

type PaginatedLyricsResponse struct {
	SongName    string   `json:"song_name"`
	VersePage   int      `json:"verse_page"`
	VerseLimit  int      `json:"verse_limit"`
	TotalVerses int      `json:"total_verses"`
	Verses      []string `json:"verses"`
}
