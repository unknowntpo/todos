package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockTaskRepo struct {
	mock.Mock
}

func NewTaskRepo() domain.TaskRepository {
	return &MockTaskRepo{}
}

func (m *MockTaskRepo) GetAll(ctx context.Context, userID int64, title string, filters domain.Filters) ([]*domain.Task, domain.Metadata, error) {
	args := m.Called(ctx, userID, title, filters)
	return args.Get(0).([]*domain.Task), args.Get(1).(domain.Metadata), args.Error(2)
}

func (m *MockTaskRepo) GetByID(ctx context.Context, userID int64, taskID int64) (*domain.Task, error) {
	args := m.Called(ctx, userID, taskID)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepo) Insert(ctx context.Context, userID int64, task *domain.Task) error {
	args := m.Called(ctx, userID, task)
	return args.Error(0)
}

func (m *MockTaskRepo) Update(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepo) Delete(ctx context.Context, userID int64, taskID int64) error {
	args := m.Called(ctx, userID, taskID)
	return args.Error(0)
}
