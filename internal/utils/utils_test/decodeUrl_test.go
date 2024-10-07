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
			name:          "Valid parameter with space",
			param:         "hello%20world",
			expectedValue: "hello world",
			expectedOk:    true,
		},
		{
			name:          "Valid parameter with plus",
			param:         "hello+world",
			expectedValue: "hello world",
			expectedOk:    true,
		},
		{
			name:          "Valid parameter with special characters",
			param:         "test%40example.com",
			expectedValue: "test@example.com",
			expectedOk:    true,
		},
		{
			name:          "Invalid parameter with incorrect encoding",
			param:         "%ZZ",
			expectedValue: "",
			expectedOk:    false,
		},
		{
			name:          "Invalid parameter with single %",
			param:         "%",
			expectedValue: "",
			expectedOk:    false,
		},
		{
			name:          "Invalid parameter with incomplete escape",
			param:         "hello%",
			expectedValue: "",
			expectedOk:    false,
		},
		{
			name:          "Valid parameter with Cyrillic",
			param:         "параметр%20с%20кириллицей",
			expectedValue: "параметр с кириллицей",
			expectedOk:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx := context.Background()

			value, ok := utils.DecodeURLParameter(ctx, tt.param, w, "Decode failed")

			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedOk, ok)

			if !tt.expectedOk {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})
	}
}
