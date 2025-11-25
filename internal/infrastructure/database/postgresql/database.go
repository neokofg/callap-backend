package postgresql

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Pool     ConfigPool
}

type ConfigPool struct {
	MaxConns          int
	MinConns          int
	MaxConnLifeTime   int
	MaxConnIdleTime   int
	HealthCheckPeriod int
}

func Conn(config Config) (*pgx.Conn, error) {
	password := url.QueryEscape(config.Password)
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.Username,
		password,
		config.Host,
		config.Port,
		config.Database,
	)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ConnPool(config Config) (*pgxpool.Pool, error) {
	password := url.QueryEscape(config.Password)
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.Username,
		password,
		config.Host,
		config.Port,
		config.Database,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(config.Pool.MaxConns)
	poolConfig.MinConns = int32(config.Pool.MinConns)
	poolConfig.MaxConnLifetime = time.Duration(config.Pool.MaxConnLifeTime) * time.Second
	poolConfig.MaxConnIdleTime = time.Duration(config.Pool.MaxConnIdleTime) * time.Second
	poolConfig.HealthCheckPeriod = time.Duration(config.Pool.HealthCheckPeriod) * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
