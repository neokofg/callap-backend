package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	UserRepository         *UserRepository
	FriendRepository       *FriendRepository
	ConversationRepository *ConversationRepository
}

func NewRepositories(pool *pgxpool.Pool, rdb *redis.Client) *Repositories {
	return &Repositories{
		UserRepository:         NewUserRepository(pool),
		FriendRepository:       NewFriendRepository(pool),
		ConversationRepository: NewConversationRepository(pool, rdb),
	}
}
