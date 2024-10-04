package db

import (
	"context"

	"music/internal/models"

	"gorm.io/gorm"

	"music/pkg/logger"
)

func Migrate(db *gorm.DB) {
	// Выполняем миграции для моделей
	err := db.AutoMigrate(&models.Artist{}, &models.SongDetail{})
	if err != nil {
		logger.Fatal(context.Background(), "failed to migrate database", err)
	}
	logger.Info(context.Background(), "Database migrated successfully!")
}
