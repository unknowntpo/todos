package domain

import (
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
	GetByID(id int64) (*Task, error)
}
