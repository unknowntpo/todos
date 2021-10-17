package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockTokenUsecase struct {
	mock.Mock
}

func (m *MockTokenUsecase) Insert(ctx context.Context, token *domain.Token) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenUsecase) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	args := m.Called(ctx, scope, userID)
	return args.Error(0)
}
