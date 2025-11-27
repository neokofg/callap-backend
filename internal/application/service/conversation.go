package service

import (
	"context"
	"time"

	"github.com/neokofg/callap-backend/internal/domain/entity"
	"github.com/neokofg/callap-backend/internal/domain/repository"
)

type ConversationService struct {
	cTimeout time.Duration
	repo     *repository.ConversationRepository
}

func NewConversationService(timeout time.Duration, repo *repository.ConversationRepository) *ConversationService {
	return &ConversationService{
		cTimeout: timeout,
		repo:     repo,
	}
}

func (cs *ConversationService) ListMessages(c context.Context, userId string, id string, limit int, offest int) ([]entity.Message, error) {
	c, cancel := context.WithTimeout(c, cs.cTimeout)
	defer cancel()

	return cs.repo.ListMessages(c, userId, id, limit, offest)
}

func (cs *ConversationService) DeleteMessage(c context.Context, userId string, id string) error {
	c, cancel := context.WithTimeout(c, cs.cTimeout)
	defer cancel()

	return cs.repo.DeleteMessage(c, userId, id)
}

func (cs *ConversationService) NewMessage(c context.Context, userId string, id string, content string) error {
	c, cancel := context.WithTimeout(c, cs.cTimeout)
	defer cancel()

	return cs.repo.NewMessage(c, userId, id, content)
}

func (cs *ConversationService) Hide(c context.Context, userId string, id string) error {
	c, cancel := context.WithTimeout(c, cs.cTimeout)
	defer cancel()

	return cs.repo.Hide(c, userId, id)
}

func (cs *ConversationService) GetConversationById(c context.Context, userId string, convId string) (*entity.ConversationDetails, error) {
	c, cancel := context.WithTimeout(c, cs.cTimeout)
	defer cancel()

	return cs.repo.GetConversationByID(c, userId, convId)
}

func (cs *ConversationService) List(c context.Context, userId string, limit int, offest int) ([]entity.ConversationSummary, error) {
	c, cancel := context.WithTimeout(c, cs.cTimeout)
	defer cancel()

	return cs.repo.List(c, userId, limit, offest)
}

func (cs *ConversationService) GetOrCreate(c context.Context, userId string, targetId string) (string, error) {
	c, cancel := context.WithTimeout(c, cs.cTimeout)
	defer cancel()

	return cs.repo.GetOrCreate(c, userId, targetId)
}
