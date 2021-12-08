package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
	_repoMock "github.com/unknowntpo/todos/internal/domain/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// init mock taskrepo
		repo := new(_repoMock.TaskRepository)

		fakeUserID := int64(1)

		fakeTasks := []*domain.Task{
			{
				UserID:  int64(1),
				Title:   "Do housework with my friend",
				Content: "It's boring!",
				Done:    false,
			},
			{
				UserID:  int64(1),
				Title:   "Learn first principle",
				Content: "It's cool!",
				Done:    true,
			},
		}

		// follow the precedure in taskAPI to create a request
		var input struct {
			Title string
			domain.Filters
		}

		input.Title = "housework"
		input.CurrentPage = 1
		input.PageSize = 10
		input.Sort = "id"
		input.SortSafelist = []string{"id", "-id", "title", "-title"}

		// We expect gotTasks contains only one task "Do housework with my friend".
		wantMeta := domain.Metadata{
			CurrentPage:  1,
			PageSize:     10,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 1,
		}

		wantTasks := fakeTasks[:1]
		repo.On("GetAll", mock.Anything, fakeUserID, input.Title, input.Filters).
			Return(wantTasks, wantMeta, nil)

		taskUserUsecase := NewTaskUsecase(repo, 3*time.Second)

		ctx := context.TODO()
		gotTasks, gotMeta, err := taskUserUsecase.GetAll(ctx, fakeUserID, input.Title, input.Filters)
		assert.NoError(t, err)

		assert.Equal(t, wantMeta, gotMeta, "metadata should be equal")
		assert.Equal(t, len(wantTasks), len(gotTasks), "length of gotTasks should be equal to wantTasks")
		assert.Equal(t, wantTasks[0].UserID, gotTasks[0].UserID, ".UserID should be equal")
		assert.Equal(t, wantTasks[0].Title, gotTasks[0].Title)
		assert.Equal(t, wantTasks[0].Content, gotTasks[0].Content)
		assert.Equal(t, wantTasks[0].Done, gotTasks[0].Done)

		repo.AssertExpectations(t)
	})

	t.Run("Fail with some errors", func(t *testing.T) {
		// init mock taskrepo
		repo := new(_repoMock.TaskRepository)

		fakeUserID := int64(1)

		// follow the precedure in taskAPI to create a request
		var input struct {
			Title string
			domain.Filters
		}

		input.Title = "housework"
		input.CurrentPage = 1
		input.PageSize = 10
		input.Sort = "id"
		input.SortSafelist = []string{"id", "-id", "title", "-title"}

		// Set expected error.
		dummyErr := errors.New("something goes wrong")
		wantErr := errors.E(errors.Op("mockTaskRepo.GetAll"), dummyErr)

		var wantTasks []*domain.Task = nil
		wantMeta := domain.CalculateMetadata(0, 0, 0)

		repo.On("GetAll", mock.Anything, fakeUserID, input.Title, input.Filters).
			Return(wantTasks, wantMeta, wantErr)

		taskUsecase := NewTaskUsecase(repo, 3*time.Second)

		ctx := context.TODO()
		gotTasks, gotMeta, err := taskUsecase.GetAll(ctx, fakeUserID, input.Title, input.Filters)
		assert.Nil(t, gotTasks)
		assert.Equal(t, wantMeta, gotMeta)
		assert.ErrorIs(t, err, dummyErr)
		assert.Equal(t, "taskUsecase.GetAll: >> mockTaskRepo.GetAll: >> something goes wrong", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	t.Fail()
}

func TestInsert(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	t.Fail()
}

func TestUpdate(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	t.Fail()
}

func TestDelete(t *testing.T) {
	t.Skip("TODO: finish the implementation")
	t.Fail()
}
