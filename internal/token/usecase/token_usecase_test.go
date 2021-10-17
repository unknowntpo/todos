package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
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

	// Make sure usecase.Insert() will return error propagated from token repo layer.
	// with some error message 'token usecase.insert' annotated.
	t.Run("Fail with some errors", func(t *testing.T) {
		const op errors.Op = "tokenUsecase.Insert"

		repo := new(_repoMock.MockTokenRepo)
		token, err := domain.GenerateToken(1, 30*time.Minute, domain.ScopeActivation)
		assert.NoError(t, err)

		ctx := context.TODO()
		// set expectations on mock repo
		dummyErr := errors.E(errors.New("error in mock token repo"))
		repo.On("Insert", mock.Anything, mock.Anything).Return(dummyErr)

		tokenUsecase := NewTokenUsecase(repo, 3*time.Second)
		gotErr := tokenUsecase.Insert(ctx, token)

		wantErr := errors.E(op, dummyErr)
		assert.Equal(t, wantErr, gotErr)

		repo.AssertExpectations(t)
	})
}

func TestDeleteAllForUser(t *testing.T) {
	// When success, it should return no error.
	t.Run("Success", func(t *testing.T) {
		// init token usecase with mock token repo
		repo := new(_repoMock.MockTokenRepo)
		// Set expectations
		repo.On("DeleteAllForUser", mock.Anything, domain.ScopeActivation, mock.MatchedBy(func(userID int64) bool {
			return userID == 1
		})).Return(nil)

		tokenUsecase := NewTokenUsecase(repo, 3*time.Second)

		ctx := context.TODO()
		gotErr := tokenUsecase.DeleteAllForUser(ctx, domain.ScopeActivation, 1)
		assert.NoError(t, gotErr)

		repo.AssertExpectations(t)
	})

	t.Run("Fail with some errors", func(t *testing.T) {
		const op errors.Op = "tokenUsecase.DeleteAllForUser"

		repo := new(_repoMock.MockTokenRepo)

		ctx := context.TODO()
		dummyUserID := int64(1)
		dummyErr := errors.E(errors.New("error in mock token repo"))
		repo.On("DeleteAllForUser", mock.Anything, domain.ScopeActivation, mock.MatchedBy(func(userID int64) bool {
			return userID == dummyUserID
		})).Return(dummyErr)

		tokenUsecase := NewTokenUsecase(repo, 3*time.Second)
		gotErr := tokenUsecase.DeleteAllForUser(ctx, domain.ScopeActivation, dummyUserID)

		wantErr := errors.E(op, dummyErr)
		assert.Equal(t, wantErr, gotErr)

		repo.AssertExpectations(t)
	})
}
