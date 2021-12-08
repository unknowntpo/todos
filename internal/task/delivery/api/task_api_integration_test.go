package api

// setup
// taskrepo
// taskusecase
// taskdelivery
import (
	"bytes"
	"context"
	"database/sql"
	"testing"

	"github.com/unknowntpo/todos/internal/domain"
	//"github.com/unknowntpo/todos/internal/domain/errors"
	//"github.com/unknowntpo/todos/internal/domain/mocks"
	"github.com/unknowntpo/todos/internal/logger"
	"github.com/unknowntpo/todos/internal/middleware"
	"github.com/unknowntpo/todos/internal/reactor"
	"github.com/unknowntpo/todos/internal/testutil"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type TaskAPIIntegrationTestSuite struct {
	suite.Suite
	container testcontainers.Container
	db        *sql.DB
	mig       *migrate.Migrate

	taskUsecase domain.TaskUsecase
	userUsecase domain.UserUsecase
	logBuf      *bytes.Buffer
	logger      logger.Logger
	mid         *middleware.Middleware
	rc          *reactor.Reactor
	fakeUser    *domain.User
}

func (suite *TaskAPIIntegrationTestSuite) SetupSuite() {
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

func (suite *TaskAPIIntegrationTestSuite) TearDownSuite() {}

func (suite *TaskAPIIntegrationTestSuite) SetupTest() {
	/*
		suite.taskUsecase = new(mocks.TaskUsecase)
		suite.userUsecase = new(mocks.UserUsecase)
		suite.logBuf = new(bytes.Buffer)
		suite.logger = zerolog.New(logbuf)
		suite.rc = reactor.NewReactor(suite.logger)
		suite.mid = middleware.New(&config.Config{}, suite.userUsecase, suite.rc)
		suite.fakeUser = testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", true)
	*/
	/*
		taskRepo := _taskRepoPostgres.NewTaskRepo(app.database)
		taskUsecase := _taskUsecase.NewTaskUsecase(taskRepo, 3*time.Second)
		genMid := _generalMiddleware.New(app.config, userUsecase, rc)
		_taskAPI.NewTaskAPI(router, taskUsecase, genMid, rc)
	*/
}

// SetupTest do migration down for each test to ensure the results of
// this test won't affect to the result of next test.
func (suite *TaskAPIIntegrationTestSuite) TearDownTest() {
	suite.taskUsecase = nil
	suite.userUsecase = nil
	suite.logBuf = nil
	suite.logger = nil
	suite.mid = nil
	suite.rc = nil
	suite.fakeUser = nil
}

func TestTaskAPIIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests...")
	}

	suite.Run(t, new(TaskAPIIntegrationTestSuite))
}
