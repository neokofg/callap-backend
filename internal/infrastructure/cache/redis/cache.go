package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Password string
	DB       int
}

func Conn(config Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
