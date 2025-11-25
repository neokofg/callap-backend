package main

import (
	"github.com/neokofg/callap-backend/internal/application"
	"github.com/neokofg/callap-backend/internal/application/config"
	"github.com/neokofg/callap-backend/pkg/zap"
)

func main() {
	logger := zap.InitZap()
	cfg := config.InitConfig(logger)

	application.Run(cfg, logger)
}
