package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neokofg/callap-backend/internal/application/service"
	"github.com/neokofg/callap-backend/internal/domain/entity"
	"github.com/neokofg/callap-backend/pkg/jwt"
	"github.com/neokofg/callap-backend/pkg/validator"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger      *zap.Logger
	jwtService  *jwt.Service
	userService *service.UserService
}

func NewAuthHandler(jwtService *jwt.Service, userService *service.UserService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		jwtService:  jwtService,
		userService: userService,
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

	var user = entity.User{
		Name:     req.Name,
		Password: req.Password,
	}

	user, err := ah.userService.Create(c.Context(), user)
	if err != nil {
		ah.logger.Error("Failed to create user", zap.Error(err))
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err.Error(),
		})
	}

	body := map[string]interface{}{
		"name":       user.Name,
		"tag":        user.Tag,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	token, err := ah.jwtService.GenerateToken(user.Id, body)
	if err != nil {
		ah.logger.Error("Failed to generate token", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"token": token,
		},
	})
}
