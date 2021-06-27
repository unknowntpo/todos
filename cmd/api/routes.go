package main

import (
	"net/http"
	"time"

	"github.com/unknowntpo/todos/internal/domain/mock"
	_healthcheckAPI "github.com/unknowntpo/todos/internal/healthcheck/delivery/api"
	_taskAPI "github.com/unknowntpo/todos/internal/task/delivery/api"
	_taskUsecase "github.com/unknowntpo/todos/internal/task/usecase"

	"github.com/julienschmidt/httprouter"
)

func (app *application) newRoutes() http.Handler {
	router := httprouter.New()
	_healthcheckAPI.NewHealthcheckAPI(router, version, app.config.env)

	//tRepo := _taskRepo.NewTaskRepository()
	//tu := _taskUsecase.NewTaskUsecase()

	// Use mockUsecase for testing.
	//tu := mock.NewTaskUsecase()
	tr := mock.NewTaskRepo()
	tu := _taskUsecase.NewTaskUsecase(tr, 3*time.Second)
	_taskAPI.NewTaskAPI(router, tu)

	// TODO: Add more api endpoints

	// TODO: Add middleware
	return router
}
