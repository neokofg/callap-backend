package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	UserRepository *UserRepository
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		UserRepository: NewUserRepository(pool),
	}
}
