package utils_test

import (
	"context"
	"music/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeURLParameter(t *testing.T) {
	tests := []struct {
		name          string
		param         string
		expectedValue string
		expectedOk    bool
	}{
		{
			name:          "Valid parameter",
			param:         "hello%20world",
			expectedValue: "hello world",
			expectedOk:    true,
		},
		{
			name:          "Invalid parameter",
			param:         "%ZZ", // некорректный URL
			expectedValue: "",
			expectedOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый контекст и записываю в него результат работы http.ResponseWriter
			w := httptest.NewRecorder()
			ctx := context.Background()

			// Вызываем функцию DecodeURLParameter
			value, ok := utils.DecodeURLParameter(ctx, tt.param, w, "Decode failed")

			// Проверяем ожидаемые результаты
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedOk, ok)

			// Если параметр некорректный, проверяем, что ответ содержит ошибку
			if !tt.expectedOk {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})
	}
}
