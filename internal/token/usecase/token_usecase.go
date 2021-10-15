package usecase

import (
	"context"

	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
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
	const op errors.Op = "tokenUsecase.Insert"
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	err := tu.tr.Insert(ctx, token)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

func (tu *tokenUsecase) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	const op errors.Op = "tokenUsecase.DeleteAllForUser"

	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	err := tu.tr.DeleteAllForUser(ctx, scope, userID)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}
