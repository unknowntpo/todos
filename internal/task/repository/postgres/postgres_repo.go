package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/unknowntpo/todos/internal/domain"
)

type taskRepo struct {
	DB *sql.DB
}

func NewTaskRepo(DB *sql.DB) domain.TaskRepository {
	return &taskRepo{DB}
}

func (tr *taskRepo) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	if id < 1 {
		return nil, domain.ErrRecordNotFound
	}

	query := `
	SELECT id, created_at, title, content, done, version
	FROM tasks
	WHERE id = $1`

	var task domain.Task

	err := tr.DB.QueryRowContext(ctx, query, id).Scan(
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
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &task, nil
}

func (tr *taskRepo) Insert(ctx context.Context, task *domain.Task) error {
	return nil
}
func (tr *taskRepo) Update(ctx context.Context, id int64, task *domain.Task) (*domain.Task, error) {
	return nil, nil
}

func (tr *taskRepo) Delete(ctx context.Context, id int64) error {
	return nil
}
