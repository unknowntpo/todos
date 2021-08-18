package api

import (
	"errors"
	"net/http"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/helpers/validator"
	"github.com/unknowntpo/todos/internal/logger"

	"github.com/julienschmidt/httprouter"
)

type taskAPI struct {
	TU domain.TaskUsecase
}

func NewTaskAPI(router *httprouter.Router, tu domain.TaskUsecase) {
	api := &taskAPI{TU: tu}
	router.HandlerFunc(http.MethodGet, "/v1/tasks", api.GetAll)
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", api.GetByID)
	router.HandlerFunc(http.MethodPost, "/v1/tasks", api.Insert)
	router.HandlerFunc(http.MethodPatch, "/v1/tasks/:id", api.Update)
	router.HandlerFunc(http.MethodDelete, "/v1/tasks/:id", api.Delete)
}

// GetAll gets all tasks.
// TODO: GetAll should get all tasks with specific user id.
func (t *taskAPI) GetAll(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"debug": "GetAll called"}, nil)

}

// GetByID gets a task by its id.
func (t *taskAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		helpers.NotFoundResponse(w, r)
		return
	}

	ctx := r.Context()
	task, err := t.TU.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			helpers.NotFoundResponse(w, r)
		default:
			logger.Log.PrintError(err, nil)
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
		logger.Log.PrintError(err, nil)
		helpers.NotFoundResponse(w, r)
		return
	}

	ctx := r.Context()
	task, err := t.TU.GetByID(ctx, id)
	if err != nil {
		// TODO: errors.Is
		logger.Log.PrintError(err, nil)
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
	if validateTask(v, task); !v.Valid() {
		helpers.FailedValidationResponse(w, r, v.Errors)
		return
	}

	// TODO: We validate input at delivery layer.

	ctx = r.Context()
	taskUpdated, err := t.TU.Update(ctx, id, task)
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
