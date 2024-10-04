package models_test

import (
	"testing"

	"music/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestSongInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   models.SongInput
		wantErr string
	}{
		{
			name: "Valid input",
			input: models.SongInput{
				Group: "Valid Artist",
				Song:  "Valid Song",
			},
			wantErr: "", // Ошибки не ожидается
		},
		{
			name: "Empty artist name",
			input: models.SongInput{
				Group: "",
				Song:  "Valid Song",
			},
			wantErr: "artist name cannot be empty", // Ожидаемая ошибка
		},
		{
			name: "Empty song name",
			input: models.SongInput{
				Group: "Valid Artist",
				Song:  "",
			},
			wantErr: "song name cannot be empty", // Ожидаемая ошибка
		},
		{
			name: "Empty artist and song name",
			input: models.SongInput{
				Group: "",
				Song:  "",
			},
			wantErr: "artist name cannot be empty", // Ожидаемая ошибка
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Вызываем метод Validate
			err := tt.input.Validate()

			// Если ожидается ошибка, проверяем её текст
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			} else {
				// Если ошибки не ожидается, проверяем, что err == nil
				assert.NoError(t, err)
			}
		})
	}
}
