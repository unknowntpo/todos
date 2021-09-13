package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	_repoMock "github.com/unknowntpo/todos/internal/domain/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsert(t *testing.T) {
	// When success, it should return no error.
	t.Run("Success", func(t *testing.T) {
		repo := new(_repoMock.MockTokenRepo)
		token, err := domain.GenerateToken(1, 30*time.Minute, domain.ScopeActivation)
		assert.NoError(t, err)

		ctx := context.TODO()
		repo.On("Insert", mock.Anything, token).Return(nil)

		tokenUsecase := NewTokenUsecase(repo, 3*time.Second)
		err = tokenUsecase.Insert(ctx, token)
		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	// When failed on db timeout, it should return with following errors.
	t.Run("Fail on db timeout", func(t *testing.T) {
		t.Fail()
	})
}

func TestDeleteAllForUser(t *testing.T) {
	t.Fail()
}
