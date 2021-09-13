package usecase

import (
	"context"

	"time"

	"github.com/unknowntpo/todos/internal/domain"

	"github.com/pkg/errors"
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
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	err := tu.tr.Insert(ctx, token)
	if err != nil {
		// TODO: Improve error message chain.
		return errors.WithMessage(err, "token usecase.insert")
	}

	return nil
}

func (tu *tokenUsecase) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	err := tu.tr.DeleteAllForUser(ctx, scope, userID)
	if err != nil {
		// TODO: Improve error message chain.
		return errors.WithMessagef(err, "token usecase.deleteallforuser.userID = %d", userID)
	}

	return nil
}
