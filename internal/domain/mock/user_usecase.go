package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
)

type mockUserUsecase struct{}

func NewUserUsecase() domain.UserUsecase {
	return &mockUserUsecase{}
}

// GetByEmail returns dummy user
func (m *mockUserUsecase) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{
		ID:        1,
		CreatedAt: time.Now(),
		Name:      "Alice Smith",
		Email:     "alice@example.com",
		Activated: true,
		Version:   1,
	}

	err := user.Password.Set("pa55word")
	if err != nil {
		return nil, fmt.Errorf("in mock user usecase,fail to set password for dummy user: %v", err)
	}

	return user, nil
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
