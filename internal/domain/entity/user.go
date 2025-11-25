package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type User struct {
	Id           ulid.ULID
	Name         string
	Tag          string
	Email        string
	RefreshToken string
	Password     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(u User) User {
	createdAt := time.Now()
	if !u.CreatedAt.IsZero() {
		createdAt = u.CreatedAt
	}

	updatedAt := time.Now()
	if !u.UpdatedAt.IsZero() {
		updatedAt = u.UpdatedAt
	}

	id := ulid.Make()
	if u.Id != ulid.Zero {
		id = u.Id
	}

	return User{
		Id:           id,
		Email:        u.Email,
		Password:     u.Password,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		RefreshToken: u.RefreshToken,
	}
}
