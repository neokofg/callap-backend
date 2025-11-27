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

func Conn() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
