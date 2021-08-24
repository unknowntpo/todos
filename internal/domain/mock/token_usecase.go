package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"
)

type MockTokenUsecase struct{}

func NewTokenUsecase() domain.TokenUsecase {
	return &MockTokenUsecase{}
}

func (m *MockTokenUsecase) ValidateTokenPlaintext(ctx context.Context, v domain.Validator, tokenPlaintext string) {
	return
}
