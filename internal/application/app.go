package application

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/neokofg/callap-backend/internal/application/config"
	"github.com/neokofg/callap-backend/internal/application/service"
	"github.com/neokofg/callap-backend/internal/domain/repository"
	"github.com/neokofg/callap-backend/internal/infrastructure/cache/redis"
	"github.com/neokofg/callap-backend/internal/infrastructure/database/postgresql"
	"github.com/neokofg/callap-backend/internal/infrastructure/http/fiber"
	"go.uber.org/zap"
)

func Run(cfg *config.Config, logger *zap.Logger) {
	defer logger.Info("application shut down")

	pool, err := postgresql.ConnPool(postgresql.Config{
		Username: cfg.PostgreSQL.Username,
		Password: cfg.PostgreSQL.Password,
		Host:     cfg.PostgreSQL.Host,
		Port:     cfg.PostgreSQL.Port,
		Database: cfg.PostgreSQL.Database,
		Pool: postgresql.ConfigPool{
			MaxConns:          cfg.PostgreSQL.Pool.MaxConns,
			MinConns:          cfg.PostgreSQL.Pool.MinConns,
			MaxConnLifeTime:   cfg.PostgreSQL.Pool.MaxConnLifeTime,
			MaxConnIdleTime:   cfg.PostgreSQL.Pool.MaxConnIdleTime,
			HealthCheckPeriod: cfg.PostgreSQL.Pool.HealthCheckPeriod,
		},
	})
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	logger.Info("database connected")
	rdb, err := redis.Conn(redis.Config{
		Host:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}
	logger.Info("redis connected")

	repositories := repository.NewRepositories(pool, rdb)
	services := service.NewServices(cfg, repositories, logger)

	fiber.InitFiber(cfg, logger, services)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down servers...")
}
