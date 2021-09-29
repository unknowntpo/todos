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

	user := testutil.NewFakeUser(suite.T(), "Alice Smith", "alice@example.com", "pa55word", true)
	// Create fake user
	query := `
	INSERT INTO users (name, email, password_hash, activated)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
	}

	ctx := context.TODO()
	err = suite.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		suite.T().Fatal(err)
	}

	// store fake user in suite
	suite.fakeuser = user
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
	suite.Run("Success", func() {
		// create activation token
		actToken, err := domain.GenerateToken(suite.fakeuser.ID, 30*time.Minute, domain.ScopeActivation)
		if err != nil {
			suite.T().Fatal("fail to generate activation token")
		}

		// create authentication token
		authToken, err := domain.GenerateToken(suite.fakeuser.ID, 30*time.Minute, domain.ScopeAuthentication)
		if err != nil {
			suite.T().Fatal("fail to generate authentication token")
		}

		// create new repo
		repo := NewTokenRepo(suite.db)

		// Insert both activation token and authentication token to repo
		ctx := context.TODO()
		err = repo.Insert(ctx, actToken)
		suite.NoError(err)
		err = repo.Insert(ctx, authToken)
		suite.NoError(err)

		// invoke DeleteAllForUser for activation
		err = repo.DeleteAllForUser(ctx, domain.ScopeActivation, suite.fakeuser.ID)
		suite.NoError(err)

		// assert activation token is deleted.
		var gotActToken domain.Token
		query := `SELECT hash, user_id, expiry, scope FROM tokens WHERE user_id = $1 AND scope = $2`
		err = suite.db.QueryRowContext(ctx, query, suite.fakeuser.ID, domain.ScopeActivation).Scan(&gotActToken.Hash, &gotActToken.UserID, &gotActToken.Expiry, &gotActToken.Scope)
		suite.ErrorIs(err, sql.ErrNoRows)

		// assert authentication token is still there.
		var gotAuthToken domain.Token
		query = `SELECT hash, user_id, expiry, scope FROM tokens WHERE user_id = $1 AND scope = $2`
		err = suite.db.QueryRowContext(ctx, query, suite.fakeuser.ID, domain.ScopeAuthentication).Scan(&gotAuthToken.Hash, &gotAuthToken.UserID, &gotAuthToken.Expiry, &gotAuthToken.Scope)
		suite.NoError(err)

		suite.Equal(authToken.Hash, gotAuthToken.Hash, "hash should be equal")
		suite.Equal(authToken.UserID, gotAuthToken.UserID, "user id should be equal")
		suite.Equal(authToken.Scope, gotAuthToken.Scope, "scope should be equal")
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
