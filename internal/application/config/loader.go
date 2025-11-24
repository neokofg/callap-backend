package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func InitConfig(logger *zap.Logger) *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		err := godotenv.Load()
		if err == nil {
			configPath = ".env"
		}
	}

	var cfg Config

	if configPath != "" {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			logger.Fatal("Cannot read .env file: %v", zap.Error(err))
		}
	} else {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			logger.Fatal("Cannot read .env param: %v", zap.Error(err))
		}
	}

	return &cfg
}
