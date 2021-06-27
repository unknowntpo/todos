package mock

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"
)

type MockTaskUsecase struct{}

func NewTaskUsecase() domain.TaskUsecase {
	return &MockTaskUsecase{}
}

/*
func (m *mockTaskUsecase) GetAll(title string, filters Filters) ([]*Task, Metadata, error) {

	tasks := []*domain.Task{
		{ID: 1, Title: "Do homework", Content: "Interesting", Done: true, Version: 1},
		{ID: 2, Title: "Wash dishes", Content: "Boring", Done: false, Version: 1},
	}

	metadata := mockMetadata{
		CurrentPage:  1,
		PageSize:     10,
		FirstPage:    1,
		LastPage:     1,
		TotalRecords: 2,
	}
	return tasks, Metadata, nil
}

*/

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

func (m *MockTaskUsecase) Update(ctx context.Context, id int64, task *domain.Task) (*domain.Task, error) {
	// Placeholder for context.
	_ = ctx

	// TODO: What fileds should we update ?
	task = &domain.Task{
		ID:      1,
		Title:   "Do homework",
		Content: "boring",
		Done:    false,
		Version: 2,
	}
	return task, nil
}
