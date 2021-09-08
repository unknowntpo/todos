package usecase

import (
	"context"
	"fmt"
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

func (tu *tokenUsecase) Insert(ctx context.Context, token *domain.Token) error {
	if token == nil {
		return domain.ErrNilObject
	}
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	err := tu.tr.Insert(ctx, token)
	if err != nil {
		// TODO: Improve error message chain.
		return fmt.Errorf("failed to insert token at token usecase: %w", err)
	}

	return nil
}

func (tu *tokenUsecase) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	err := tu.tr.DeleteAllForUser(ctx, scope, userID)
	if err != nil {
		// TODO: Improve error message chain.
		return fmt.Errorf("failed to delete all token for user %d at token usecase: %w", userID, err)
	}

	return nil
}
