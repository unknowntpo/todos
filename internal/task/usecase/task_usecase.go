package usecase

import (
	"context"
	"time"

	"github.com/unknowntpo/todos/internal/domain"

	"github.com/pkg/errors"
)

type taskUsecase struct {
	taskRepo       domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecase(t domain.TaskRepository, timeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		taskRepo:       t,
		contextTimeout: timeout,
	}
}

func (tu *taskUsecase) GetAll(ctx context.Context, title string, filters domain.Filters) ([]*domain.Task, domain.Metadata, error) {
	return nil, domain.Metadata{}, nil
}

// Just call repo layer method for now.
func (tu *taskUsecase) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	// TODO: What does it canceled ? The slow function?
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	task, err := tu.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.WithMessage(err, "task usecase.GetByID")
	}
	return task, nil
}
func (tu *taskUsecase) Insert(ctx context.Context, task *domain.Task) error {
	// TODO: Implement it.
	return nil
}
func (tu *taskUsecase) Update(ctx context.Context, id int64, taskUpdated *domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	taskUpdated, err := tu.taskRepo.Update(ctx, id, taskUpdated)
	if err != nil {
		return nil, errors.WithMessage(err, "task usecase.Update")

	}
	return taskUpdated, nil
}

func (tu *taskUsecase) Delete(ctx context.Context, id int64) error {
	// TODO: Implement it.
	return nil
}
