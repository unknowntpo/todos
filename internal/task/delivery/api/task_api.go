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
func (t *taskAPI) GetAll(c *reactor.Context) error {
	const op errors.Op = "taskAPI.GetAll"

	user := helpers.ContextGetUser(c.GetRequest())

	var input struct {
		Title string
		domain.Filters
	}

	v := validator.New()

	input.Title = c.ReadString("title", "")
	input.CurrentPage = c.ReadInt("page", 1, v)
	input.PageSize = c.ReadInt("page_size", 20, v)

	input.Sort = c.ReadString("sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "-id", "-title"}

	if domain.ValidateFilters(v, input.Filters); !v.Valid() {
		return c.FailedValidationResponse(v.Err())
	}

	ctx := c.GetRequest().Context()
	tasks, metadata, err := t.tu.GetAll(ctx, user.ID, input.Title, input.Filters)
	if err != nil {
		return errors.E(op, err)
	}

	err = c.WriteJSON(http.StatusOK, &GetAllTasksResponse{
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
func (t *taskAPI) GetByID(c *reactor.Context) error {
	const op errors.Op = "GetByID"
	user := helpers.ContextGetUser(c.GetRequest())

	id, err := c.ReadIDParam()
	if err != nil {
		return c.NotFoundResponse()
	}

	ctx := c.GetRequest().Context()
	task, err := t.tu.GetByID(ctx, user.ID, id)
	if err != nil {
		switch {
		case errors.KindIs(err, errors.ErrRecordNotFound):
			return c.NotFoundResponse()
		default:
			return errors.E(op, err)
		}
	}

	return c.WriteJSON(http.StatusOK, GetTaskByIDResponse{task})
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
func (t *taskAPI) Insert(c *reactor.Context) error {
	const op errors.Op = "taksAPI.Insert"

	var input CreateTaskRequest
	user := helpers.ContextGetUser(c.GetRequest())
	_ = user

	err := c.ReadJSON(&input)
	if err != nil {
		return c.BadRequestResponse(err)
	}

	task := &domain.Task{
		Title:   input.Title,
		Content: input.Content,
		Done:    input.Done,
	}

	// validate the request
	v := validator.New()

	if domain.ValidateTask(v, task); !v.Valid() {
		return c.FailedValidationResponse(v.Err())
	}

	ctx := c.GetRequest().Context()
	err = t.tu.Insert(ctx, user.ID, task)
	if err != nil {
		return errors.E(op, err)
	}

	c.GetResponseWriter().Header().Set("Location", fmt.Sprintf("/v1/tasks/%d", task.ID))

	// Write a JSON response with a 201 Created status code, the task data in the
	// response body, and the Location header.
	err = c.WriteJSON(http.StatusCreated, &CreateTaskResponse{task})
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
func (t *taskAPI) Update(c *reactor.Context) error {
	const op errors.Op = "taskAPI.Update"

	user := helpers.ContextGetUser(c.GetRequest())

	taskID, err := c.ReadIDParam()
	if err != nil {
		return c.NotFoundResponse()
	}

	ctx := c.GetRequest().Context()
	task, err := t.tu.GetByID(ctx, user.ID, taskID)
	if err != nil {
		return errors.E(op, err)
	}

	var input struct {
		Title   *string `json:"title"`   // task title
		Content *string `json:"content"` // task content
		Done    *bool   `json:"done"`    // true if task is done
	}

	err = c.ReadJSON(&input)
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
		return c.FailedValidationResponse(v.Err())
	}

	ctx = c.GetRequest().Context()
	err = t.tu.Update(ctx, task)
	if err != nil {
		return errors.E(op, err)
	}

	return c.WriteJSON(http.StatusOK, &UpdateTaskByIDResponse{task})
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
func (t *taskAPI) Delete(c *reactor.Context) error {
	const op errors.Op = "taskAPI.Delete"

	user := helpers.ContextGetUser(c.GetRequest())

	// Extract the task ID from the URL.
	taskID, err := c.ReadIDParam()
	if err != nil {
		return c.NotFoundResponse()
	}

	ctx := c.GetRequest().Context()
	err = t.tu.Delete(ctx, user.ID, taskID)
	if err != nil {
		switch {
		case errors.KindIs(err, errors.ErrRecordNotFound):
			return c.NotFoundResponse()
		default:
			return errors.E(op, err)
		}
	}

	return c.WriteJSON(http.StatusOK, &DeleteTaskByIDResponse{"task successfully deleted"})
}
