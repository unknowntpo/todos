package mock

import (
	"context"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
)

type MockTokenUsecase struct{}

func NewTokenUsecase() domain.TokenUsecase {
	return &MockTokenUsecase{}
}

func (m *MockTokenUsecase) New(ctx context.Context, userID int64, ttl time.Duration, scope string) (*domain.Token, error) {
	return nil, nil
}
func (m *MockTokenUsecase) Insert(ctx context.Context, token *domain.Token) error {
	return nil
}
func (m *MockTokenUsecase) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	return nil
}

func (m *MockTokenUsecase) ValidateTokenPlaintext(ctx context.Context, v domain.Validator, tokenPlaintext string) {
	return
}
