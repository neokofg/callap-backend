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
	id := ulid.Make()
	if u.Id != ulid.Zero {
		id = u.Id
	}

	return User{
		Id:           id,
		Name:         u.Name,
		Tag:          u.Tag,
		Email:        u.Email,
		Password:     u.Password,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		RefreshToken: u.RefreshToken,
	}
}

type FriendUser struct {
	Id   ulid.ULID
	Name string
	Tag  string
}
