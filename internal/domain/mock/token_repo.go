package mock

import (
	"context"
	"time"

	"github.com/unknowntpo/todos/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockTokenRepo struct {
	mock.Mock
}

func (m *MockTokenRepo) New(ctx context.Context, userID int64, ttl time.Duration, scope string) (*domain.Token, error) {
	args := m.Called(ctx, userID, ttl, scope)
	return args.Get(0).(*domain.Token), args.Error(1)
}

func (m *MockTokenRepo) Insert(ctx context.Context, token *domain.Token) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepo) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	args := m.Called(ctx, scope, userID)
	return args.Error(0)
}
