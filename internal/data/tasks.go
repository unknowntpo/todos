package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the task struct.
	return t.DB.QueryRowContext(ctx, query, args...).Scan(&task.ID, &task.CreatedAt, &task.Version)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, id).Scan(
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
        WHERE id = $4 AND version = $5
        RETURNING version`

	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		task.Title,
		task.Content,
		task.Done,
		task.ID,
		task.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&task.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// Delete delets a specific record from the tasks table.
func (t TaskModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM tasks
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, we know that the tasks table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// GetAll returns a slice of tasks. Although we're not
// using them right now, we've set this up to accept the various filter parameters as
// arguments.
func (t TaskModel) GetAll(title string, filters Filters) ([]*Task, Metadata, error) {
	// Construct the SQL query to retrieve all task records.
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), id, created_at, title, content, version
        FROM tasks
        WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') 
        ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, filters.limit(), filters.offset()}

	rows, err := t.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	tasks := []*Task{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		var task Task

		err := rows.Scan(
			&totalRecords,
			&task.ID,
			&task.CreatedAt,
			&task.Title,
			&task.Content,
			&task.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Task struct to the slice.
		tasks = append(tasks, &task)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tasks, metadata, nil
}
