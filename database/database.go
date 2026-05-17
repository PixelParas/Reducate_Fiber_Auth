package database

import (
	"fiber_auth/config"
	"fiber_auth/models"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) error {

	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		return fmt.Errorf("Failed to connect to databse: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("Failed to get the underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to AutoMigrate: %w", err)
	}

	DB = db
	return nil
}
