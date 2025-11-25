package service

import (
	"context"
	"time"

	"github.com/neokofg/callap-backend/internal/domain/entity"
	"github.com/neokofg/callap-backend/internal/domain/repository"
)

type UserService struct {
	cTimeout time.Duration
	repo     *repository.UserRepository
}

func NewUserService(cTimeout time.Duration, repo *repository.UserRepository) *UserService {
	return &UserService{
		cTimeout: cTimeout,
		repo:     repo,
	}
}

func (us *UserService) GetById(c context.Context, id string) (entity.User, error) {
	c, cancel := context.WithTimeout(c, us.cTimeout)
	defer cancel()

	return us.repo.GetById(c, id)
}

func (us *UserService) Create(c context.Context, user entity.User) (entity.User, error) {
	c, cancel := context.WithTimeout(c, us.cTimeout)
	defer cancel()

	return us.repo.Create(c, user)
}
