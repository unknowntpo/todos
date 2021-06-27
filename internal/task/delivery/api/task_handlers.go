package api

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/helpers/validator"

	"github.com/julienschmidt/httprouter"
)

type TaskAPI struct {
	TU domain.TaskUsecase
}

func NewTaskAPI(router *httprouter.Router, tu domain.TaskUsecase) {
	handler := &TaskAPI{TU: tu}
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", handler.GetByID)
	router.HandlerFunc(http.MethodPatch, "/v1/tasks/:id", handler.Update)
}

func (t *TaskAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		helpers.NotFoundResponse(w, r)
		return
	}

	ctx := r.Context()
	task, err := t.TU.GetByID(ctx, id)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"task": task}, nil)
}

func (t *TaskAPI) Update(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		helpers.NotFoundResponse(w, r)
		return
	}

	ctx := r.Context()
	task, err := t.TU.GetByID(ctx, id)
	if err != nil {
		// TODO: errors.Is
		helpers.ServerErrorResponse(w, r, err)
	}

	var input struct {
		Title   *string `json:"title"`   // task title
		Content *string `json:"content"` // task content
		Done    *bool   `json:"done"`    // true if task is done
	}
	// readJSON
	err = helpers.ReadJSON(w, r, &input)

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
