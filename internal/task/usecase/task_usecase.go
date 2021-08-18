package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
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

func (tu *taskUsecase) GetAll(ctx context.Context, title string, filters helpers.Filters) ([]*domain.Task, helpers.Metadata, error) {
	return nil, helpers.Metadata{}, nil
}

// Just call repo layer method for now.
func (tu *taskUsecase) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	// TODO: What does it canceled ? The slow function?
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	task, err := tu.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fail to get task by id from task repo: %w", err)
	}
	return task, nil
}
func (tu *taskUsecase) Insert(ctx context.Context, task *domain.Task) error {
	return nil
}
func (tu *taskUsecase) Update(ctx context.Context, id int64, taskUpdated *domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	taskUpdated, err := tu.taskRepo.Update(ctx, id, taskUpdated)
	if err != nil {
		return nil, fmt.Errorf("fail to update task from task usecase: %v", err)
	}
	return taskUpdated, nil
}

func (tu *taskUsecase) Delete(ctx context.Context, id int64) error {
	return nil
}
