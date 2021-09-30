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
		suite.TearDownTest()
		suite.SetupTest()

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
		suite.TearDownTest()
		suite.SetupTest()

		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "Ben Johnson", "ben@example.com", "pa55word", true)

		// Create a deadline-exceeded context
		ctx, cancel := context.WithTimeout(context.Background(), -7*time.Minute)
		defer cancel()

		err := repo.Insert(ctx, user)
		suite.ErrorIs(err, context.DeadlineExceeded)
	})
	suite.Run("Fail on duplicate email", func() {
		suite.TearDownTest()
		suite.SetupTest()

		// create two user with same email
		repo := NewUserRepo(suite.db)
		userBen := testutil.NewFakeUser(suite.T(), "Ben Johnson", "ben@example.com", "pa55word", true)
		userAlan := testutil.NewFakeUser(suite.T(), "Alan Johnson", "ben@example.com", "pa55word", true)

		// insert user Ben
		ctx := context.TODO()
		if err := repo.Insert(ctx, userBen); err != nil {
			suite.T().Fatalf("failed to insert user %q into database", userBen.Name)
		}

		err := repo.Insert(ctx, userAlan)
		suite.ErrorIs(err, domain.ErrDuplicateEmail)
	})

}

func (suite *RepoTestSuite) TestGetByEmail() {
	suite.Run("Success", func() {
		suite.TearDownTest()
		suite.SetupTest()

		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "Alice Smith", "alice@example.com", "pa55word", true)

		// insert user into testdb
		ctx := context.TODO()
		if err := repo.Insert(ctx, user); err != nil {
			// means our implementation of Insert method has some bug !
			suite.T().Fatalf("failed to insert user into database: %v", err)
		}

		// call GetByEmail
		gotUser, err := repo.GetByEmail(ctx, user.Email)
		suite.NoError(err)
		// assert user we got is want we want.
		suite.Equal(user.ID, gotUser.ID, "user ID should be equal")
		suite.Equal(user.Name, gotUser.Name, "user name should be equal")
		suite.Equal(user.Email, gotUser.Email, "email should be equal")
		suite.Equal(user.Password.Hash, gotUser.Password.Hash, "password_hash should be equal")
	})
	suite.Run("Fail on timeout", func() {
		suite.TearDownTest()
		suite.SetupTest()

		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "Ben Johnson", "ben@example.com", "pa55word", true)

		// insert user into testdb
		ctx := context.TODO()
		if err := repo.Insert(ctx, user); err != nil {
			// means our implementation of Insert method has some bug !
			suite.T().Fatalf("failed to insert user into database: %v", err)
		}

		// Create a deadline-exceeded context
		ctx, cancel := context.WithTimeout(context.Background(), -7*time.Minute)
		defer cancel()

		_, err := repo.GetByEmail(ctx, user.Email)
		suite.ErrorIs(err, context.DeadlineExceeded)
	})
	suite.Run("Fail on record not found", func() {
		suite.TearDownTest()
		suite.SetupTest()

		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "Ben Johnson", "ben@example.com", "pa55word", true)

		// insert user into testdb
		ctx := context.TODO()
		_, err := repo.GetByEmail(ctx, user.Email)
		suite.ErrorIs(err, domain.ErrRecordNotFound)
	})
}

func (suite *RepoTestSuite) TestUpdate() {
	suite.Run("Success", func() {
		suite.TearDownTest()
		suite.SetupTest()

		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "Alice Smith", "alice@example.com", "pa55word", true)

		// Insert new user to db
		ctx := context.TODO()
		if err := repo.Insert(ctx, user); err != nil {
			suite.T().Fatalf("failed to insert user into database: %v", err)
		}

		// Update user email and check if user email is updated.
		oldVersion := user.Version
		newName := "Alice Smith Jr."
		user.Name = newName
		err := repo.Update(ctx, user)
		suite.NoError(err)

		// check if user name is updated
		var updatedUser domain.User
		query := `SELECT id, name, version FROM users
		WHERE id = $1`
		err = suite.db.QueryRowContext(ctx, query, user.ID).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Version)
		suite.NoError(err)

		suite.Equal(updatedUser.ID, user.ID, "ID should be equal")
		suite.Equal(updatedUser.Name, newName, "user name should be equal")
		suite.Equal(oldVersion+1, user.Version)
	})
	suite.Run("Fail on timeout", func() {
		suite.TearDownTest()
		suite.SetupTest()

		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "Kevin Smith", "kevin@example.com", "pa55word", true)

		// Insert new user to db
		ctx := context.TODO()
		if err := repo.Insert(ctx, user); err != nil {
			suite.T().Fatalf("failed to insert user into database: %v", err)
		}

		// Update user email and check if error is deadline exceeded.
		newName := "Kevin Smith Jr."
		user.Name = newName
		ctx, cancel := context.WithTimeout(context.Background(), -7*time.Minute)
		defer cancel()

		err := repo.Update(ctx, user)
		suite.ErrorIs(err, context.DeadlineExceeded)
	})
	suite.Run("Fail on edit conflict", func() {
		suite.TearDownTest()
		suite.SetupTest()

		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "Brian Smith", "brian@example.com", "pa55word", true)

		// Insert new user to db
		ctx := context.TODO()
		if err := repo.Insert(ctx, user); err != nil {
			suite.T().Fatalf("failed to insert user into database: %v", err)
		}

		// update user with wrong userID
		user.ID = 333

		err := repo.Update(ctx, user)
		suite.ErrorIs(err, domain.ErrEditConflict)
	})
	suite.Run("Fail on duplicate email", func() {
		suite.TearDownTest()
		suite.SetupTest()

		repo := NewUserRepo(suite.db)
		userAlan := testutil.NewFakeUser(suite.T(), "Alan Smith", "alan@example.com", "pa55word", true)

		userBrian := testutil.NewFakeUser(suite.T(), "Brian Smith", "brian@example.com", "pa55word", true)

		// Insert new user to db
		ctx := context.TODO()
		if err := repo.Insert(ctx, userAlan); err != nil {
			suite.T().Fatalf("failed to insert user into database: %v", err)
		}
		if err := repo.Insert(ctx, userBrian); err != nil {
			suite.T().Fatalf("failed to insert user into database: %v", err)
		}

		// update user with existed email
		userAlan.Email = "brian@example.com"

		err := repo.Update(ctx, userAlan)
		suite.ErrorIs(err, domain.ErrDuplicateEmail)
	})
}

func (suite *RepoTestSuite) TestGetForToken() {
	suite.Run("Success", func() {
		suite.TearDownTest()
		suite.SetupTest()

		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "Alice Smith", "alice@example.com", "pa55word", true)

		// insert user into testdb
		ctx := context.TODO()
		if err := repo.Insert(ctx, user); err != nil {
			// means our implementation of Insert method has some bug !
			suite.T().Fatalf("failed to insert user into database: %v", err)
		}

		// create new token
		token, err := domain.GenerateToken(user.ID, 30*time.Minute, domain.ScopeActivation)
		if err != nil {
			suite.T().Fatalf("failed to generate token: %v", err)
		}

		// insert into token table
		query := `
		INSERT INTO tokens (hash, user_id, expiry, scope) 
		VALUES ($1, $2, $3, $4)`

		args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}
		res, err := suite.db.ExecContext(ctx, query, args...)
		suite.NoError(err)
		rowsAffected, err := res.RowsAffected()
		suite.NoError(err)
		if rowsAffected != int64(1) {
			suite.T().Fatalf("failed to insert token into database: got rowsAffected = %d, want %d", rowsAffected, 1)
		}

		// call GetForToken()
		gotUser, err := repo.GetForToken(ctx, token.Scope, token.Plaintext)
		suite.NoError(err)

		// assert user we got is want we want.
		suite.Equal(user.ID, gotUser.ID, "user ID should be equal")
		suite.Equal(user.Name, gotUser.Name, "user name should be equal")
		suite.Equal(user.Email, gotUser.Email, "email should be equal")
		suite.Equal(user.Password.Hash, gotUser.Password.Hash, "password_hash should be equal")
	})

	suite.Run("Fail on timeout", func() {
		suite.TearDownTest()
		suite.SetupTest()

		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "John Smith", "john@example.com", "pa55word", true)

		// insert user into testdb
		ctx := context.TODO()
		if err := repo.Insert(ctx, user); err != nil {
			// means our implementation of Insert method has some bug !
			suite.T().Fatalf("failed to insert user into database: %v", err)
		}

		// create new token
		token, err := domain.GenerateToken(user.ID, 30*time.Minute, domain.ScopeActivation)
		if err != nil {
			suite.T().Fatalf("failed to generate token: %v", err)
		}

		// insert into token table
		query := `
		INSERT INTO tokens (hash, user_id, expiry, scope) 
		VALUES ($1, $2, $3, $4)`

		args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}
		res, err := suite.db.ExecContext(ctx, query, args...)
		suite.NoError(err)
		rowsAffected, err := res.RowsAffected()
		suite.NoError(err)
		if rowsAffected != int64(1) {
			suite.T().Fatalf("failed to insert token into database: got rowsAffected = %d, want %d", rowsAffected, 1)
		}

		// call GetForToken()
		ctx, cancel := context.WithTimeout(context.Background(), -7*time.Minute)
		defer cancel()

		_, err = repo.GetForToken(ctx, token.Scope, token.Plaintext)
		suite.ErrorIs(err, context.DeadlineExceeded)
	})
	suite.Run("Fail on record not found", func() {
		suite.TearDownTest()
		suite.SetupTest()

		repo := NewUserRepo(suite.db)
		user := testutil.NewFakeUser(suite.T(), "John Smith", "john@example.com", "pa55word", true)

		token, err := domain.GenerateToken(user.ID, 30*time.Minute, domain.ScopeActivation)
		if err != nil {
			suite.T().Fatalf("failed to generate token: %v", err)
		}

		ctx := context.TODO()
		_, err = repo.GetForToken(ctx, token.Scope, token.Plaintext)
		suite.ErrorIs(err, domain.ErrRecordNotFound)
	})
}
