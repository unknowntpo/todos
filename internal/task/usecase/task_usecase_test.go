package usecase

import (
	"testing"

	"github.com/unknowntpo/todos/internal/domain/mock"
)

func TestGetByID(t *testing.T) {
	// Test success
	// TODO: Test failure.
	repo := mock.NewTaskRepo()

	repo.On("GetByID", ctx, id).Return(&domain.Task{
		ID:        1,
		CreatedAt: mock.Anything,
		Title:     "",
		Content:   "",
		Done:      false,
		Version:   1,
	}, nil)
}
