package usecase

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/unknowntpo/todos/config"
	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
	_repoMock "github.com/unknowntpo/todos/internal/domain/mocks"
	"github.com/unknowntpo/todos/internal/logger"
	"github.com/unknowntpo/todos/internal/logger/zerolog"
	"github.com/unknowntpo/todos/internal/mailer"
	"github.com/unknowntpo/todos/internal/testutil"
	"github.com/unknowntpo/todos/pkg/naivepool"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	userRepo   *_repoMock.UserRepository
	tokenRepo  *_repoMock.TokenRepository
	logBuf     *bytes.Buffer
	logger     logger.Logger
	pool       *naivepool.Pool
	poolCancel context.CancelFunc
	mailer     *mailer.Mailer
	fakeUser   *domain.User
}

func (suite *UserUsecaseTestSuite) SetupSuite() {
	maxJobs := 5
	maxWorkers := 5
	workerChanSize := 5
	suite.pool = naivepool.New(maxJobs, maxWorkers, workerChanSize)
	// Starting worker pool.
	poolCtx, poolCancel := context.WithCancel(context.Background())
	suite.poolCancel = poolCancel

	suite.pool.Start(poolCtx)

}

func (suite *UserUsecaseTestSuite) TearDownSuite() {
	suite.poolCancel()
	suite.pool.Wait()
}

// SetupTest do migration up for each test.
func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.userRepo = new(_repoMock.UserRepository)
	suite.tokenRepo = new(_repoMock.TokenRepository)
	suite.logBuf = new(bytes.Buffer)
	suite.logger = zerolog.New(suite.logBuf)
	suite.mailer = mailer.New(&config.Smtp{})
	suite.fakeUser = testutil.NewFakeUser(suite.T(), "Alice Smith", "alice@example.com", "pa55word", false)

}

// SetupTest do migration down for each test to ensure the results of
// this test won't affect to the result of next test.
func (suite *UserUsecaseTestSuite) TearDownTest() {
	suite.userRepo = nil
	suite.tokenRepo = nil
	suite.logBuf = nil
	suite.logger = nil
	suite.mailer = nil
	suite.fakeUser = nil
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

func (suite *UserUsecaseTestSuite) TestInsert() {
	suite.Run("Success", func() {
		// Manually setup test because suite won't do it for us at each sub-test.
		suite.SetupTest()
		suite.userRepo.On("Insert", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
			return user.Name == suite.fakeUser.Name && user.Email == suite.fakeUser.Email
		})).Return(nil)

		userUsecase := NewUserUsecase(suite.userRepo, suite.tokenRepo, suite.pool, suite.mailer, suite.logger, 3*time.Second)

		ctx := context.TODO()
		err := userUsecase.Insert(ctx, suite.fakeUser)
		suite.NoError(err)

		suite.userRepo.AssertExpectations(suite.T())

		// Manually tear down test because suite won't do it for us at each sub-test.
		suite.TearDownTest()
	})

	suite.Run("Fail with some errors", func() {
		suite.SetupTest()
		dummyErr := errors.New("something goes wrong")
		wantErr := errors.E(errors.Op("mockUserRepo.Insert"), dummyErr)

		suite.userRepo.On("Insert", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
			return user.Name == suite.fakeUser.Name && user.Email == suite.fakeUser.Email
		})).Return(wantErr)

		userUsecase := NewUserUsecase(suite.userRepo, suite.tokenRepo, suite.pool, suite.mailer, suite.logger, 3*time.Second)

		ctx := context.TODO()

		// FIXME: Why err == nil ??
		err := userUsecase.Insert(ctx, suite.fakeUser)

		suite.ErrorIs(err, dummyErr)
		suite.Equal("userUsecase.Insert: mockUserRepo.Insert: something goes wrong", err.Error(), "error message should be equal")

		suite.userRepo.AssertExpectations(suite.T())
		suite.TearDownTest()
	})
}

func (suite *UserUsecaseTestSuite) TestAuthenticate() {
	suite.Run("Success", func() {
		suite.SetupTest()

		token, err := domain.GenerateToken(suite.fakeUser.ID, 30*time.Minute, domain.ScopeActivation)
		if err != nil {
			suite.T().Fatal("fail to generate token")
		}

		// when GetForToken is called with token.Scope and token.Plaintext as arguments,
		// it should return user we defined and nil error.
		suite.userRepo.On("GetForToken", mock.Anything, token.Scope, token.Plaintext).Return(suite.fakeUser, nil)

		userUsecase := NewUserUsecase(suite.userRepo, suite.tokenRepo, suite.pool, suite.mailer, suite.logger, 3*time.Second)

		ctx := context.TODO()
		gotUser, err := userUsecase.Authenticate(ctx, token.Scope, token.Plaintext)
		suite.NoError(err)

		suite.Equal(suite.fakeUser.ID, gotUser.ID, "user ID should be equal")
		suite.Equal(suite.fakeUser.Name, gotUser.Name, "user name should be equal")
		suite.Equal(suite.fakeUser.Email, gotUser.Email, "email should be equal")
		suite.Equal(suite.fakeUser.Password.Hash, gotUser.Password.Hash, "password_hash should be equal")

		suite.userRepo.AssertExpectations(suite.T())

		suite.TearDownTest()
	})

	suite.Run("Fail with some errors", func() {
		suite.SetupTest()

		token, err := domain.GenerateToken(suite.fakeUser.ID, 30*time.Minute, domain.ScopeActivation)
		if err != nil {
			suite.T().Fatal("fail to generate token")
		}

		// Define the error we expect
		dummyErr := errors.New("something goes wrong")
		wantErr := errors.E(errors.Op("mockUserRepo.GetForToken"), dummyErr)

		// when GetForToken is called with token.Scope and token.Plaintext as arguments,
		// it should return user we defined and nil error.
		suite.userRepo.On("GetForToken", mock.Anything, token.Scope, token.Plaintext).Return(nil, wantErr)

		userUsecase := NewUserUsecase(suite.userRepo, suite.tokenRepo, suite.pool, suite.mailer, suite.logger, 3*time.Second)

		ctx := context.TODO()
		gotUser, err := userUsecase.Authenticate(ctx, token.Scope, token.Plaintext)
		suite.Nil(gotUser, "gotUser should be nil due to provided error")

		suite.ErrorIs(err, dummyErr)
		suite.Equal("userUsecase.Authenticate: mockUserRepo.GetForToken: something goes wrong", err.Error(), "error message should be equal")

		suite.userRepo.AssertExpectations(suite.T())

		suite.TearDownTest()
	})
}

func (suite *UserUsecaseTestSuite) TestUpdate() {
	suite.Run("Success", func() {
		suite.SetupTest()

		wantUpdatedName := "Alice Smith Jr."

		// When Update is called, we expect no error is returned
		suite.userRepo.On("Update", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
			return user.Name == suite.fakeUser.Name && user.Email == suite.fakeUser.Email
		})).Return(nil)

		userUsecase := NewUserUsecase(suite.userRepo, suite.tokenRepo, suite.pool, suite.mailer, suite.logger, 3*time.Second)

		ctx := context.TODO()

		// set the new name to fakeUser
		suite.fakeUser.Name = wantUpdatedName

		err := userUsecase.Update(ctx, suite.fakeUser)
		suite.NoError(err)

		// and fakeUser.Name should be updated to "Alice Smith Jr."
		suite.Equalf(wantUpdatedName, suite.fakeUser.Name, "user name should be updated to %s", wantUpdatedName)

		suite.userRepo.AssertExpectations(suite.T())

		suite.TearDownTest()
	})

	suite.Run("Fail with some errors", func() {
		suite.SetupTest()

		wantUpdatedName := "Alice Smith Jr."

		// Define the error we expect.
		dummyErr := errors.New("something goes wrong")
		wantErr := errors.E(errors.Op("mockUserRepo.Update"), dummyErr)

		// When Update is called, we expect no error is returned
		suite.userRepo.On("Update", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
			return user.Name == suite.fakeUser.Name && user.Email == suite.fakeUser.Email
		})).Return(wantErr)

		userUsecase := NewUserUsecase(suite.userRepo, suite.tokenRepo, suite.pool, suite.mailer, suite.logger, 3*time.Second)

		ctx := context.TODO()

		// set the new name to fakeUser
		suite.fakeUser.Name = wantUpdatedName

		err := userUsecase.Update(ctx, suite.fakeUser)

		suite.ErrorIs(err, dummyErr)
		suite.Equal("userUsecase.Update: mockUserRepo.Update: something goes wrong", err.Error(), "error message should be equal")

		suite.userRepo.AssertExpectations(suite.T())

		suite.TearDownTest()
	})
}
