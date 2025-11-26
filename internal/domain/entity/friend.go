package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Status string

const (
	Pending  Status = "pending"
	Accepted Status = "accepted"
	Rejected Status = "rejected"
	Blocked  Status = "blocked"
)

type Friend struct {
	Id        ulid.ULID
	UserId    ulid.ULID
	FriendId  ulid.ULID
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewFriend(f Friend) Friend {
	id := ulid.Make()
	if f.Id != ulid.Zero {
		id = f.Id
	}
	return Friend{
		Id:        id,
		UserId:    f.UserId,
		FriendId:  f.FriendId,
		Status:    f.Status,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

type PendingFriend struct {
	ID       string `json:"id"`
	SenderID string `json:"sender_id"`
	Name     string `json:"name"`
	Tag      string `json:"tag"`
}
