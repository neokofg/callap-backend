package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/neokofg/callap-backend/internal/application/service"
	"github.com/neokofg/callap-backend/internal/infrastructure/http/fiber/utils"
	"github.com/neokofg/callap-backend/pkg/validator"
	"go.uber.org/zap"
)

type FriendHandler struct {
	logger        *zap.Logger
	friendService *service.FriendService
	userService   *service.UserService
}

func NewFriendHandler(friendService *service.FriendService, userService *service.UserService, logger *zap.Logger) *FriendHandler {
	return &FriendHandler{
		logger:        logger,
		friendService: friendService,
		userService:   userService,
	}
}

type ListRequest struct {
	Limit  int `query:"limit" validate:"required,gte=1,lte=100"`
	Offset int `query:"offset" validate:"omitempty,gt=0"`
}

func (fh *FriendHandler) ListFriends(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		fh.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &ListRequest{}

	err := utils.ParseQuery(c, fh.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(fh.logger, req)
	if err != nil {
		return err
	}

	list, err := fh.friendService.List(c.Context(), userId, req.Limit, req.Offset)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponseWithData(list))
}

type DeleteRequest struct {
	FriendId string `json:"friend_id" validate:"required"`
}

func (fh *FriendHandler) Delete(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		fh.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &DeleteRequest{}

	err := utils.ParseBody(c, fh.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(fh.logger, req)
	if err != nil {
		return err
	}

	err = fh.friendService.Delete(c.Context(), userId, req.FriendId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponse())
}

type DeclineRequest struct {
	Id string `json:"id" validate:"required"`
}

func (fh *FriendHandler) Decline(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		fh.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &DeclineRequest{}

	err := utils.ParseBody(c, fh.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(fh.logger, req)
	if err != nil {
		return err
	}

	err = fh.friendService.Decline(c.Context(), userId, req.Id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponse())
}

type AcceptRequest struct {
	Id string `json:"id" validate:"required"`
}

func (fh *FriendHandler) Accept(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		fh.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &AcceptRequest{}

	err := utils.ParseBody(c, fh.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(fh.logger, req)
	if err != nil {
		return err
	}

	err = fh.friendService.Accept(c.Context(), userId, req.Id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponse())
}

func (fh *FriendHandler) GetPending(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		fh.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	pending, err := fh.friendService.GetPending(c.Context(), userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponseWithData(pending))
}

type AddFriendRequest struct {
	Nametag string `json:"nametag" validate:"required,max=255"`
}

func (fh *FriendHandler) AddFriend(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(string)
	if !exists {
		fh.logger.Warn("User ID required")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token")
	}

	req := &AddFriendRequest{}

	err := utils.ParseBody(c, fh.logger, req)
	if err != nil {
		return err
	}

	err = validator.Validate(fh.logger, req)
	if err != nil {
		return err
	}

	userFriend, err := fh.userService.GetByNameTag(c.Context(), req.Nametag)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if userFriend.Id.String() == userId {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Cannot add yourself")
	}

	err = fh.friendService.AddFriend(c.Context(), userId, userFriend.Id.String())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(utils.MakeSuccessResponse())
}
