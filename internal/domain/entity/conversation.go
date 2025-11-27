package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Type string

const (
	Private Type = "private"
	Group   Type = "group"
)

type Conversation struct {
	Id        ulid.ULID
	Type      Type
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ConversationSummary struct {
	ID            string     `json:"id"`
	OtherUserID   string     `json:"other_user_id"`
	OtherUserName string     `json:"other_user_name"`
	OtherUserTag  string     `json:"other_user_tag"`
	LastMessage   *string    `json:"last_message"`
	LastMessageAt *time.Time `json:"last_message_at"`
	UnreadCount   int        `json:"unread_count"`
}

type ConversationDetails struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Participants []string  `json:"participants"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
