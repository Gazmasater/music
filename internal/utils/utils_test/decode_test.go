package utils_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"music/internal/utils"

	"github.com/stretchr/testify/assert"
)

type TestInput struct {
	Name  string `json:"name"`
	Name1 string `json:"name1"`
}

func TestDecodeInput(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expected     TestInput
		expectError  bool
		expectedCode int
	}{
		{
			name:         "Valid JSON data",
			input:        `{"name": "John", "name1":"Put"}`,
			expected:     TestInput{Name: "John", Name1: "Put"},
			expectError:  false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Valid JSON data",
			input:        `{"name": "John numbe1", "name1":"Put number2"}`,
			expected:     TestInput{Name: "John numbe1", Name1: "Put number2"},
			expectError:  false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "inValid JSON data",
			input:        `{"name": "John numbe1", "name1": 30}`,
			expected:     TestInput{Name: "John numbe1", Name1: "Put number2"},
			expectError:  false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Преобразуем строку в буфер для использования в HTTP-запросе
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer([]byte(tt.input)))
			req.Header.Set("Content-Type", "application/json")

			// Создание ResponseRecorder для фиксации HTTP-ответа
			w := httptest.NewRecorder()

			// Создание контекста
			ctx := context.Background()

			// Переменная для декодирования
			var decodedData TestInput

			// Сохранение начального состояния структуры
			initialData := decodedData

			// Вызов функции DecodeInput
			err := utils.DecodeInput(w, req, ctx, &decodedData, "Decoded test input")

			// Проверка результата
			if tt.expectError {
				assert.Error(t, err)
				// При ошибке декодирования проверяем, что данные остались на месте
				assert.Equal(t, initialData, decodedData)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, decodedData)
			}

			// Проверка статуса ответа
			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)
		})
	}
}
