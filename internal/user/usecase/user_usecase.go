package usecase

import (
	"context"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(ur domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{userRepo: ur, contextTimeout: timeout}
}

func (uu *userUsecase) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const op errors.Op = "userUsecase.GetByEmail"

	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return user, nil
}
func (uu *userUsecase) Insert(ctx context.Context, user *domain.User) error {
	const op errors.Op = "userUsecase.Insert"

	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	err := uu.userRepo.Insert(ctx, user)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}
func (uu *userUsecase) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	const op errors.Op = "userUsecase.GetForToken"

	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.userRepo.GetForToken(ctx, tokenScope, tokenPlaintext)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return user, nil
}
func (uu *userUsecase) Update(ctx context.Context, user *domain.User) error {
	const op errors.Op = "userUsecase.Update"

	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	err := uu.userRepo.Update(ctx, user)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}
