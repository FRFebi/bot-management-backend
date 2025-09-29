package main

import (
	"fmt"
	"log"

	"github.com/FRFebi/bot-management-backend/internal/config"
	"github.com/FRFebi/bot-management-backend/internal/database"
	"github.com/FRFebi/bot-management-backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize config
	cfg := config.New()

	// Initialize logger
	log := logger.New()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Info("Database connected successfully")

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Info("Database migrated successfully")

	// Seed database in development
	if cfg.Server.Env == "development" {
		if err := database.Seed(); err != nil {
			log.Errorf("Failed to seed database: %v", err)
		} else {
			log.Info("Database seeded successfully")
		}
	}

	// Ensure database connection is closed on exit
	defer func() {
		if err := database.Close(); err != nil {
			log.Errorf("Failed to close database connection: %v", err)
		}
	}()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Bot Management Backend",
	})

	// Add middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "ok"
		if err := database.HealthCheck(); err != nil {
			dbStatus = "error"
		}

		return c.JSON(fiber.Map{
			"status":      "ok",
			"service":     "bot-management-backend",
			"version":     "1.0.0",
			"environment": cfg.Server.Env,
			"database":    dbStatus,
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Placeholder route
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Bot Management API v1.0",
		})
	})

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Infof("Server starting on %s in %s mode", addr, cfg.Server.Env)

	if err := app.Listen(addr); err != nil {
		log.Errorf("Failed to start server: %v", err)
	}
}