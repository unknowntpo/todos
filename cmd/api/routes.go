package main

import (
	"net/http"
	"time"

	"github.com/unknowntpo/todos/internal/domain/mock"
	_healthcheckAPI "github.com/unknowntpo/todos/internal/healthcheck/delivery/api"

	_taskAPI "github.com/unknowntpo/todos/internal/task/delivery/api"
	_taskUsecase "github.com/unknowntpo/todos/internal/task/usecase"
	//_taskRepoPostgres "github.com/unknowntpo/todos/internal/task/repository/postgres"

	//_userRepoPostgres "github.com/unknowntpo/todos/internal/user/repository/postgres"
	_userAPI "github.com/unknowntpo/todos/internal/user/delivery/api"
	_userUsecase "github.com/unknowntpo/todos/internal/user/usecase"

	//_tokenRepoPostgres "github.com/unknowntpo/todos/internal/token/repository/postgres"
	_tokenAPI "github.com/unknowntpo/todos/internal/token/delivery/api"
	_tokenUsecase "github.com/unknowntpo/todos/internal/token/usecase"

	_generalMiddleware "github.com/unknowntpo/todos/internal/middleware"

	"github.com/julienschmidt/httprouter"
)

func (app *application) newRoutes() http.Handler {
	router := httprouter.New()
	_healthcheckAPI.NewHealthcheckAPI(router, version, app.config.env)

	// mock the repo
	// maybe its stub ?
	taskRepo := mock.NewTaskRepo()
	taskUsecase := _taskUsecase.NewTaskUsecase(taskRepo, 3*time.Second)
	_taskAPI.NewTaskAPI(router, taskUsecase)

	userRepo := mock.NewUserRepo()
	userUsecase := _userUsecase.NewUserUsecase(userRepo, 3*time.Second)
	_userAPI.NewUserAPI(router, userUsecase, &app.bg)

	// TODO: Add more api endpoints
	tokenRepo := mock.NewTokenRepo()
	tokenUsecase := _tokenUsecase.NewTokenUsecase(tokenRepo, 3*time.Second)
	_tokenAPI.NewTokenAPI(router, tokenUsecase, userUsecase)

	// TODO: Add middleware
	genMid := _generalMiddleware.New()

	return genMid.RecoverPanic(router)
}
