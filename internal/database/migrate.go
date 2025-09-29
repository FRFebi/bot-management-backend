package database

import (
	"fmt"

	"github.com/FRFebi/bot-management-backend/internal/models"
)

func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.Bot{},
		&models.Schedule{},
		&models.Run{},
		&models.AuditLog{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}