package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neokofg/callap-backend/internal/domain/entity"
)

type UserRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool:      pool,
		tableName: userTableName,
	}
}

func (ur *UserRepository) GetByNameTag(c context.Context, nametag string) (entity.User, error) {
	var user entity.User
	name, tag, found := strings.Cut(nametag, "#")
	if !found {
		return entity.User{}, errors.New("invalid nametag")
	}

	query := fmt.Sprintf(
		"SELECT id, name, tag, email, created_at, updated_at FROM %s WHERE name = $1 AND tag = $2",
		ur.tableName,
	)

	err := ur.pool.QueryRow(c, query, name, tag).
		Scan(&user.Id, &user.Name, &user.Tag, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
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

func (ur *UserRepository) GetByEmail(c context.Context, email string) (entity.User, error) {
	var user entity.User
	query := fmt.Sprintf(
		"SELECT id, name, tag, email, password, created_at, updated_at, refresh_token FROM %s WHERE email = $1",
		ur.tableName,
	)
	err := ur.pool.QueryRow(c, query, email).
		Scan(&user.Id, &user.Name, &user.Tag, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.RefreshToken)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (ur *UserRepository) Create(c context.Context, user entity.User) (entity.User, error) {
	newUser := entity.NewUser(user)
	query := fmt.Sprintf(
		"INSERT INTO %s (id, name, tag, email, password) VALUES ($1, $2, $3, $4, $5)",
		ur.tableName,
	)
	_, err := ur.pool.Exec(
		c, query,
		newUser.Id.String(), newUser.Name, newUser.Tag, newUser.Email, newUser.Password,
	)
	if err != nil {
		return entity.User{}, err
	}
	return newUser, nil
}

func (ur *UserRepository) Update(c context.Context, user entity.User) (entity.User, error) {
	query := fmt.Sprintf(
		"UPDATE %s SET email = $2, tag = $3, name = $4, password = $5, refresh_token = $6 WHERE id = $1",
		ur.tableName,
	)
	_, err := ur.pool.Exec(
		c, query,
		user.Id.String(), user.Email, user.Tag, user.Name, user.Password, user.RefreshToken,
	)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
