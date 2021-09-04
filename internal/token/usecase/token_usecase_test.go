package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/mock"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := new(mock.MockTokenRepo)
		ctx := context.TODO()
		token, err := domain.GenerateToken(1, 30*time.Minute, domain.ScopeActivation)
		assert.NoError(t, err)

		repo.On("Insert", ctx, token).Return(nil)

		err = repo.Insert(ctx, token)
		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("Timeout", func(t *testing.T) {
		repo := new(mock.MockTokenRepo)
		ctx := context.TODO()
		token, err := domain.GenerateToken(1, 30*time.Minute, domain.ScopeActivation)
		assert.NoError(t, err)

		// Simulate the deadline exceeded context.
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(-7*time.Minute))
		defer cancel()

		wantErr := fmt.Errorf("failed to insert token at token usecase: context deadline exceeded")
		repo.On("Insert", ctx, token).Return(wantErr)

		gotErr := repo.Insert(ctx, token)
		assert.Equal(t, wantErr, gotErr)

		repo.AssertExpectations(t)
	})

}

func TestDeleteAllForUser(t *testing.T) {
	t.Fail()
}
