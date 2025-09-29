package database

import (
	"fmt"

	"github.com/FRFebi/bot-management-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

func Seed() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	var userCount int64
	if err := DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if userCount > 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	adminUser := &models.User{
		Name:         "Admin User",
		Email:        "admin@example.com",
		PasswordHash: string(hashedPassword),
		Role:         "admin",
	}

	if err := DB.Create(adminUser).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	viewerPassword, err := bcrypt.GenerateFromPassword([]byte("viewer123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash viewer password: %w", err)
	}

	viewerUser := &models.User{
		Name:         "Viewer User",
		Email:        "viewer@example.com",
		PasswordHash: string(viewerPassword),
		Role:         "viewer",
	}

	if err := DB.Create(viewerUser).Error; err != nil {
		return fmt.Errorf("failed to create viewer user: %w", err)
	}

	sampleBot := &models.Bot{
		Name:        "Sample Scraper Bot",
		Description: "A sample web scraping bot for demonstration",
		Version:     "1.0.0",
		Config: datatypes.JSON([]byte(`{
			"target_url": "https://example.com",
			"timeout": 30,
			"retry_count": 3
		}`)),
		Status: "stopped",
	}

	if err := DB.Create(sampleBot).Error; err != nil {
		return fmt.Errorf("failed to create sample bot: %w", err)
	}

	sampleSchedule := &models.Schedule{
		BotID:          sampleBot.ID,
		CronExpression: "0 0 * * *",
		IsActive:       true,
	}

	if err := DB.Create(sampleSchedule).Error; err != nil {
		return fmt.Errorf("failed to create sample schedule: %w", err)
	}

	return nil
}