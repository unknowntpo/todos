package data

import (
	"time"

	"github.com/unknowntpo/todos/internal/validator"
)

// Task represent the data structure of our task object.
type Task struct {
	ID        int64     `json:"id"` // Unique integer ID for the task
	CreatedAt time.Time `json:"-"`  // Timestamp for when the task is added to our database
	// hide this from clients using `-` directive
	Title   string `json:"title"`   // task title
	Content string `json:"content"` // task content
	Done    bool   `json:"done"`    // true if task is done
	Version int32  `json:"version"` // The version number starts at 1 and will be incremented each
	// time the task information is updated
}

// ValidateTask check if task match the constrains.
func ValidateTask(v *validator.Validator, task *Task) {
	v.Check(task.Title != "", "title", "must be provided")
	v.Check(len(task.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(task.Content != "", "content", "must be provided")
	v.Check(len(task.Content) <= 500, "title", "must not be more than 500 bytes long")
}
