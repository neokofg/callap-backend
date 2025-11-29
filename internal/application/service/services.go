package service

import (
	"time"

	"github.com/neokofg/callap-backend/internal/application/config"
	"github.com/neokofg/callap-backend/internal/domain/repository"
	"github.com/neokofg/callap-backend/pkg/jwt"
	"go.uber.org/zap"
)

type Services struct {
	JWT                 *jwt.Service
	UserService         *UserService
	PasswordService     *PasswordService
	FriendService       *FriendService
	ConversationService *ConversationService
	WebsocketService    *WebsocketService
}

func NewServices(cfg *config.Config, repositories *repository.Repositories, logger *zap.Logger) *Services {
	c := time.Duration(cfg.ContextTimeout) * time.Second

	wsService := NewWebsocketService(c, logger)

	return &Services{
		JWT:                 jwt.NewService(jwt.Config(cfg.JWT)),
		UserService:         NewUserService(c, repositories.UserRepository),
		PasswordService:     NewPasswordService(),
		FriendService:       NewFriendService(c, repositories.FriendRepository),
		ConversationService: NewConversationService(c, repositories.ConversationRepository, wsService),
		WebsocketService:    wsService,
	}
}
