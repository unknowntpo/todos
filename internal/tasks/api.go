package tasks

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/entity"
)

func NewAPI(route http.Handler, TU entity.TaskUseCase) http.Handler {}

// showTaskHandler shows the detail of specific task.
// Request URL: GET /v1/tasks/:id
// TODO: tidy up resource
func (r resource) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// task is *data.Task
	//task, err := app.models.Tasks.Get(id)
	task, err := t.Service
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
