package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Insert(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	args := m.Called(ctx, tokenScope, tokenPlaintext)
	return args.Get(0).(*domain.User), args.Error(1)
}
