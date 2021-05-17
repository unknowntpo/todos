package data

import (
	"database/sql"
	"errors"
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

// TaskModel wraps a sql.DB connection pool.
type TaskModel struct {
	DB *sql.DB
}

// ValidateTask check if task match the constrains.
func ValidateTask(v *validator.Validator, task *Task) {
	v.Check(task.Title != "", "title", "must be provided")
	v.Check(len(task.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(task.Content != "", "content", "must be provided")
	v.Check(len(task.Content) <= 500, "title", "must not be more than 500 bytes long")
}

// Insert accepts a pointer to a task struct, which should contain the
// data for the new record.
func (t TaskModel) Insert(task *Task) error {
	query := `
        INSERT INTO tasks (title, content, done) 
        VALUES ($1, $2, $3)
        RETURNING id, created_at, version`

	// Create an args slice containing the values for the placeholder parameters from
	// the task struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{task.Title, task.Content, task.Done}

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the task struct.
	return t.DB.QueryRow(query, args...).Scan(&task.ID, &task.CreatedAt, &task.Version)
}

// Add a placeholder method for fetching a specific record from the tasks table.
func (t TaskModel) Get(id int64) (*Task, error) {
	// The PostgreSQL bigserial type that we're using for the task ID starts
	// auto-incrementing at 1 by default, so we know that no tasks will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, title, content, done, version
        FROM tasks
        WHERE id = $1`

	var task Task

	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := t.DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.Title,
		&task.Content,
		&task.Done,
		&task.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &task, nil
}

func (t TaskModel) Update(task *Task) error {
	query := `
        UPDATE tasks 
        SET title = $1, content = $2, done = $3, version = version + 1
        WHERE id = $4
        RETURNING version`

	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		task.Title,
		task.Content,
		task.Done,
		task.ID,
	}

	return t.DB.QueryRow(query, args...).Scan(&task.Version)
}

// Add a placeholder method for deleting a specific record from the tasks table.
func (t TaskModel) Delete(id int64) error {
	return nil
}
