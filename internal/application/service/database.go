package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseService struct {
	cTimeout time.Duration
	pool     *pgxpool.Pool
}

func (ds *DatabaseService) StartTransaction(c context.Context) (*pgx.Tx, error) {
	tx, err := ds.pool.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if tx != nil {
			tx.Rollback(c)
		}
	}()

	return &tx, nil
}
