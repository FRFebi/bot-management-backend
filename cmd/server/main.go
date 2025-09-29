package main

import (
	"fmt"

	"github.com/FRFebi/bot-management-backend/internal/config"
	"github.com/FRFebi/bot-management-backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize config
	cfg := config.New()

	// Initialize logger
	log := logger.New()

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
		return c.JSON(fiber.Map{
			"status":     "ok",
			"service":    "bot-management-backend",
			"version":    "1.0.0",
			"environment": cfg.Server.Env,
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