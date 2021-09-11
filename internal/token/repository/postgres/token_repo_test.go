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
	fakeuser  *domain.User
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

	// Create fake user
	query := `
	INSERT INTO users (name, email, password_hash, activated)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	args := []interface{}{
		"Alice Smith",
		"alice@example.com",
		`\x24326124313224765a682f57676e6246446c4f696757654d62616849756874642f3363526b75576558655469782e374f71524c372e78746570635479`,
		true,
	}

	// Insert the fake user into user table.
	var user domain.User

	ctx := context.TODO()
	err = suite.db.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		suite.T().Fatal(err)
	}

	// store fake user in suite
	suite.fakeuser = &user
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
		wantToken, err := domain.GenerateToken(suite.fakeuser.ID, 30*time.Minute, domain.ScopeActivation)
		if err != nil {
			suite.T().Fatal("fail to generate token")
		}

		ctx := context.TODO()
		repo := NewTokenRepo(suite.db)
		err = repo.Insert(ctx, wantToken)
		suite.NoError(err)

		// do db query to get inserted token
		var gotToken domain.Token
		query := `SELECT hash, user_id, expiry, scope FROM tokens WHERE user_id = $1`
		err = suite.db.QueryRowContext(ctx, query, suite.fakeuser.ID).Scan(&gotToken.Hash, &gotToken.UserID, &gotToken.Expiry, &gotToken.Scope)
		suite.NoError(err)

		// We don't care about Expiry, we just make sure that hash, userid,
		// scope are the same, them we can say that these two token are equal.
		suite.Equal(wantToken.Hash, gotToken.Hash, "should be equal")
		suite.Equal(wantToken.UserID, gotToken.UserID, "should be equal")
		suite.Equal(wantToken.Scope, gotToken.Scope, "should be equal")
	})

	suite.Run("Fail on timeout", func() {
		// Create a deadline-exceeded context
		ctx, cancel := context.WithTimeout(context.Background(), -7*time.Minute)
		defer cancel()
		token, err := domain.GenerateToken(suite.fakeuser.ID, 30*time.Minute, domain.ScopeActivation)
		if err != nil {
			suite.T().Fatal("fail to generate token")
		}

		repo := NewTokenRepo(suite.db)
		err = repo.Insert(ctx, token)
		suite.ErrorIs(err, context.DeadlineExceeded)
	})
}

func (suite *RepoTestSuite) TestDeleteAllForUser() {
	// TODO: Implement tests.
	suite.T().Fail()
}
