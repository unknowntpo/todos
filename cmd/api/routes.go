package main

import (
	"net/http"
	"time"

	"github.com/unknowntpo/todos/internal/domain/mock"
	_healthcheckAPI "github.com/unknowntpo/todos/internal/healthcheck/delivery/api"
	_taskAPI "github.com/unknowntpo/todos/internal/task/delivery/api"
	_taskRepoPostgres "github.com/unknowntpo/todos/internal/task/repository/postgres"
	_taskUsecase "github.com/unknowntpo/todos/internal/task/usecase"
	_userAPI "github.com/unknowntpo/todos/internal/user/delivery/api"

	"github.com/julienschmidt/httprouter"
)

func (app *application) newRoutes() http.Handler {
	router := httprouter.New()
	_healthcheckAPI.NewHealthcheckAPI(router, version, app.config.env)

	//tRepo := _taskRepo.NewTaskRepository()
	//tu := _taskUsecase.NewTaskUsecase()

	// Use mockUsecase for testing.
	//tu := mock.NewTaskUsecase()
	tr := _taskRepoPostgres.NewTaskRepo(app.database)
	//tr := mock.NewTaskRepo()
	tu := _taskUsecase.NewTaskUsecase(tr, 3*time.Second)
	_taskAPI.NewTaskAPI(router, tu)

	// Use mockUseCase for testing.
	uu := mock.NewUserUsecase()
	_userAPI.NewUserAPI(router, uu)

	// TODO: Add more api endpoints

	// TODO: Add middleware
	return router
}
