package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/testutil"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type RepoTestSuite struct {
	suite.Suite
	container testcontainers.Container
	db        *sql.DB
	mig       *migrate.Migrate
}

func (suite *RepoTestSuite) SetupSuite() {
	ctx := context.Background()

	container, db, err := testutil.CreatePostgresTestContainer(ctx, "testdb")
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.container = container
	suite.db = db

	mig, err := testutil.NewPgMigrator(db)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.mig = mig
}

// TearDownSuite tears down the test suite by closing db connection,
// terminates container.
func (suite *RepoTestSuite) TearDownSuite() {
	defer suite.db.Close()
	ctx := context.Background()
	defer suite.container.Terminate(ctx)
}

// SetupTest do migration up for each test.
func (suite *RepoTestSuite) SetupTest() {
	err := suite.mig.Up()
	if err != nil {
		suite.T().Fatal(err)
	}
}

// SetupTest do migration down for each test to ensure the results of
// this test won't affect to the result of next test.
func (suite *RepoTestSuite) TearDownTest() {
	err := suite.mig.Down()
	if err != nil {
		suite.T().Fatal(err)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}

func (suite *RepoTestSuite) TestInsert() {
	suite.Run("Success", func() {
		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "Alice Smith", "alice@example.com", "pa55word", true)

		ctx := context.TODO()
		err := repo.Insert(ctx, user)
		suite.NoError(err, "error should be nil")
		// Check if the user is inserted.
		query := `SELECT id, name, email, password_hash FROM users
		WHERE id = $1`

		var gotUser domain.User
		err = suite.db.QueryRowContext(ctx, query, user.ID).Scan(&gotUser.ID, &gotUser.Name, &gotUser.Email, &gotUser.Password.Hash)
		suite.NoError(err)

		suite.Equal(user.ID, gotUser.ID, "user ID should be equal")
		suite.Equal(user.Name, gotUser.Name, "user name should be equal")
		suite.Equal(user.Email, gotUser.Email, "email should be equal")
		suite.Equal(user.Password.Hash, gotUser.Password.Hash, "password_hash should be equal")
	})
	suite.Run("Fail on timeout", func() {
		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "Alice Smith", "alice@example.com", "pa55word", true)

		// Create a deadline-exceeded context
		ctx, cancel := context.WithTimeout(context.Background(), -7*time.Minute)
		defer cancel()

		err := repo.Insert(ctx, user)
		suite.ErrorIs(err, context.DeadlineExceeded)
	})
}

func (suite *RepoTestSuite) TestGetByEmail() {
	// TODO: Implement tests.
	suite.T().Fail()
}

func (suite *RepoTestSuite) TestUpdate() {
	// TODO: Implement tests.
	suite.T().Fail()
}

func (suite *RepoTestSuite) TestGetForToken() {
	// TODO: Implement tests.
	suite.T().Fail()
}
