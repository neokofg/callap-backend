package zap

import (
	"log"

	"go.uber.org/zap"
)

func InitZap() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	return logger
}
