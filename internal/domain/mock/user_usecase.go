package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
)

type mockUserUsecase struct{}

func NewUserUsecase() domain.UserUsecase {
	return &mockUserUsecase{}
}

// GetByEmail returns dummy user
func (m *mockUserUsecase) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserUsecase) Insert(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *mockUserUsecase) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	args := m.Called(ctx, tokenScope, tokenPlaintext)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserUsecase) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
