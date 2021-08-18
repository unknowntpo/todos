package mock

import (
	"context"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
)

type MockTokenRepo struct{}

func NewTokenRepo() domain.TokenRepository {
	return &MockTokenRepo{}
}

func (m *MockTokenRepo) New(ctx context.Context, userID int64, ttl time.Duration, scope string) (*domain.Token, error) {
	return nil, nil
}

func (m *MockTokenRepo) Insert(ctx context.Context, token *domain.Token) error {
	return nil
}

func (m *MockTokenRepo) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	return nil
}
