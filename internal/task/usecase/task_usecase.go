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
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	tasks, metadata, err := tu.taskRepo.GetAll(ctx, title, filters)
	if err != nil {
		return nil, domain.Metadata{}, errors.WithMessage(err, "task usecase.GetAll")
	}
	return tasks, metadata, nil
}

// Just call repo layer method for now.
func (tu *taskUsecase) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	task, err := tu.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.WithMessage(err, "task usecase.GetByID")
	}
	return task, nil
}
func (tu *taskUsecase) Insert(ctx context.Context, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	err := tu.taskRepo.Insert(ctx, task)
	if err != nil {
		return errors.WithMessage(err, "task usecase.Insert")
	}
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
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	err := tu.taskRepo.Delete(ctx, id)
	if err != nil {
		return errors.WithMessage(err, "task usecase.Delete")
	}
	return nil
}
