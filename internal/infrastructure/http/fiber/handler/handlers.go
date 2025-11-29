package handler

import (
	"github.com/neokofg/callap-backend/internal/application/service"
	"go.uber.org/zap"
)

type Handlers struct {
	AuthHandler         *AuthHandler
	UserHandler         *UserHandler
	FriendHandler       *FriendHandler
	ConversationHandler *ConversationHandler
	WebsocketHandler    *WebsocketHandler
}

func NewHandlers(services *service.Services, logger *zap.Logger) *Handlers {
	return &Handlers{
		AuthHandler:         NewAuthHandler(services.JWT, services.UserService, services.PasswordService, logger),
		UserHandler:         NewUserHandler(services.UserService, logger),
		FriendHandler:       NewFriendHandler(services.FriendService, services.UserService, logger),
		ConversationHandler: NewConversationHandler(services.ConversationService, logger),
		WebsocketHandler:    NewWebsocketHandler(services.WebsocketService, logger),
	}
}
