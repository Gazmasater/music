package utils

import (
	"context"
	"music/pkg/logger"
	"net/http"
	"net/url"
)

func DecodeURLParameter(ctx context.Context, param string, w http.ResponseWriter, errorMessage string) (string, bool) {
	decodedParam, err := url.QueryUnescape(param)
	if err != nil {
		// Логируем ошибку с использованием logger
		logger.ErrorKV(ctx, errorMessage, "param", param, "error", err)
		http.Error(w, "Bad Request: "+errorMessage, http.StatusBadRequest)
		return "", false
	}
	return decodedParam, true
}
