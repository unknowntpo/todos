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

type TaskRepoTestSuite struct {
	suite.Suite
	container testcontainers.Container
	db        *sql.DB
	mig       *migrate.Migrate
	fakeuser  *domain.User
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

	// Setup fake user.
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
func (suite *TaskRepoTestSuite) TearDownTest() {
	err := suite.mig.Down()
	if err != nil {
		suite.T().Fatal(err)
	}

	// set it to nil to prevent it from affecting next test.
	suite.fakeuser = nil
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTaskRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepoTestSuite))
}

func (suite *TaskRepoTestSuite) TestGetAll() {
	suite.Run("Success", func() {
		suite.Run("search with title", func() {
			suite.TearDownTest()
			suite.SetupTest()

			repo := NewTaskRepo(suite.db)

			fakeTasks := []*domain.Task{
				{
					Title:   "Do housework with my friend",
					Content: "It's boring!",
					Done:    false,
				},
				{
					Title:   "Learn first principle",
					Content: "It's cool!",
					Done:    true,
				},
			}
			// insert dummy task
			// insert it into db
			ctx := context.TODO()

			for _, task := range fakeTasks {
				if err := repo.Insert(ctx, suite.fakeuser.ID, task); err != nil {
					suite.T().Fatalf("failed to insert dummy task %v to database: %v", task, err)
				}
			}

			// follow the precedure in taskAPI to create a request
			var input struct {
				Title string
				domain.Filters
			}

			input.Title = "housework"
			input.Page = 1
			input.PageSize = 10
			input.Sort = "id"
			input.SortSafelist = []string{"id", "-id", "title", "-title"}

			gotTasks, gotMeta, err := repo.GetAll(ctx, suite.fakeuser.ID, input.Title, input.Filters)
			suite.NoError(err)

			// We expect gotTasks contains only one task "Do housework with my friend".
			wantMeta := domain.Metadata{
				CurrentPage:  1,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 1,
			}

			wantTask := fakeTasks[0]
			suite.Equal(wantMeta, gotMeta, "metadata should be equal")
			suite.Equal(wantTask.Title, gotTasks[0].Title)
			suite.Equal(wantTask.Content, gotTasks[0].Content)
			suite.Equal(wantTask.Done, gotTasks[0].Done)
		})

		suite.Run("search with filter:sort by '-id'", func() {
			suite.TearDownTest()
			suite.SetupTest()

			// test with page = 2, pagesize = 1
			// sort = -id

			repo := NewTaskRepo(suite.db)

			fakeTasks := []*domain.Task{
				{
					Title:   "Do housework with my friend",
					Content: "It's boring!",
					Done:    false,
				},
				{
					Title:   "Learn first principle",
					Content: "It's cool!",
					Done:    true,
				},
			}
			// insert dummy task
			// insert it into db
			ctx := context.TODO()

			for _, task := range fakeTasks {
				if err := repo.Insert(ctx, suite.fakeuser.ID, task); err != nil {
					suite.T().Fatalf("failed to insert dummy task %v to database: %v", task, err)
				}
			}

			// follow the precedure in taskAPI to create a request
			var input struct {
				Title string
				domain.Filters
			}

			input.Page = 2
			input.PageSize = 1
			input.Sort = "-id"
			input.SortSafelist = []string{"id", "-id", "title", "-title"}

			gotTasks, gotMeta, err := repo.GetAll(ctx, suite.fakeuser.ID, input.Title, input.Filters)
			suite.NoError(err)

			// We expect gotTasks contains only one task "Do housework with my friend".
			wantMeta := domain.Metadata{
				CurrentPage:  2,
				PageSize:     1,
				FirstPage:    1,
				LastPage:     2,
				TotalRecords: 2,
			}

			wantTask := fakeTasks[0]

			suite.Equal(wantMeta, gotMeta, "metadata should be equal")
			suite.Equal(wantTask.Title, gotTasks[0].Title)
			suite.Equal(wantTask.Content, gotTasks[0].Content)
			suite.Equal(wantTask.Done, gotTasks[0].Done)
		})
	})
	// FIXME: Maybe using failed on database error to test errors.ErrDatabase ?
	suite.Run("Fail on database error-timeout", func() {
		suite.TearDownTest()
		suite.SetupTest()

		suite.T().Fail()
	})
	suite.Run("Fail on record not found", func() {
		suite.TearDownTest()
		suite.SetupTest()

		suite.T().Fail()
	})
}

func (suite *TaskRepoTestSuite) TestGetByID() {
	suite.Run("Success", func() {
		suite.TearDownTest()
		suite.SetupTest()

		suite.T().Fail()
	})

	// FIXME: Maybe using failed on database error to test errors.ErrDatabase ?
	suite.Run("Fail on database error-timeout", func() {
		suite.TearDownTest()
		suite.SetupTest()

		suite.T().Fail()
	})
	suite.Run("Fail on record not found", func() {
		suite.TearDownTest()
		suite.SetupTest()

		suite.T().Fail()
	})
}

func (suite *TaskRepoTestSuite) TestInsert() {
	suite.Run("Success", func() {
		suite.TearDownTest()
		suite.SetupTest()

		// init repo
		repo := NewTaskRepo(suite.db)
		// prepare a task
		wantTask := &domain.Task{
			Title:   "test insert",
			Content: "Yeah",
			Done:    false,
		}
		// insert it into db
		ctx := context.TODO()
		if err := repo.Insert(ctx, suite.fakeuser.ID, wantTask); err != nil {
			suite.T().Fatalf("failed to insert task %v to database: %v", wantTask, err)
		}

		// call GetByID with title
		gotTask, err := repo.GetByID(ctx, suite.fakeuser.ID, wantTask.ID)
		suite.NoError(err)

		suite.Equal(wantTask.Title, gotTask.Title)
		suite.Equal(wantTask.Content, gotTask.Content)
		suite.Equal(wantTask.Done, gotTask.Done)
	})

	// FIXME: Maybe using failed on database error to test errors.ErrDatabase ?
	suite.Run("Fail on database error-timeout", func() {
		suite.TearDownTest()
		suite.SetupTest()

		// init repo
		repo := NewTaskRepo(suite.db)
		// prepare a task
		wantTask := &domain.Task{
			Title:   "test insert",
			Content: "Yeah",
			Done:    false,
		}
		// insert it into db
		ctx := context.TODO()
		if err := repo.Insert(ctx, suite.fakeuser.ID, wantTask); err != nil {
			suite.T().Fatalf("failed to insert task %v to database: %v", wantTask, err)
		}

		// Prepare a deadline-exceeded context.
		ctx, cancel := context.WithTimeout(context.Background(), -7*time.Minute)
		defer cancel()

		_, err := repo.GetByID(ctx, suite.fakeuser.ID, wantTask.ID)
		suite.ErrorIs(err, context.DeadlineExceeded)
	})
}

func (suite *TaskRepoTestSuite) TestUpdate() {
	suite.Run("Success", func() {
		suite.TearDownTest()
		suite.SetupTest()

		suite.T().Fail()
	})

	// FIXME: Maybe using failed on database error to test errors.ErrDatabase ?
	suite.Run("Fail on database error-timeout", func() {
		suite.TearDownTest()
		suite.SetupTest()

		suite.T().Fail()
	})
	suite.Run("Fail on database error: edit conflict", func() {
		suite.TearDownTest()
		suite.SetupTest()

		suite.T().Fail()
	})
}

func (suite *TaskRepoTestSuite) TestDelete() {
	suite.Run("Success", func() {
		suite.TearDownTest()
		suite.SetupTest()

		suite.T().Fail()
	})

	// FIXME: Maybe using failed on database error to test errors.ErrDatabase ?
	suite.Run("Fail on database error-timeout", func() {
		suite.TearDownTest()
		suite.SetupTest()

		suite.T().Fail()
	})
	suite.Run("Fail on record not found", func() {
		suite.TearDownTest()
		suite.SetupTest()

		suite.T().Fail()
	})
}
