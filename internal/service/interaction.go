package service

import (
	"context"
	"time"

	"github.com/LXJ0000/clean-backend/internal/domain"
)

type interactionService struct {
	repo           domain.InteractionRepository
	contextTimeout time.Duration
}

func NewInteractionService(repo domain.InteractionRepository, timeout time.Duration) domain.InteractionService {
	return &interactionService{repo: repo, contextTimeout: timeout}
}

func (uc *interactionService) IncrReadCount(c context.Context, biz string, id int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.IncrReadCount(ctx, biz, id)
}

func (uc *interactionService) Like(c context.Context, biz string, bizID, userID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Like(ctx, biz, bizID, userID)
}

func (uc *interactionService) CancelLike(c context.Context, biz string, bizID, userID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.CancelLike(ctx, biz, bizID, userID)
}

func (uc *interactionService) Stat(c context.Context, biz string, bizID, userID int64) (domain.Interaction, domain.UserInteractionStat, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Stat(ctx, biz, bizID, userID)
}

func (uc *interactionService) Collect(c context.Context, biz string, bizID, userID, collectionID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Collect(ctx, biz, bizID, userID, collectionID)
}

func (uc *interactionService) CancelCollect(c context.Context, biz string, bizID, userID, collectionID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.CancelCollect(ctx, biz, bizID, userID, collectionID)
}

func (uc *interactionService) GetByIDs(c context.Context, biz string, bizIDs []int64) (map[int64]domain.Interaction, error) {
	return nil, nil
}
