package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/unknowntpo/todos/internal/domain"

	"github.com/pkg/errors"
)

type taskRepo struct {
	DB *sql.DB
}

func NewTaskRepo(DB *sql.DB) domain.TaskRepository {
	return &taskRepo{DB}
}

func (tr *taskRepo) GetAll(ctx context.Context, title string, filters domain.Filters) ([]*domain.Task, domain.Metadata, error) {
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), id, created_at, title, content, done, version
        FROM tasks
        WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') 
        ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3`, filters.SortColumn(), filters.SortDirection())

	args := []interface{}{title, filters.Limit(), filters.Offset()}

	rows, err := tr.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, domain.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	tasks := []*domain.Task{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		var task domain.Task

		err := rows.Scan(
			&totalRecords,
			&task.ID,
			&task.CreatedAt,
			&task.Title,
			&task.Content,
			&task.Done,
			&task.Version,
		)
		if err != nil {
			return nil, domain.Metadata{}, err
		}

		// Add the Task struct to the slice.
		tasks = append(tasks, &task)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, domain.Metadata{}, err
	}

	metadata := domain.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tasks, metadata, nil
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
