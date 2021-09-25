package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockTaskUsecase struct {
	mock.Mock
}

func NewTaskUsecase() domain.TaskUsecase {
	return &MockTaskUsecase{}
}

func (m *MockTaskUsecase) GetAll(ctx context.Context, userID int64, title string, filters domain.Filters) ([]*domain.Task, domain.Metadata, error) {
	args := m.Called(ctx, userID, title, filters)
	return args.Get(0).([]*domain.Task), args.Get(1).(domain.Metadata), args.Error(2)
}

func (m *MockTaskUsecase) GetByID(ctx context.Context, userID int64, taskID int64) (*domain.Task, error) {
	args := m.Called(ctx, userID, taskID)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) Insert(ctx context.Context, userID int64, task *domain.Task) error {
	args := m.Called(ctx, userID, task)
	return args.Error(0)
}

func (m *MockTaskUsecase) Update(ctx context.Context, userID int64, taskID int64, task *domain.Task) (*domain.Task, error) {
	args := m.Called(ctx, userID, taskID, task)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) Delete(ctx context.Context, userID int64, taskID int64) error {
	args := m.Called(ctx, userID, taskID)
	return args.Error(0)
}
