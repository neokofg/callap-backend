package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neokofg/callap-backend/internal/domain/entity"
)

type UserRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (ur *UserRepository) GetById(c context.Context, id string) (entity.User, error) {
	var user entity.User
	query := fmt.Sprintf(
		"SELECT id, name, tag, email, password, created_at, updated_at, refresh_token FROM %s WHERE id = $1",
		ur.tableName,
	)
	err := ur.pool.QueryRow(c, query, id).
		Scan(&user.Id, &user.Name, &user.Tag, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.RefreshToken)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (ur *UserRepository) Create(c context.Context, user entity.User) (entity.User, error) {
	newUser := entity.NewUser(user)
	query := fmt.Sprintf(
		"INSERT INTO %s (id, name, tag, email, password, created_at, updated_at, refresh_token) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		ur.tableName,
	)
	_, err := ur.pool.Exec(
		c, query,
		newUser.Id, newUser.Name, newUser.Tag, newUser.Email, newUser.Password, newUser.CreatedAt, newUser.UpdatedAt, newUser.RefreshToken,
	)
	if err != nil {
		return entity.User{}, err
	}
	return newUser, nil
}
