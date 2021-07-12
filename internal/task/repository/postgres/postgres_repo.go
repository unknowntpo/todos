package postgres

import (
	"database/sql"

	"github.com/unknowntpo/todos/internal/domain"
)

type taskRepo struct {
	DB *sql.DB
}

func NewTaskRepo(DB *sql.DB) domain.TaskRepository {
	return &taskRepo{DB}
}

func (tr *taskRepo) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	return nil, nil
}

func (tr *taskRepo) Insert(ctx context.Context, task *domain.Task) error {
	return nil
}
func (tr *taskRepo) Update(ctx context.Context, id int64, task *domain.Task) (*domain.Task, error) {
	return nil, nil
}

func (tr *taskRepo) Delete(ctx context.Context, id int64) error {
	return nil, nil
}
