package postgres

import (
	"context"
	"database/sql"
	"testing"
	//	"time"

	//	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/testutil"

	"github.com/golang-migrate/migrate/v4"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type TaskRepoTestSuite struct {
	suite.Suite
	container testcontainers.Container
	db        *sql.DB
	mig       *migrate.Migrate
}

func (suite *TaskRepoTestSuite) SetupSuite() {
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
func (suite *TaskRepoTestSuite) TearDownSuite() {
	defer suite.db.Close()
	ctx := context.Background()
	defer suite.container.Terminate(ctx)
}

// SetupTest do migration up for each test.
func (suite *TaskRepoTestSuite) SetupTest() {
	err := suite.mig.Up()
	if err != nil {
		suite.T().Fatal(err)
	}
}

// SetupTest do migration down for each test to ensure the results of
// this test won't affect to the result of next test.
func (suite *TaskRepoTestSuite) TearDownTest() {
	err := suite.mig.Down()
	if err != nil {
		suite.T().Fatal(err)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTaskRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepoTestSuite))
}

func (suite *TaskRepoTestSuite) TestGetAll() {
	suite.T().Fail()
}

func (suite *TaskRepoTestSuite) TestGetByID() {
	suite.T().Fail()
}

func (suite *TaskRepoTestSuite) TestInsert() {
	suite.T().Fail()
}

func (suite *TaskRepoTestSuite) TestUpdate() {
	suite.T().Fail()
}

func (suite *TaskRepoTestSuite) TestDelete() {
	suite.T().Fail()
}
