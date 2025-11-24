package application

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/neokofg/callap-backend/internal/application/config"
	"github.com/neokofg/callap-backend/internal/application/service"
	"github.com/neokofg/callap-backend/internal/infrastructure/http/fiber"
	"go.uber.org/zap"
)

func Run(cfg *config.Config, logger *zap.Logger) {
	defer logger.Info("application shut down")

	services := service.NewServices(cfg)

	fiber.InitFiber(cfg, logger, services)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down servers...")
}
