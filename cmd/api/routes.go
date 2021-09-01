package main

import (
	"expvar"
	"net/http"
	"time"

	"github.com/unknowntpo/todos/internal/domain/mock"
	_healthcheckAPI "github.com/unknowntpo/todos/internal/healthcheck/delivery/api"

	_taskAPI "github.com/unknowntpo/todos/internal/task/delivery/api"
	_taskUsecase "github.com/unknowntpo/todos/internal/task/usecase"
	//_taskRepoPostgres "github.com/unknowntpo/todos/internal/task/repository/postgres"

	_userAPI "github.com/unknowntpo/todos/internal/user/delivery/api"
	_userRepoPostgres "github.com/unknowntpo/todos/internal/user/repository/postgres"
	_userUsecase "github.com/unknowntpo/todos/internal/user/usecase"

	_tokenAPI "github.com/unknowntpo/todos/internal/token/delivery/api"
	_tokenRepoPostgres "github.com/unknowntpo/todos/internal/token/repository/postgres"
	_tokenUsecase "github.com/unknowntpo/todos/internal/token/usecase"

	_generalMiddleware "github.com/unknowntpo/todos/internal/middleware"

	"github.com/julienschmidt/httprouter"
)

func (app *application) newRoutes() http.Handler {
	router := httprouter.New()
	_healthcheckAPI.NewHealthcheckAPI(router, version, app.config.Env)

	// mock the repo
	// maybe its stub ?
	taskRepo := mock.NewTaskRepo()
	taskUsecase := _taskUsecase.NewTaskUsecase(taskRepo, 3*time.Second)
	_taskAPI.NewTaskAPI(router, taskUsecase)

	userRepo := _userRepoPostgres.NewUserRepo(app.database)
	userUsecase := _userUsecase.NewUserUsecase(userRepo, 3*time.Second)
	_userAPI.NewUserAPI(router, userUsecase, &app.bg)

	tokenRepo := _tokenRepoPostgres.NewTokenRepo(app.database)
	tokenUsecase := _tokenUsecase.NewTokenUsecase(tokenRepo, 3*time.Second)
	_tokenAPI.NewTokenAPI(router, tokenUsecase, userUsecase)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// TODO: Add middleware
	genMid := _generalMiddleware.New(app.config, userUsecase)

	return genMid.Metrics(genMid.RecoverPanic(genMid.EnableCORS(genMid.RateLimit(genMid.Authenticate(router)))))
	//return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))

}
