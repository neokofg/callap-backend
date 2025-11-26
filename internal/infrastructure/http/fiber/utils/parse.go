package utils

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ParseBody(c *fiber.Ctx, logger *zap.Logger, req interface{}) error {
	if err := c.BodyParser(req); err != nil {
		logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid JSON",
		})
	}
	return nil
}

func ParseQuery(c *fiber.Ctx, logger *zap.Logger, req interface{}) error {
	if err := c.QueryParser(req); err != nil {
		logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid Query",
		})
	}
	return nil
}
