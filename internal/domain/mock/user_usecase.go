package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"
)

type mockUserUsecase struct{}

func NewUserUsecase() domain.UserUsecase {
	return &mockUserUsecase{}
}

func (m *mockUserUsecase) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
func (m *mockUserUsecase) Insert(ctx context.Context, user *domain.User) error {
	return nil

}
func (m *mockUserUsecase) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	return nil, nil
}
func (m *mockUserUsecase) Update(ctx context.Context, user *domain.User) error {
	return nil
}
func (m *mockUserUsecase) ValidatePasswordPlaintext(ctx context.Context, v domain.Validator, password string) {
	return
}
func (m *mockUserUsecase) ValidateEmail(v domain.Validator, email string) {
	return
}
func (m *mockUserUsecase) ValidateUser(v domain.Validator, user *domain.User) {
	return
}
