package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers/validator"
)

type MockTokenUsecase struct{}

func NewTokenUsecase() domain.TokenUsecase {
	return &MockTokenUsecase{}
}

func (m *MockTokenUsecase) ValidateTokenPlaintext(ctx context.Context, v *validator.Validator, tokenPlaintext string) {
	return
}
