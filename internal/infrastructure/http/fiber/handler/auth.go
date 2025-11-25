package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/neokofg/callap-backend/internal/application/service"
	"github.com/neokofg/callap-backend/internal/application/utils"
	"github.com/neokofg/callap-backend/internal/domain/entity"
	"github.com/neokofg/callap-backend/pkg/jwt"
	"github.com/neokofg/callap-backend/pkg/validator"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger          *zap.Logger
	jwtService      *jwt.Service
	userService     *service.UserService
	passwordService *service.PasswordService
}

func NewAuthHandler(
	jwtService *jwt.Service,
	userService *service.UserService,
	passwordService *service.PasswordService,
	logger *zap.Logger,
) *AuthHandler {
	return &AuthHandler{
		logger:          logger,
		jwtService:      jwtService,
		userService:     userService,
		passwordService: passwordService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=255"`
}

func (ah *AuthHandler) Login(c *fiber.Ctx) error {
	req := &LoginRequest{}

	err := utils.ParseBody(c, ah.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(ah.logger, req)
	if err != nil {
		return err
	}

	user, err := ah.userService.GetByEmail(c.Context(), req.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid password or email")
	}

	if err = ah.passwordService.CheckPassword(req.Password, user.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid password or email")
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
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}
	refreshBody := map[string]interface{}{
		"type": "refresh",
	}

	refreshToken, err := ah.jwtService.GenerateToken(user.Id, refreshBody)
	if err != nil {
		ah.logger.Error("Failed to generate refresh token", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate refresh token")
	}
	user.RefreshToken = refreshToken
	_, err = ah.userService.Update(c.Context(), user)
	if err != nil {
		ah.logger.Error("Failed to update user", zap.Error(err))
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Failed to update user")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"access_token":  token,
			"refresh_token": refreshToken,
		},
	})
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
	Tag      string `json:"tag" validate:"required,min=3,max=3"`
	Email    string `json:"email" validate:"required,email,max=255"`
}

func (ah *AuthHandler) Register(c *fiber.Ctx) error {
	req := &RegisterRequest{}
	err := utils.ParseBody(c, ah.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(ah.logger, req)
	if err != nil {
		return err
	}

	hashedPassword, err := ah.passwordService.HashPassword(req.Password)
	if err != nil {
		return err
	}

	var user = entity.User{
		Name:     req.Name,
		Tag:      req.Tag,
		Password: hashedPassword,
		Email:    req.Email,
	}

	user, err = ah.userService.Create(c.Context(), user)
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
	refreshBody := map[string]interface{}{
		"type": "refresh",
	}

	refreshToken, err := ah.jwtService.GenerateToken(user.Id, refreshBody)
	if err != nil {
		ah.logger.Error("Failed to generate refresh token", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"errors":  err.Error(),
		})
	}
	user.RefreshToken = refreshToken
	_, err = ah.userService.Update(c.Context(), user)
	if err != nil {
		ah.logger.Error("Failed to update user", zap.Error(err))
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"access_token":  token,
			"refresh_token": refreshToken,
		},
	})
}
