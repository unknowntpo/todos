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
	// FIXME: Just a stub! Need actual implementation
	user := &domain.User{
		ID:        1,
		CreatedAt: time.Now(),
		Name:      "Alice Smith",
		Email:     "alice@example.com",
		Activated: true,
		Version:   1,
	}

	err := user.Password.Set("pa55word")
	if err != nil {
		return nil, fmt.Errorf("in user usecase,fail to set password for dummy user: %v", err)
	}

	return user, nil
}
func (uu *userUsecase) Insert(ctx context.Context, user *domain.User) error {
	return nil
}
func (uu *userUsecase) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	return nil, nil
}
func (uu *userUsecase) Update(ctx context.Context, user *domain.User) error {
	return nil
}
