package handlers

import (
	"strconv"

	"github.com/FRFebi/bot-management-backend/internal/database"
	"github.com/FRFebi/bot-management-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

type BotHandler struct{}

func NewBotHandler() *BotHandler {
	return &BotHandler{}
}

type CreateBotRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Version     string         `json:"version"`
	Config      datatypes.JSON `json:"config"`
}

type UpdateBotRequest struct {
	Name        *string         `json:"name,omitempty"`
	Description *string         `json:"description,omitempty"`
	Version     *string         `json:"version,omitempty"`
	Config      *datatypes.JSON `json:"config,omitempty"`
}

func (h *BotHandler) GetBots(c *fiber.Ctx) error {
	var bots []models.Bot
	if err := database.DB.Find(&bots).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch bots",
		})
	}

	return c.JSON(bots)
}

func (h *BotHandler) GetBot(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bot ID",
		})
	}

	var bot models.Bot
	if err := database.DB.First(&bot, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Bot not found",
		})
	}

	return c.JSON(bot)
}

func (h *BotHandler) CreateBot(c *fiber.Ctx) error {
	var req CreateBotRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bot name is required",
		})
	}

	bot := models.Bot{
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		Config:      req.Config,
		Status:      "stopped",
	}

	if err := database.DB.Create(&bot).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create bot",
		})
	}

	// Log audit
	h.logAudit(c, "bot.create", fiber.Map{
		"bot_id":   bot.ID,
		"bot_name": bot.Name,
	})

	return c.Status(fiber.StatusCreated).JSON(bot)
}

func (h *BotHandler) UpdateBot(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bot ID",
		})
	}

	var bot models.Bot
	if err := database.DB.First(&bot, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Bot not found",
		})
	}

	var req UpdateBotRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name != nil {
		bot.Name = *req.Name
	}
	if req.Description != nil {
		bot.Description = *req.Description
	}
	if req.Version != nil {
		bot.Version = *req.Version
	}
	if req.Config != nil {
		bot.Config = *req.Config
	}

	if err := database.DB.Save(&bot).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update bot",
		})
	}

	// Log audit
	h.logAudit(c, "bot.update", fiber.Map{
		"bot_id":   bot.ID,
		"bot_name": bot.Name,
	})

	return c.JSON(bot)
}

func (h *BotHandler) DeleteBot(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bot ID",
		})
	}

	var bot models.Bot
	if err := database.DB.First(&bot, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Bot not found",
		})
	}

	if err := database.DB.Delete(&bot).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete bot",
		})
	}

	// Log audit
	h.logAudit(c, "bot.delete", fiber.Map{
		"bot_id":   bot.ID,
		"bot_name": bot.Name,
	})

	return c.JSON(fiber.Map{
		"message": "Bot deleted successfully",
	})
}

func (h *BotHandler) StartBot(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bot ID",
		})
	}

	var bot models.Bot
	if err := database.DB.First(&bot, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Bot not found",
		})
	}

	if bot.Status == "running" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bot is already running",
		})
	}

	bot.Status = "running"
	if err := database.DB.Save(&bot).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start bot",
		})
	}

	// Log audit
	h.logAudit(c, "bot.start", fiber.Map{
		"bot_id":   bot.ID,
		"bot_name": bot.Name,
	})

	return c.JSON(fiber.Map{
		"message": "Bot started successfully",
		"bot":     bot,
	})
}

func (h *BotHandler) StopBot(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bot ID",
		})
	}

	var bot models.Bot
	if err := database.DB.First(&bot, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Bot not found",
		})
	}

	if bot.Status == "stopped" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bot is already stopped",
		})
	}

	bot.Status = "stopped"
	if err := database.DB.Save(&bot).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to stop bot",
		})
	}

	// Log audit
	h.logAudit(c, "bot.stop", fiber.Map{
		"bot_id":   bot.ID,
		"bot_name": bot.Name,
	})

	return c.JSON(fiber.Map{
		"message": "Bot stopped successfully",
		"bot":     bot,
	})
}

func (h *BotHandler) RestartBot(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bot ID",
		})
	}

	var bot models.Bot
	if err := database.DB.First(&bot, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Bot not found",
		})
	}

	bot.Status = "running"
	if err := database.DB.Save(&bot).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to restart bot",
		})
	}

	// Log audit
	h.logAudit(c, "bot.restart", fiber.Map{
		"bot_id":   bot.ID,
		"bot_name": bot.Name,
	})

	return c.JSON(fiber.Map{
		"message": "Bot restarted successfully",
		"bot":     bot,
	})
}

func (h *BotHandler) DeployBot(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bot ID",
		})
	}

	var bot models.Bot
	if err := database.DB.First(&bot, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Bot not found",
		})
	}

	var req struct {
		Version string `json:"version"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Version == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Version is required",
		})
	}

	bot.Version = req.Version
	if err := database.DB.Save(&bot).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to deploy bot",
		})
	}

	// Log audit
	h.logAudit(c, "bot.deploy", fiber.Map{
		"bot_id":   bot.ID,
		"bot_name": bot.Name,
		"version":  req.Version,
	})

	return c.JSON(fiber.Map{
		"message": "Bot deployed successfully",
		"bot":     bot,
	})
}

func (h *BotHandler) GetBotStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bot ID",
		})
	}

	var bot models.Bot
	if err := database.DB.First(&bot, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Bot not found",
		})
	}

	return c.JSON(fiber.Map{
		"id":      bot.ID,
		"name":    bot.Name,
		"status":  bot.Status,
		"version": bot.Version,
	})
}

func (h *BotHandler) logAudit(c *fiber.Ctx, action string, details fiber.Map) {
	userID := c.Locals("userID")
	if userID == nil {
		return
	}

	detailsJSON, _ := datatypes.NewJSONType(details).MarshalJSON()

	audit := models.AuditLog{
		UserID:  toUintPtr(userID.(uint)),
		Action:  action,
		Details: datatypes.JSON(detailsJSON),
	}

	database.DB.Create(&audit)
}

func toUintPtr(val uint) *uint {
	return &val
}