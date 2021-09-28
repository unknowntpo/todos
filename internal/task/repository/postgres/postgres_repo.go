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

func (tr *taskRepo) GetAll(ctx context.Context, userID int64, title string, filters domain.Filters) ([]*domain.Task, domain.Metadata, error) {
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), id, user_id, created_at, title, content, done, version
        FROM tasks
        WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') 
	AND user_id = $2
        ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.SortColumn(), filters.SortDirection())

	args := []interface{}{title, userID, filters.Limit(), filters.Offset()}

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
			&task.UserID,
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.Metadata{}, errors.Errorf("taskRepo.GetAll: %v", domain.ErrRecordNotFound)
		default:
			return nil, domain.Metadata{}, errors.WithMessage(err, "taskRepo.GetAll:")
		}
	}

	metadata := domain.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tasks, metadata, nil
}

func (tr *taskRepo) GetByID(ctx context.Context, userID int64, taskID int64) (*domain.Task, error) {
	if taskID < 1 {
		return nil, domain.ErrRecordNotFound
	}

	query := `
	SELECT id, user_id, created_at, title, content, done, version
	FROM tasks
	WHERE id = $1
	AND user_id = $2`

	var task domain.Task

	err := tr.DB.QueryRowContext(ctx, query, taskID, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.CreatedAt,
		&task.Title,
		&task.Content,
		&task.Done,
		&task.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.Errorf("taskRepo.GetByID: %v", domain.ErrRecordNotFound)
		default:
			return nil, errors.WithMessage(err, "taskRepo.GetByID:")
		}
	}

	return &task, nil
}

func (tr *taskRepo) Insert(ctx context.Context, userID int64, task *domain.Task) error {
	query := `INSERT INTO tasks (user_id, title, content, done)
	      VALUES ($1, $2, $3, $4)
	      RETURNING id, created_at, version`
	args := []interface{}{userID, task.Title, task.Content, task.Done}

	return tr.DB.QueryRowContext(ctx, query, args...).Scan(&task.ID, &task.CreatedAt, &task.Version)
}

func (tr *taskRepo) Update(ctx context.Context, task *domain.Task) error {
	query := `UPDATE tasks
        SET title = $1, content = $2, done = $3, version = version + 1
	WHERE id = $4 AND user_id = $5 AND version = $6
	RETURNING version`

	args := []interface{}{
		task.Title,
		task.Content,
		task.Done,
		task.ID,
		task.UserID,
		task.Version,
	}

	if err := tr.DB.QueryRowContext(ctx, query, args...).Scan(&task.Version); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errors.Errorf("taskRepo.Update: %v", domain.ErrRecordNotFound)
		default:
			return errors.WithMessage(err, "taskRepo.Update:")
		}
	}

	return nil
}

func (tr *taskRepo) Delete(ctx context.Context, userID int64, taskID int64) error {
	query := `DELETE FROM tasks
        WHERE id = $1 AND user_id = $2`

	result, err := tr.DB.ExecContext(ctx, query, taskID, userID)
	if err != nil {
		return errors.WithMessage(err, "taskRepo.Delete:")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.WithMessage(err, "taskRepo.Delete:")
	}

	// If no rows were affected, we know that the tasks table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return errors.Errorf("taskRepo.Delete: %v", domain.ErrRecordNotFound)
	}

	return nil
}
