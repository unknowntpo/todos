package api

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/helpers/validator"
	"github.com/unknowntpo/todos/internal/logger"

	"github.com/julienschmidt/httprouter"

	"github.com/pkg/errors"
)

type taskAPI struct {
	tu     domain.TaskUsecase
	logger logger.Logger
}

type CreateTaskRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

type CreateTaskResponse struct {
	Task *domain.Task `json:"task"`
}

type DeleteTaskByIdResponse struct {
	Message string `json:"message"`
}

type GetAllTasksResponse struct {
	Metadata *domain.Metadata `json:"metadata"`
	Tasks    []*domain.Task   `json:"tasks"`
}

type GetTaskByIdResponse struct {
	Task *domain.Task `json:"task"`
}

type UpdateTaskByIdRequest struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Done    bool   `json:"done,omitempty"`
}

type UpdateTaskByIdResponse struct {
	Task *domain.Task `json:"updated_task"`
}

func NewTaskAPI(router *httprouter.Router, tu domain.TaskUsecase, logger logger.Logger) {
	api := &taskAPI{tu: tu, logger: logger}
	router.HandlerFunc(http.MethodGet, "/v1/tasks", api.GetAll)
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", api.GetByID)
	router.HandlerFunc(http.MethodPost, "/v1/tasks", api.Insert)
	router.HandlerFunc(http.MethodPatch, "/v1/tasks/:id", api.Update)
	router.HandlerFunc(http.MethodDelete, "/v1/tasks/:id", api.Delete)
}

// GetAll gets all tasks.
// TODO: GetAll should get all tasks with specific user id.
// @Summary Gets all tasks for specific user.
// @Description: None.
// @Accept  json
// @Produce  json
// @Param userId query int true "User Id"
// @Param title query string false "title filter"
// @Param sort query string false "sort filter"
// @Param id query string false "id filter"
// @Param page query string false "page filter"
// @Param page_size query string false "page size filter"
// @Success 200 {object} GetAllTasksResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /v1/tasks [get]
func (t *taskAPI) GetAll(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		domain.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = helpers.ReadString(qs, "title", "")
	input.Page = helpers.ReadInt(qs, "page", 1, v)
	input.PageSize = helpers.ReadInt(qs, "page_size", 20, v)

	input.Sort = helpers.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "-id", "-title"}

	if domain.ValidateFilters(v, input.Filters); !v.Valid() {
		helpers.FailedValidationResponse(w, r, v.Errors)
		return
	}

	ctx := r.Context()
	tasks, metadata, err := t.tu.GetAll(ctx, input.Title, input.Filters)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{
		"metadata": metadata,
		"tasks":    tasks,
	}, nil)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
	}
}

// GetByID gets a task by its id.
func (t *taskAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		helpers.NotFoundResponse(w, r)
		return
	}

	ctx := r.Context()
	task, err := t.tu.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			helpers.NotFoundResponse(w, r)
		default:
			t.logger.PrintError(err, nil)
			helpers.ServerErrorResponse(w, r, err)
		}
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"task": task}, nil)
}

// Insert inserts a new task.
func (t *taskAPI) Insert(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"debug": "Insert called"}, nil)
}

// Update updates an exist task.
func (t *taskAPI) Update(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		t.logger.PrintError(err, nil)
		helpers.NotFoundResponse(w, r)
		return
	}

	ctx := r.Context()
	task, err := t.tu.GetByID(ctx, id)
	if err != nil {
		// TODO: errors.Is
		t.logger.PrintError(err, nil)
		helpers.ServerErrorResponse(w, r, err)
	}

	var input struct {
		Title   *string `json:"title"`   // task title
		Content *string `json:"content"` // task content
		Done    *bool   `json:"done"`    // true if task is done
	}
	// readJSON
	err = helpers.ReadJSON(w, r, &input)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
	}

	if input.Title != nil {
		task.Title = *input.Title
	}

	if input.Content != nil {
		task.Content = *input.Content
	}

	if input.Done != nil {
		task.Done = *input.Done
	}

	v := validator.New()

	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if domain.ValidateTask(v, task); !v.Valid() {
		helpers.FailedValidationResponse(w, r, v.Errors)
		return
	}

	ctx = r.Context()
	taskUpdated, err := t.tu.Update(ctx, id, task)
	if err != nil {
		// TODO: errors.Is() to determine which error we got.
		helpers.ServerErrorResponse(w, r, err)
	}

	// write updated JSON to
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"task": taskUpdated}, nil)
}

// Delete delets an exist task.
func (t *taskAPI) Delete(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"debug": "Delete called"}, nil)
}
