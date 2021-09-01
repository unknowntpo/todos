package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(ur domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{userRepo: ur, contextTimeout: timeout}
}

func (uu *userUsecase) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// TODO: Improve error chain message
		return nil, fmt.Errorf("fail to get user by email at user usecase: %w", err)
	}
	return user, nil
}
func (uu *userUsecase) Insert(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	err := uu.userRepo.Insert(ctx, user)
	if err != nil {
		// TODO: Improve error chain message
		return fmt.Errorf("fail to insert user at user usecase: %w", err)
	}
	return nil
}
func (uu *userUsecase) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.userRepo.GetForToken(ctx, tokenScope, tokenPlaintext)
	if err != nil {
		// TODO: Improve error chain message
		return nil, fmt.Errorf("fail to get user for given token at user usecase: %w", err)
	}
	return user, nil
}
func (uu *userUsecase) Update(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	err := uu.userRepo.Update(ctx, user)
	if err != nil {
		// TODO: Improve error chain message
		return fmt.Errorf("fail to update user at user usecase: %w", err)
	}
	return nil
}
