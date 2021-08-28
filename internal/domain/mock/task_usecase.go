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

func (m *MockTaskUsecase) GetAll(ctx context.Context, title string, filters domain.Filters) ([]*domain.Task, domain.Metadata, error) {
	args := m.Called(ctx, title, filters)
	return args.Get(0).([]*domain.Task), args.Get(1).(domain.Metadata), args.Error(2)
}

// GetByID gets the task when id is 1, otherwise it returns nil, and an ErrRecordNotFound error.
func (m *MockTaskUsecase) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Task), args.Error(1)
}

// Insert inserts the task with id = 2
func (m *MockTaskUsecase) Insert(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

// Update updates the task where id = 1, change the following fields:
// Content: boring, Done: true, and increment Version field to 2.
func (m *MockTaskUsecase) Update(ctx context.Context, id int64, task *domain.Task) (*domain.Task, error) {
	args := m.Called(ctx, id, task)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
