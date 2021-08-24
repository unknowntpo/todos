package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"
)

type mockUserUsecase struct{}

func NewUserUsecase() domain.UserUsecase {
	return &mockUserUsecase{}
}

func (m *mockUserUsecase) ValidatePasswordPlaintext(ctx context.Context, v domain.Validator, password string) {
	return
}
