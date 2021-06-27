package postgres

import (
	"database/sql"
	"fmt"
)

type taskRepo struct {
	DB *sql.DB
}

func NewTaskRepo(DB *sql.DB) domain.TaskRepository {
	return &taskRepo{DB}
}

func (tr *taskRepo) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
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

	var task domain.Task

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
func (tr *taskRepo) Update(ctx content.Context, id int64, task *Task) (*Task, error)
