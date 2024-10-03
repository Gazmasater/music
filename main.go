package main

import (
	"context"
	"fmt"
	"net/http"

	"music.com/config"
	"music.com/internal/db"

	"music.com/internal/router"
	"music.com/pkg/logger"
)

// @title Music API
// @version 1.0
// @description Это API для работы с музыкальной библиотекой, позволяющее получать, добавлять, обновлять и удалять песни.
func main() {
	// Загружаем переменные окружения
	config.LoadEnv()

	// Получаем конфигурацию сервера
	port, readTimeout, writeTimeout := config.GetServerConfig()
	config.SetLogLevel()

	// Создание нового контекста с логгером
	ctx := context.Background()
	ctx = logger.ToContext(ctx, logger.Global())

	// Подключение к базе данных
	database, err := db.Connect()
	if err != nil {
		logger.Fatal(ctx, "failed to connect to the database", err) // Используем ваш логгер
	}

	logger.Info(ctx, "Database connection established successfully!") // Логируем успешное подключение

	db.Migrate(database)

	// Передаем соединение базы данных в маршрутизатор
	r := router.NewRouter(database) // Обновите функцию NewRouter, чтобы принимать db
	fmt.Printf("Server started at :%s\n", port)

	// Настройка сервера с таймаутами
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Обработаем ошибку от ListenAndServe
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(ctx, "Server failed to start", err)
	}
}
