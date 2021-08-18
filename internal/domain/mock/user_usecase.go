package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers/validator"
)

type mockUserUsecase struct{}

func NewUserUsecase() domain.UserUsecase {
	return &mockUserUsecase{}
}

func (m *mockUserUsecase) ValidatePasswordPlaintext(ctx context.Context, v *validator.Validator, password string) {
	return
}
