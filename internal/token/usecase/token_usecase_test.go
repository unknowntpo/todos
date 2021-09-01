package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/mock"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Fail()
}

func TestInsert(t *testing.T) {
	// Test success
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

	// TODO: Test failure.
}

func TestDeleteAllForUser(t *testing.T) {
	t.Fail()
}
