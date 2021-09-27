package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"
)

type mockUserRepo struct{}

func NewUserRepo() domain.UserRepository {
	return &mockUserRepo{}
}

func (m *mockUserRepo) Insert(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepo) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	args := m.Called(ctx, tokenScope, tokenPlaintext)
	return args.Get(0).(*domain.User), args.Error(1)
}
