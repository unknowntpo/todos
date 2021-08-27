package usecase

import (
	"context"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
)

type tokenUsecase struct {
	tr             domain.TokenRepository
	contextTimeout time.Duration
}

func NewTokenUsecase(tr domain.TokenRepository, timeout time.Duration) domain.TokenUsecase {
	return &tokenUsecase{
		tr:             tr,
		contextTimeout: timeout,
	}
}

func (tu *tokenUsecase) New(ctx context.Context, userID int64, ttl time.Duration, scope string) (*domain.Token, error) {
	return nil, nil
}
func (tu *tokenUsecase) Insert(ctx context.Context, token *domain.Token) error {
	return nil
}
func (tu *tokenUsecase) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	return nil
}
