package api

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"

	"github.com/julienschmidt/httprouter"
)

type TaskAPI struct {
	TU domain.TaskUsecase
}

func NewTaskAPI(router *httprouter.Router, tu domain.TaskUsecase) {
	handler := &TaskAPI{TU: tu}
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", handler.GetByID)
}

func (t *TaskAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		helpers.NotFoundResponse(w, r)
		return
	}

	task, err := t.TU.GetByID(id)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"task": task}, nil)
}
