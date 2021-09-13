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

	taskRepo := mock.NewTaskRepo()
	userRepo := _userRepoPostgres.NewUserRepo(app.database)
	tokenRepo := _tokenRepoPostgres.NewTokenRepo(app.database)

	// usecase
	taskUsecase := _taskUsecase.NewTaskUsecase(taskRepo, 3*time.Second)
	userUsecase := _userUsecase.NewUserUsecase(userRepo, 3*time.Second)
	tokenUsecase := _tokenUsecase.NewTokenUsecase(tokenRepo, 3*time.Second)

	// delivery
	_taskAPI.NewTaskAPI(router, taskUsecase, app.logger)
	_userAPI.NewUserAPI(router, userUsecase, tokenUsecase, app.pool, app.mailer, app.logger)
	_tokenAPI.NewTokenAPI(router, tokenUsecase, userUsecase, app.logger)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// TODO: Add middleware
	genMid := _generalMiddleware.New(app.config, userUsecase, app.logger)

	return genMid.Metrics(genMid.RecoverPanic(genMid.EnableCORS(genMid.RateLimit(genMid.Authenticate(router)))))
	//return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))

}
