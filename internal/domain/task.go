package domain

import (
	"context"
	"time"
)

// Task represent the data structure of our task object.
type Task struct {
	ID        int64     `json:"id"`      // Unique integer ID for the task
	CreatedAt time.Time `json:"-"`       // Timestamp for when the task is added to our database
	Title     string    `json:"title"`   // task title
	Content   string    `json:"content"` // task content
	Done      bool      `json:"done"`    // true if task is done
	Version   int32     `json:"version"` // The version number starts at 1 and will be incremented each
	// time the task information is updated
}

type TaskUsecase interface {
	GetAll(ctx context.Context, title string, filters Filters) ([]*Task, Metadata, error)
	GetByID(ctx context.Context, id int64) (*Task, error)
	Insert(ctx context.Context, task *Task) error
	Update(ctx context.Context, id int64, task *Task) (*Task, error)
	Delete(ctx context.Context, id int64) error
}

type TaskRepository interface {
	GetAll(ctx context.Context, title string, filters Filters) ([]*Task, Metadata, error)
	GetByID(ctx context.Context, id int64) (*Task, error)
	Insert(ctx context.Context, task *Task) error
	Update(ctx context.Context, id int64, task *Task) (*Task, error)
	Delete(ctx context.Context, id int64) error
}
