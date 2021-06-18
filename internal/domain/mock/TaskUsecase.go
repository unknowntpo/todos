package mock

import (
	"github.com/unknowntpo/todos/internal/domain"
)

type MockTaskUsecase struct{}

func NewTaskUsecase() *MockTaskUsecase {
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

func (m *MockTaskUsecase) GetByID(id int64) (*domain.Task, error) {
	task := &domain.Task{
		ID:      1,
		Title:   "Do homework",
		Content: "Interesting",
		Done:    true,
		Version: 1,
	}

	return task, nil
}
