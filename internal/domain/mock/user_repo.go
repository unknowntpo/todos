package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"
)

type mockUserRepo struct{}

func NewUserRepo() domain.UserRepository {
	return &mockUserRepo{}
}

func (m *mockUserRepo) Insert(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *mockUserRepo) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	return nil, nil
}
