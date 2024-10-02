package config

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"music.com/pkg/logger"
)

const (
	defaultPort         = "8080"
	defaultReadTimeout  = 10
	defaultWriteTimeout = 10
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		logger.Warn(context.Background(), "No .env file found", err)
	}
}

// getDurationFromEnv получает значение из переменной окружения и возвращает его как Duration
func getDurationFromEnv(envKey string, defaultValue int) time.Duration {
	envValue := os.Getenv(envKey)
	if envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			return time.Duration(value) * time.Second
		}
	}
	return time.Duration(defaultValue) * time.Second
}

func GetServerConfig() (string, time.Duration, time.Duration) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	readTimeout := getDurationFromEnv("READ_TIMEOUT", defaultReadTimeout)
	writeTimeout := getDurationFromEnv("WRITE_TIMEOUT", defaultWriteTimeout)

	return port, readTimeout, writeTimeout
}

// SetLogLevel устанавливает уровень логирования на основе переменной окружения
func SetLogLevel() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "debug"
	}

	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.DebugLevel
	}

	logger.SetLogger(logger.New(zap.NewAtomicLevelAt(level)))
}
