package service

import (
	"context"
	"time"

	"github.com/neokofg/callap-backend/internal/domain/entity"
	"github.com/neokofg/callap-backend/internal/domain/repository"
)

type FriendService struct {
	cTimeout time.Duration
	repo     *repository.FriendRepository
}

func NewFriendService(cTimeout time.Duration, repo *repository.FriendRepository) *FriendService {
	return &FriendService{
		cTimeout: cTimeout,
		repo:     repo,
	}
}

func (fs *FriendService) List(c context.Context, userId string, limit int, offset int) ([]entity.FriendUser, error) {
	c, cancel := context.WithTimeout(c, fs.cTimeout)
	defer cancel()

	return fs.repo.List(c, userId, limit, offset)
}

func (fs *FriendService) Delete(c context.Context, userId string, friendId string) error {
	c, cancel := context.WithTimeout(c, fs.cTimeout)
	defer cancel()

	return fs.repo.Delete(c, userId, friendId)
}

func (fs *FriendService) Decline(c context.Context, userId string, id string) error {
	c, cancel := context.WithTimeout(c, fs.cTimeout)
	defer cancel()

	return fs.repo.Decline(c, userId, id)
}

func (fs *FriendService) Accept(c context.Context, userId string, id string) error {
	c, cancel := context.WithTimeout(c, fs.cTimeout)
	defer cancel()

	return fs.repo.Accept(c, userId, id)
}

func (fs *FriendService) GetPending(c context.Context, userId string) ([]*entity.PendingFriend, error) {
	c, cancel := context.WithTimeout(c, fs.cTimeout)
	defer cancel()

	return fs.repo.GetPending(c, userId)
}

func (fs *FriendService) AddFriend(c context.Context, userId string, friendId string) error {
	c, cancel := context.WithTimeout(c, fs.cTimeout)
	defer cancel()

	return fs.repo.AddFriend(c, userId, friendId)
}
