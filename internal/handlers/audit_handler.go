package handlers

import (
	"strconv"

	"github.com/FRFebi/bot-management-backend/internal/database"
	"github.com/FRFebi/bot-management-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type AuditHandler struct{}

func NewAuditHandler() *AuditHandler {
	return &AuditHandler{}
}

func (h *AuditHandler) GetAuditLogs(c *fiber.Ctx) error {
	var audits []models.AuditLog

	query := database.DB.Preload("User").Order("created_at DESC")

	// Optional filter by user ID
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Optional filter by action
	if action := c.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}

	// Pagination
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if err := query.Limit(limit).Offset(offset).Find(&audits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch audit logs",
		})
	}

	return c.JSON(audits)
}

func (h *AuditHandler) GetAuditLog(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid audit log ID",
		})
	}

	var audit models.AuditLog
	if err := database.DB.Preload("User").First(&audit, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Audit log not found",
		})
	}

	return c.JSON(audit)
}
