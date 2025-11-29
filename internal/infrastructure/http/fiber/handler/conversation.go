package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/neokofg/callap-backend/internal/application/service"
	"github.com/neokofg/callap-backend/internal/infrastructure/http/fiber/utils"
	"github.com/neokofg/callap-backend/pkg/validator"
	"go.uber.org/zap"
)

type ConversationHandler struct {
	logger              *zap.Logger
	conversationService *service.ConversationService
}

func NewConversationHandler(conversationService *service.ConversationService, logger *zap.Logger) *ConversationHandler {
	return &ConversationHandler{
		logger:              logger,
		conversationService: conversationService,
	}
}

type ListMessagesRequest struct {
	Id     string `query:"id" validate:"required"`
	Limit  int    `query:"limit" validate:"required,gte=1,lte=100"`
	Offset int    `query:"offset" validate:"omitempty,gt=0"`
}

func (ch *ConversationHandler) ListMessages(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		ch.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &ListMessagesRequest{}

	err := utils.ParseQuery(c, ch.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(ch.logger, req)
	if err != nil {
		return err
	}

	messages, err := ch.conversationService.ListMessages(c.Context(), userId, req.Id, req.Limit, req.Offset)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponseWithData(messages))
}

type DeleteMessageRequest struct {
	Id string `json:"id"`
}

func (ch *ConversationHandler) DeleteMessage(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		ch.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &DeleteMessageRequest{}

	err := utils.ParseBody(c, ch.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(ch.logger, req)
	if err != nil {
		return err
	}

	err = ch.conversationService.DeleteMessage(c.Context(), userId, req.Id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponse())
}

type NewMessageRequest struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

func (ch *ConversationHandler) NewMessage(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		ch.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &NewMessageRequest{}

	err := utils.ParseBody(c, ch.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(ch.logger, req)
	if err != nil {
		return err
	}

	msg, err := ch.conversationService.NewMessage(c.Context(), userId, req.Id, req.Content)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponseWithData(msg))
}

type HideConversationRequest struct {
	Id string `json:"id"`
}

func (ch *ConversationHandler) Hide(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		ch.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &HideConversationRequest{}

	err := utils.ParseBody(c, ch.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(ch.logger, req)
	if err != nil {
		return err
	}

	err = ch.conversationService.Hide(c.Context(), userId, req.Id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponse())
}

type GetConversationRequest struct {
	Id string `query:"id"`
}

func (ch *ConversationHandler) GetConversation(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		ch.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &GetConversationRequest{}

	err := utils.ParseQuery(c, ch.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(ch.logger, req)
	if err != nil {
		return err
	}

	conv, err := ch.conversationService.GetConversationById(c.Context(), userId, req.Id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponseWithData(conv))
}

func (ch *ConversationHandler) ListConversations(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		ch.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &ListRequest{}

	err := utils.ParseQuery(c, ch.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(ch.logger, req)
	if err != nil {
		return err
	}

	convs, err := ch.conversationService.List(c.Context(), userId, req.Limit, req.Offset)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponseWithData(convs))
}

type GetOrCreateRequest struct {
	TargetId string `json:"target_id"`
}

func (ch *ConversationHandler) GetOrCreate(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		ch.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &GetOrCreateRequest{}

	err := utils.ParseBody(c, ch.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(ch.logger, req)
	if err != nil {
		return err
	}

	convId, err := ch.conversationService.GetOrCreate(c.Context(), userId, req.TargetId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponseWithData(convId))
}
