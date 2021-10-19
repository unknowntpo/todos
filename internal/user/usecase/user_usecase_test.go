package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
	_repoMock "github.com/unknowntpo/todos/internal/domain/mock"
	"github.com/unknowntpo/todos/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsert(t *testing.T) {
	// When success, it should return no error.
	t.Run("Success", func(t *testing.T) {
		repo := new(_repoMock.UserUsecase)

		fakeUser := testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", false)

		repo.On("Insert", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
			return user.Name == "Alice Smith" && user.Email == "alice@example.com"
		})).Return(nil)

		userUsecase := NewUserUsecase(repo, 3*time.Second)

		ctx := context.TODO()
		err := userUsecase.Insert(ctx, fakeUser)
		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("Fail with some errors", func(t *testing.T) {
		repo := new(_repoMock.UserUsecase)

		fakeUser := testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", false)

		dummyErr := errors.New("something goes wrong")
		wantErr := errors.E(errors.Op("mockUserRepo.Insert"), dummyErr)

		repo.On("Insert", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
			return user.Name == "Alice Smith" && user.Email == "alice@example.com"
		})).Return(wantErr)

		userUsecase := NewUserUsecase(repo, 3*time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), -7*time.Minute)
		defer cancel()

		err := userUsecase.Insert(ctx, fakeUser)

		assert.ErrorIs(t, err, dummyErr)
		assert.Equal(t, "userUsecase.Insert: mockUserRepo.Insert: something goes wrong", err.Error(), "error message should be equal")

		repo.AssertExpectations(t)
	})

}
