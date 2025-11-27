package service

import (
	"time"

	"github.com/neokofg/callap-backend/internal/application/config"
	"github.com/neokofg/callap-backend/internal/domain/repository"
	"github.com/neokofg/callap-backend/pkg/jwt"
)

type Services struct {
	JWT                 *jwt.Service
	UserService         *UserService
	PasswordService     *PasswordService
	FriendService       *FriendService
	ConversationService *ConversationService
}

func NewServices(cfg *config.Config, repositories *repository.Repositories) *Services {
	c := time.Duration(cfg.ContextTimeout) * time.Second

	return &Services{
		JWT:                 jwt.NewService(jwt.Config(cfg.JWT)),
		UserService:         NewUserService(c, repositories.UserRepository),
		PasswordService:     NewPasswordService(),
		FriendService:       NewFriendService(c, repositories.FriendRepository),
		ConversationService: NewConversationService(c, repositories.ConversationRepository),
	}
}
