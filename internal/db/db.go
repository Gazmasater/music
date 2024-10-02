// db/db.go
package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect - функция для подключения к базе данных
func Connect() (*gorm.DB, error) {
	// Строка подключения к базе данных из переменных окружения
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err // Возвращаем ошибку, если подключение не удалось
	}

	return db, nil // Возвращаем подключение, если успешно
}
