package models

import (
	"errors"
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
	ID          uint `gorm:"primaryKey"`
	ArtistID    uint
	GroupName   string
	SongName    string
	ReleaseDate time.Time
	Text        string
	SongURL     string    `gorm:"column:song_url"` // Убедитесь, что это поле присутствует
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

// Song представляет минимальную информацию о песне для создания
type Song struct {
	ID        uint      `json:"id" gorm:"primaryKey"`             // Уникальный идентификатор песни
	ArtistID  uint      `json:"artist_id" gorm:"not null"`        // Ссылка на ID исполнителя
	SongName  string    `json:"song_name" gorm:"not null"`        // Название песни
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // Дата создания записи
}

type SongInput struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

// SongUpdateResponse представляет структуру ответа на обновление песни
type SongUpdateResponse struct {
	ArtistName  string    `json:"artist_name,omitempty"`
	SongName    string    `json:"song_name,omitempty"`
	ReleaseDate time.Time `json:"release_date,omitempty"`
	GroupLink   string    `json:"group_link,omitempty"`
	Text        SongText  `json:"text,omitempty"` // Измените на SongText
}

type PaginatedLyricsRespons struct {
	SongName    string   `json:"song_name"`    // Название песни
	VersePage   int      `json:"verse_page"`   // Номер страницы куплетов
	VerseLimit  int      `json:"verse_limit"`  // Количество куплетов на странице
	TotalVerses int      `json:"total_verses"` // Общее количество куплетов
	Verses      []string `json:"verses"`       // Пагинированные куплеты
}

type SongsResponse struct {
	TotalItems int          `json:"total_items"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	Songs      []SongDetail `json:"songs"`
}

// Validate проверяет, что поля в SongInput не пустые.
func (si *SongInput) Validate() error {
	if si.Group == "" {
		return errors.New("artist name cannot be empty")
	}
	if si.Song == "" {
		return errors.New("song name cannot be empty")
	}
	return nil
}
