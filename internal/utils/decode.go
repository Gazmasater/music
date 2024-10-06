package utils

import (
	"context"
	"encoding/json"
	"music/pkg/logger"
	"net/http"
)

func DecodeInput[T any](w http.ResponseWriter, r *http.Request, ctx context.Context, input *T, logMsg string) error {
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		logger.ErrorKV(ctx, "Failed to decode input", err)
		http.Error(w, "Bad Request: Invalid data format", http.StatusBadRequest)
		return err
	}
	logger.DebugKV(ctx, logMsg, "input", *input)
	return nil
}
