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

func (us *UserService) GetByNameTag(c context.Context, nametag string) (entity.User, error) {
	c, cancel := context.WithTimeout(c, us.cTimeout)
	defer cancel()

	return us.repo.GetByNameTag(c, nametag)
}

func (us *UserService) GetById(c context.Context, id string) (entity.User, error) {
	c, cancel := context.WithTimeout(c, us.cTimeout)
	defer cancel()

	return us.repo.GetById(c, id)
}

func (us *UserService) GetByEmail(c context.Context, email string) (entity.User, error) {
	c, cancel := context.WithTimeout(c, us.cTimeout)
	defer cancel()

	return us.repo.GetByEmail(c, email)
}

func (us *UserService) Create(c context.Context, user entity.User) (entity.User, error) {
	c, cancel := context.WithTimeout(c, us.cTimeout)
	defer cancel()

	return us.repo.Create(c, user)
}

func (us *UserService) Update(c context.Context, user entity.User) (entity.User, error) {
	c, cancel := context.WithTimeout(c, us.cTimeout)
	defer cancel()

	return us.repo.Update(c, user)
}
