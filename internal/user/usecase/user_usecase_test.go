package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
	_repoMock "github.com/unknowntpo/todos/internal/domain/mocks"
	"github.com/unknowntpo/todos/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsert(t *testing.T) {
	// When success, it should return no error.
	t.Run("Success", func(t *testing.T) {
		userRepo := new(_repoMock.UserRepository)
		// We don't use tokenRepo, just a placeholder.
		tokenRepo := new(_repoMock.TokenRepository)

		fakeUser := testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", false)

		userRepo.On("Insert", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
			return user.Name == "Alice Smith" && user.Email == "alice@example.com"
		})).Return(nil)

		userUsecase := NewUserUsecase(userRepo, tokenRepo, 3*time.Second)

		ctx := context.TODO()
		err := userUsecase.Insert(ctx, fakeUser)
		assert.NoError(t, err)

		userRepo.AssertExpectations(t)
	})

	t.Run("Fail with some errors", func(t *testing.T) {
		userRepo := new(_repoMock.UserRepository)
		tokenRepo := new(_repoMock.TokenRepository)

		fakeUser := testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", false)

		dummyErr := errors.New("something goes wrong")
		wantErr := errors.E(errors.Op("mockUserRepo.Insert"), dummyErr)

		userRepo.On("Insert", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
			return user.Name == "Alice Smith" && user.Email == "alice@example.com"
		})).Return(wantErr)

		userUsecase := NewUserUsecase(userRepo, tokenRepo, 3*time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), -7*time.Minute)
		defer cancel()

		err := userUsecase.Insert(ctx, fakeUser)

		assert.ErrorIs(t, err, dummyErr)
		assert.Equal(t, "userUsecase.Insert: mockUserRepo.Insert: something goes wrong", err.Error(), "error message should be equal")

		userRepo.AssertExpectations(t)
	})
}

func TestAuthenticate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userRepo := new(_repoMock.UserRepository)
		tokenRepo := new(_repoMock.TokenRepository)

		// new fake user
		user := testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", false)

		token, err := domain.GenerateToken(user.ID, 30*time.Minute, domain.ScopeActivation)
		if err != nil {
			t.Fatal("fail to generate token")
		}

		// when GetForToken is called with token.Scope and token.Plaintext as arguments,
		// it should return user we defined and nil error.
		userRepo.On("GetForToken", mock.Anything, token.Scope, token.Plaintext).Return(user, nil)

		userUsecase := NewUserUsecase(userRepo, tokenRepo, 3*time.Second)

		ctx := context.TODO()
		gotUser, err := userUsecase.Authenticate(ctx, token.Scope, token.Plaintext)
		assert.NoError(t, err)

		assert.Equal(t, user.ID, gotUser.ID, "user ID should be equal")
		assert.Equal(t, user.Name, gotUser.Name, "user name should be equal")
		assert.Equal(t, user.Email, gotUser.Email, "email should be equal")
		assert.Equal(t, user.Password.Hash, gotUser.Password.Hash, "password_hash should be equal")

		userRepo.AssertExpectations(t)
	})

	t.Run("Fail with some errors", func(t *testing.T) {
		userRepo := new(_repoMock.UserRepository)
		tokenRepo := new(_repoMock.TokenRepository)

		// new fake user
		user := testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", false)

		token, err := domain.GenerateToken(user.ID, 30*time.Minute, domain.ScopeActivation)
		if err != nil {
			t.Fatal("fail to generate token")
		}

		// Define the error we expect
		dummyErr := errors.New("something goes wrong")
		wantErr := errors.E(errors.Op("mockUserRepo.GetForToken"), dummyErr)

		// when GetForToken is called with token.Scope and token.Plaintext as arguments,
		// it should return user we defined and nil error.
		userRepo.On("GetForToken", mock.Anything, token.Scope, token.Plaintext).Return(nil, wantErr)

		userUsecase := NewUserUsecase(userRepo, tokenRepo, 3*time.Second)

		ctx := context.TODO()
		gotUser, err := userUsecase.Authenticate(ctx, token.Scope, token.Plaintext)
		assert.Nil(t, gotUser, "gotUser should be nil due to provided error")

		assert.ErrorIs(t, err, dummyErr)
		assert.Equal(t, "userUsecase.Authenticate: mockUserRepo.GetForToken: something goes wrong", err.Error(), "error message should be equal")

		userRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userRepo := new(_repoMock.UserRepository)
		tokenRepo := new(_repoMock.TokenRepository)

		fakeUser := testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", false)

		wantUpdatedName := "Alice Smith Jr."

		// When Update is called, we expect no error is returned
		userRepo.On("Update", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
			return user.Name == wantUpdatedName && user.Email == "alice@example.com"
		})).Return(nil)

		userUsecase := NewUserUsecase(userRepo, tokenRepo, 3*time.Second)

		ctx := context.TODO()

		// set the new name to fakeUser
		fakeUser.Name = wantUpdatedName

		err := userUsecase.Update(ctx, fakeUser)
		assert.NoError(t, err)

		// and fakeUser.Name should be updated to "Alice Smith Jr."
		assert.Equalf(t, wantUpdatedName, fakeUser.Name, "user name should be updated to %s", wantUpdatedName)

		userRepo.AssertExpectations(t)
	})

	t.Run("Fail with some errors", func(t *testing.T) {
		userRepo := new(_repoMock.UserRepository)
		tokenRepo := new(_repoMock.TokenRepository)

		fakeUser := testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", false)

		wantUpdatedName := "Alice Smith Jr."

		// Define the error we expect.
		dummyErr := errors.New("something goes wrong")
		wantErr := errors.E(errors.Op("mockUserRepo.Update"), dummyErr)

		// When Update is called, we expect no error is returned
		userRepo.On("Update", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
			return user.Name == wantUpdatedName && user.Email == "alice@example.com"
		})).Return(wantErr)

		userUsecase := NewUserUsecase(userRepo, tokenRepo, 3*time.Second)

		ctx := context.TODO()

		// set the new name to fakeUser
		fakeUser.Name = wantUpdatedName

		err := userUsecase.Update(ctx, fakeUser)

		assert.ErrorIs(t, err, dummyErr)
		assert.Equal(t, "userUsecase.Update: mockUserRepo.Update: something goes wrong", err.Error(), "error message should be equal")

		userRepo.AssertExpectations(t)
	})
}
