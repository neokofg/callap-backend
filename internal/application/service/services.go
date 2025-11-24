package service

import (
	"github.com/neokofg/callap-backend/internal/application/config"
	"github.com/neokofg/callap-backend/internal/application/context"
	"github.com/neokofg/callap-backend/pkg/jwt"
)

type Services struct {
	JWT context.JwtService
}

func NewServices(cfg *config.Config) *Services {
	return &Services{
		JWT: jwt.NewService(jwt.Config(cfg.JWT)),
	}
}
