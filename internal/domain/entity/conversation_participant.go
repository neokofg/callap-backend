package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type ConversationParticipant struct {
	Id             ulid.ULID
	ConversationId ulid.ULID
	UserId         ulid.ULID
	JoinedAt       time.Time
}
