package mock

import (
	"context"
	"sync"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
)

type MockTaskRepo struct {
	taskList []*domain.Task
	id       int64
	mu       sync.Mutex
}

func NewTaskRepo() domain.TaskRepository {
	return &MockTaskRepo{
		taskList: []*domain.Task{
			{ID: 1, Title: "Do homework", Content: "Interesting", Done: true, Version: 1},
		},
		id: 1,
	}
}

func (m *MockTaskRepo) GetAll(ctx context.Context, title string, filters helpers.Filters) ([]*domain.Task, helpers.Metadata, error) {
	metadata := helpers.Metadata{
		CurrentPage:  1,
		PageSize:     10,
		FirstPage:    1,
		LastPage:     1,
		TotalRecords: 2,
	}

	return nil, metadata, nil
}

func (m *MockTaskRepo) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	for _, task := range m.taskList {
		if task.ID == id {
			return task, nil
		}
	}
	return nil, domain.ErrRecordNotFound
}

func (m *MockTaskRepo) appendTask(t *domain.Task) {
	// set id, created time, version at here.
	t.ID = m.id
	t.CreatedAt = time.Now()
	t.Version = 1

	m.mu.Lock()
	m.taskList = append(m.taskList, t)
	m.id++
	m.mu.Unlock()
}

func (m *MockTaskRepo) Insert(ctx context.Context, task *domain.Task) error {
	_ = ctx

	newTask := &domain.Task{
		Title:   task.Title,
		Content: task.Content,
		Done:    task.Done,
	}

	m.appendTask(newTask)

	return nil
}

func (m *MockTaskRepo) Update(ctx context.Context, id int64, task *domain.Task) (*domain.Task, error) {
	_ = ctx
	// check if task is in m.taskList
	for i := range m.taskList {
		if m.taskList[i].ID == id {
			// update the task
			m.taskList[i].ID = task.ID
			m.taskList[i].Title = task.Title
			m.taskList[i].Content = task.Content
			m.taskList[i].Done = task.Done
			m.taskList[i].Version++

			return m.taskList[i], nil
		}
	}

	return nil, domain.ErrRecordNotFound
}

func (m *MockTaskRepo) Delete(ctx context.Context, id int64) error {
	return nil
}
