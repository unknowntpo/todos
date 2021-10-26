package api

import (
	"fmt"
	"net/http"

	"github.com/unknowntpo/todos/internal/domain"

	"github.com/unknowntpo/todos/internal/domain/errors"
	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/middleware"
	"github.com/unknowntpo/todos/internal/reactor"

	"github.com/unknowntpo/todos/pkg/validator"

	"github.com/julienschmidt/httprouter"
)

type CreateTaskRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

type CreateTaskResponse struct {
	Task *domain.Task `json:"task"`
}

type DeleteTaskByIDResponse struct {
	Message string `json:"message"`
}

type GetAllTasksResponse struct {
	Metadata *domain.Metadata `json:"metadata"`
	Tasks    []*domain.Task   `json:"tasks"`
}

type GetTaskByIDResponse struct {
	Task *domain.Task `json:"task"`
}

type UpdateTaskByIDRequest struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Done    bool   `json:"done,omitempty"`
}

type UpdateTaskByIDResponse struct {
	Task *domain.Task `json:"updated_task"`
}

type taskAPI struct {
	tu  domain.TaskUsecase
	mid *middleware.Middleware
	rc  *reactor.Reactor
}

func NewTaskAPI(router *httprouter.Router, tu domain.TaskUsecase, mid *middleware.Middleware, rc *reactor.Reactor) {
	api := &taskAPI{tu: tu, mid: mid, rc: rc}
	router.Handler(http.MethodGet, "/v1/tasks", mid.RequireAuthenticatedUser(rc.HandlerWrapper(api.GetAll)))
	router.Handler(http.MethodGet, "/v1/tasks/:id", mid.RequireAuthenticatedUser(rc.HandlerWrapper(api.GetByID)))
	router.Handler(http.MethodPost, "/v1/tasks", mid.RequireAuthenticatedUser(rc.HandlerWrapper(api.Insert)))
	router.Handler(http.MethodPatch, "/v1/tasks/:id", mid.RequireAuthenticatedUser(rc.HandlerWrapper(api.Update)))
	router.Handler(http.MethodDelete, "/v1/tasks/:id", mid.RequireAuthenticatedUser(rc.HandlerWrapper(api.Delete)))
}

// GetAll gets all tasks.
// TODO: GetAll should get all tasks for specific user.
// @Summary Get all tasks for specific user.
// @Description: None.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param title query string false "title filter"
// @Param sort query string false "sort filter"
// @Param id query string false "id filter"
// @Param page query string false "page filter"
// @Param page_size query string false "page size filter"
// @Success 200 {object} GetAllTasksResponse
// @Failure 400 {object} reactor.ErrorResponse
// @Failure 404 {object} reactor.ErrorResponse
// @Failure 500 {object} reactor.ErrorResponse
// @Router /v1/tasks [get]
func (t *taskAPI) GetAll(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "taskAPI.GetAll"

	user := helpers.ContextGetUser(r)

	var input struct {
		Title string
		domain.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.Title = reactor.ReadString(qs, "title", "")
	input.CurrentPage = reactor.ReadInt(qs, "page", 1, v)
	input.PageSize = reactor.ReadInt(qs, "page_size", 20, v)

	input.Sort = reactor.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "-id", "-title"}

	if domain.ValidateFilters(v, input.Filters); !v.Valid() {
		return reactor.FailedValidationResponse(w, r, v.Err())
	}

	ctx := r.Context()
	tasks, metadata, err := t.tu.GetAll(ctx, user.ID, input.Title, input.Filters)
	if err != nil {
		return errors.E(op, err)
	}

	err = reactor.WriteJSON(w, http.StatusOK, &GetAllTasksResponse{
		Metadata: &metadata,
		Tasks:    tasks,
	})
	if err != nil {
		// TODO: tweak the message shown in error chain.
		return errors.E(op, err)
	}

	return nil
}

// GetByID gets a task by its ID.
// @Summary Get task by ID for specific user.
// @Description: None.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param taskID path int true "Task ID"
// @Success 200 {object} GetAllTasksResponse
// @Failure 400 {object} reactor.ErrorResponse
// @Failure 404 {object} reactor.ErrorResponse
// @Failure 500 {object} reactor.ErrorResponse
// @Router /v1/tasks/{taskID} [get]
func (t *taskAPI) GetByID(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "GetByID"
	user := helpers.ContextGetUser(r)

	id, err := reactor.ReadIDParam(r)
	if err != nil {
		return reactor.NotFoundResponse(w, r)
	}

	ctx := r.Context()
	task, err := t.tu.GetByID(ctx, user.ID, id)
	if err != nil {
		switch {
		case errors.KindIs(err, errors.ErrRecordNotFound):
			return reactor.NotFoundResponse(w, r)
		default:
			return errors.E(op, err)
		}
	}

	return reactor.WriteJSON(w, http.StatusOK, GetTaskByIDResponse{task})
}

// Insert inserts a new task.
// @Summary Create a new task for specific user.
// @Description: None.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param userID query int true "User ID"
// @Param reqBody body CreateTaskRequest true "create task request body"
// @Success 201 {object} domain.Task
// @Failure 400 {object} reactor.ErrorResponse
// @Failure 404 {object} reactor.ErrorResponse
// @Failure 500 {object} reactor.ErrorResponse
// @Router /v1/tasks [post]
func (t *taskAPI) Insert(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "taksAPI.Insert"

	var input CreateTaskRequest
	user := helpers.ContextGetUser(r)

	err := reactor.ReadJSON(w, r, &input)
	if err != nil {
		return reactor.BadRequestResponse(w, r, err)
	}

	task := &domain.Task{
		Title:   input.Title,
		Content: input.Content,
		Done:    input.Done,
	}

	// validate the request
	v := validator.New()

	if domain.ValidateTask(v, task); !v.Valid() {
		return reactor.FailedValidationResponse(w, r, v.Err())
	}

	ctx := r.Context()
	err = t.tu.Insert(ctx, user.ID, task)
	if err != nil {
		return errors.E(op, err)
	}

	w.Header().Set("Location", fmt.Sprintf("/v1/tasks/%d", task.ID))

	// Write a JSON response with a 201 Created status code, the task data in the
	// response body, and the Location header.
	err = reactor.WriteJSON(w, http.StatusCreated, &CreateTaskResponse{task})
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// Update updates an exist task for specific user.
// @Summary Update task for specific user.
// @Description: None.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param taskID path int true "Task ID"
// @Param reqBody body UpdateTaskByIDRequest true "request body"
// @Success 200 {object} domain.Task
// @Failure 400 {object} reactor.ErrorResponse
// @Failure 404 {object} reactor.ErrorResponse
// @Failure 500 {object} reactor.ErrorResponse
// @Router /v1/tasks/{taskID} [patch]
func (t *taskAPI) Update(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "taskAPI.Update"

	user := helpers.ContextGetUser(r)

	taskID, err := reactor.ReadIDParam(r)
	if err != nil {
		return reactor.NotFoundResponse(w, r)
	}

	ctx := r.Context()
	task, err := t.tu.GetByID(ctx, user.ID, taskID)
	if err != nil {
		return errors.E(op, err)
	}

	var input struct {
		Title   *string `json:"title"`   // task title
		Content *string `json:"content"` // task content
		Done    *bool   `json:"done"`    // true if task is done
	}

	err = reactor.ReadJSON(w, r, &input)
	if err != nil {
		return errors.E(op, err)
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

	if domain.ValidateTask(v, task); !v.Valid() {
		return reactor.FailedValidationResponse(w, r, v.Err())
	}

	ctx = r.Context()
	err = t.tu.Update(ctx, task)
	if err != nil {
		return errors.E(op, err)
	}

	return reactor.WriteJSON(w, http.StatusOK, &UpdateTaskByIDResponse{task})
}

// Delete delets an exist task.
// @Summary Delete task for specific user.
// @Description: None.
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param userID query int true "User ID"
// @Param taskID path int true "Task ID"
// @Success 200 {object} DeleteTaskByIDResponse
// @Failure 400 {object} reactor.ErrorResponse
// @Failure 404 {object} reactor.ErrorResponse
// @Failure 500 {object} reactor.ErrorResponse
// @Router /v1/tasks/{taskID} [delete]
func (t *taskAPI) Delete(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "taskAPI.Delete"

	user := helpers.ContextGetUser(r)

	// Extract the task ID from the URL.
	taskID, err := reactor.ReadIDParam(r)
	if err != nil {
		return reactor.NotFoundResponse(w, r)
	}

	ctx := r.Context()
	err = t.tu.Delete(ctx, user.ID, taskID)
	if err != nil {
		switch {
		case errors.KindIs(err, errors.ErrRecordNotFound):
			return reactor.NotFoundResponse(w, r)
		default:
			return errors.E(op, err)
		}
	}

	return reactor.WriteJSON(w, http.StatusOK, &DeleteTaskByIDResponse{"task successfully deleted"})
}
