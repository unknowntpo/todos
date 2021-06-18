package api

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type TaskAPI struct {
	TU     domain.TaskUsecase
	logger *log.Logger
}

func NewTaskAPI(router *httprouter.Router, tu domain.TaskUsecase, logger *log.Logger) {
	handler := &TaskAPI{TU: tu, logger: logger}
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", handler.GetByID)
}

func (t *TaskAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	// Use mockTaskUsecase
	mockID := 1
	task, err := t.TU.GetByID(int64(mockID))
	t.logger.Info("Debug GetByID")
	if err != nil {
		t.logger.Error(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"task": task}, nil)
}
