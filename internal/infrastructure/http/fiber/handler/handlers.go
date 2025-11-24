package handler

import (
	"github.com/neokofg/callap-backend/internal/application/service"
	"go.uber.org/zap"
)

type Handlers struct {
	AuthHandler *AuthHandler
}

func NewHandlers(services *service.Services, logger *zap.Logger) *Handlers {
	return &Handlers{
		AuthHandler: NewAuthHandler(services.JWT, logger),
	}
}
