package mock

import (
	"context"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	t.Run("test id == 1", func(t *testing.T) {
		mock_repo := NewTaskRepo()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		task, err := mock_repo.GetByID(ctx, 1)
		if assert.NoError(t, err) {
			assert.Equal(t, task, &domain.Task{
				ID:      1,
				Title:   "Do homework",
				Content: "Interesting",
				Done:    true,
				Version: 1,
			})
		}
	})

	t.Run("test record not found", func(t *testing.T) {
	})
}

func TestInsert(t *testing.T) {
	// just insert task and assert no errors.

}

func TestUpdate(t *testing.T) {
	t.Run("test successfully update the task", func(t *testing.T) {

	})

	t.Run("test ErrRecordNotFound", func(t *testing.T) {

	})

}

func TestDelete(t *testing.T) {
	t.Run("test successfully delete the task", func(t *testing.T) {

	})

	t.Run("test ErrRecordNotFound", func(t *testing.T) {

	})

}
