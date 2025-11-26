package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/neokofg/callap-backend/internal/application/service"
	"github.com/neokofg/callap-backend/internal/infrastructure/http/fiber/utils"
	"go.uber.org/zap"
)

type UserHandler struct {
	logger      *zap.Logger
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (uh *UserHandler) Me(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		uh.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	user, err := uh.userService.GetById(c.Context(), userId)
	if err != nil {
		uh.logger.Warn("User not found", zap.String("userId", userId), zap.Error(err))
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponseWithData(user))
}
