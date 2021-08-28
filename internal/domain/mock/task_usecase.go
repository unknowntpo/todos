package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
)

type MockTaskUsecase struct{}

func NewTaskUsecase() domain.TaskUsecase {
	return &MockTaskUsecase{}
}

func (m *MockTaskUsecase) GetAll(ctx context.Context, title string, filters domain.Filters) ([]*domain.Task, domain.Metadata, error) {
	return nil, &helpers.Metadata{}, nil
}

// GetByID gets the task when id is 1, otherwise it returns nil, and an ErrRecordNotFound error.
func (m *MockTaskUsecase) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	// Placeholder for context.
	_ = ctx
	if id != 1 {
		return nil, domain.ErrRecordNotFound
	}
	task := &domain.Task{
		ID:      1,
		Title:   "Do homework",
		Content: "Interesting",
		Done:    true,
		Version: 1,
	}

	return task, nil
}

// Insert inserts the task with id = 2
func (m *MockTaskUsecase) Insert(ctx context.Context, task *domain.Task) error {
	// Placeholder for context.
	_ = ctx
	_ = task

	return nil
}

// Update updates the task where id = 1, change the following fields:
// Content: boring, Done: true, and increment Version field to 2.
func (m *MockTaskUsecase) Update(ctx context.Context, id int64, task *domain.Task) (*domain.Task, error) {
	// Placeholder for context.
	_ = ctx
	_ = task

	task = &domain.Task{
		ID:      1,
		Title:   "Do homework",
		Content: "boring",
		Done:    true,
		Version: 2,
	}
	return task, nil
}

func (m *MockTaskUsecase) Delete(ctx context.Context, id int64) error {
	return nil
}
