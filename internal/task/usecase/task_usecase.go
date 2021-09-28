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

func (tu *taskUsecase) GetAll(ctx context.Context, userID int64, title string, filters domain.Filters) ([]*domain.Task, domain.Metadata, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	tasks, metadata, err := tu.taskRepo.GetAll(ctx, userID, title, filters)
	if err != nil {
		return nil, domain.Metadata{}, errors.WithMessage(err, "taskUsecase.GetAll")
	}
	return tasks, metadata, nil
}

// Just call repo layer method for now.
func (tu *taskUsecase) GetByID(ctx context.Context, userID int64, taskID int64) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	task, err := tu.taskRepo.GetByID(ctx, userID, taskID)
	if err != nil {
		return nil, errors.WithMessage(err, "taskUsecase.GetByID")
	}
	return task, nil
}
func (tu *taskUsecase) Insert(ctx context.Context, userID int64, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	err := tu.taskRepo.Insert(ctx, userID, task)
	if err != nil {
		return errors.WithMessage(err, "taskUsecase.Insert")
	}
	return nil
}
func (tu *taskUsecase) Update(ctx context.Context, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	err := tu.taskRepo.Update(ctx, task)
	if err != nil {
		return errors.WithMessage(err, "taskUsecase.Update")
	}
	return nil
}

func (tu *taskUsecase) Delete(ctx context.Context, userID int64, taskID int64) error {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	err := tu.taskRepo.Delete(ctx, userID, taskID)
	if err != nil {
		return errors.WithMessage(err, "taskUsecase.Delete")
	}
	return nil
}
