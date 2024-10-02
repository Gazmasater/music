package models

import (
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
	ID          int       `json:"id"`
	ArtistID    uint      `json:"artist_id"`
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date"`
	Text        string    `json:"text"` // Заменяем json.RawMessage на string
	CreatedAt   time.Time `json:"created_at"`
}

// Song представляет минимальную информацию о песне для создания
type Song struct {
	ID        uint      `json:"id" gorm:"primaryKey"`             // Уникальный идентификатор песни
	ArtistID  uint      `json:"artist_id" gorm:"not null"`        // Ссылка на ID исполнителя
	SongName  string    `json:"song_name" gorm:"not null"`        // Название песни
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // Дата создания записи
}

type SongInput struct {
	Group string `json:"group" binding:"required"` // Имя исполнителя
	Song  string `json:"song" binding:"required"`  // Название песни
}

// SongUpdateResponse представляет структуру ответа на обновление песни
type SongUpdateResponse struct {
	ArtistName  string    `json:"artist_name,omitempty"`
	SongName    string    `json:"song_name,omitempty"`
	ReleaseDate time.Time `json:"release_date,omitempty"`
	GroupLink   string    `json:"group_link,omitempty"`
	Text        string    `json:"text,omitempty"`
}

type PaginatedLyricsRespons struct {
	SongName    string   `json:"song_name"`    // Название песни
	VersePage   int      `json:"verse_page"`   // Номер страницы куплетов
	VerseLimit  int      `json:"verse_limit"`  // Количество куплетов на странице
	TotalVerses int      `json:"total_verses"` // Общее количество куплетов
	Verses      []string `json:"verses"`       // Пагинированные куплеты
}
