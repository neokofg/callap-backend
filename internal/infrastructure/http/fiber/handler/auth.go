package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neokofg/callap-backend/internal/application/context"
	"github.com/neokofg/callap-backend/pkg/validator"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger     *zap.Logger
	jwtService context.JwtService
}

func NewAuthHandler(jwtService context.JwtService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		logger:     logger,
		jwtService: jwtService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func (ah *AuthHandler) Register(c *fiber.Ctx) error {
	req := &RegisterRequest{}

	if err := c.BodyParser(req); err != nil {
		ah.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid JSON",
		})
	}

	validationErrors := validator.Validate(req)
	if len(validationErrors) > 0 {
		errMsgs := make([]string, 0, len(validationErrors))
		for _, err := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf("[%s]: %s", err.FailedField, err.Tag))
		}
		ah.logger.Warn("Validation failed", zap.Strings("errors", errMsgs))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  strings.Join(errMsgs, ", "),
		})
	}

	return c.Status(fiber.StatusOK).JSON()
}
