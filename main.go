package main

import (
	"context"
	"fmt"
	"net/http"

	"music/config"
	"music/internal/db"

	"music/internal/router"
	"music/pkg/logger"
)

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
