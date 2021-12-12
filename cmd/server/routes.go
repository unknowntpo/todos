package main

import (
	"expvar"
	"net/http"
	"time"

	_healthcheckAPI "github.com/unknowntpo/todos/internal/healthcheck/delivery/api"

	_taskAPI "github.com/unknowntpo/todos/internal/task/delivery/api"
	_taskRepoPostgres "github.com/unknowntpo/todos/internal/task/repository/postgres"
	_taskUsecase "github.com/unknowntpo/todos/internal/task/usecase"

	_userAPI "github.com/unknowntpo/todos/internal/user/delivery/api"
	_userRepoPostgres "github.com/unknowntpo/todos/internal/user/repository/postgres"
	_userUsecase "github.com/unknowntpo/todos/internal/user/usecase"

	_tokenAPI "github.com/unknowntpo/todos/internal/token/delivery/api"
	_tokenRepoPostgres "github.com/unknowntpo/todos/internal/token/repository/postgres"
	_tokenUsecase "github.com/unknowntpo/todos/internal/token/usecase"

	_generalMiddleware "github.com/unknowntpo/todos/internal/middleware"

	"github.com/unknowntpo/todos/internal/reactor"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/unknowntpo/todos/docs"
)

func (app *application) newRoutes() http.Handler {
	router := httprouter.New()

	taskRepo := _taskRepoPostgres.NewTaskRepo(app.database)
	userRepo := _userRepoPostgres.NewUserRepo(app.database)
	tokenRepo := _tokenRepoPostgres.NewTokenRepo(app.database)

	// usecase
	taskUsecase := _taskUsecase.NewTaskUsecase(taskRepo, 3*time.Second)
	userUsecase := _userUsecase.NewUserUsecase(userRepo, tokenRepo, app.pool, app.mailer, app.logger, 3*time.Second)
	tokenUsecase := _tokenUsecase.NewTokenUsecase(tokenRepo, 3*time.Second)

	// reactor
	rc := reactor.NewReactor(app.logger)

	// middleware

	genMid := _generalMiddleware.New(app.config, userUsecase, rc)

	// delivery

	_healthcheckAPI.NewHealthcheckAPI(router, version, app.config.Env, rc)

	_taskAPI.NewTaskAPI(router, taskUsecase, genMid, rc)
	_userAPI.NewUserAPI(router, userUsecase, tokenUsecase, rc)
	_tokenAPI.NewTokenAPI(router, userUsecase, rc)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// Set up http-swagger
	urlFunc := httpSwagger.URL("https://todos.unknowntpo.net/swagger/doc.json") //The url pointing to API definition
	if app.config.Env == "production" {
		urlFunc = httpSwagger.URL("https://todos.unknowntpo.net/swagger/doc.json") //The url pointing to API definition
		// if in production env, change the host.
		docs.SwaggerInfo.Host = "todos.unknowntpo.net"
	}
	router.Handler(http.MethodGet, "/swagger/:any", httpSwagger.Handler(
		urlFunc,
	))

	//return genMid.Metrics(genMid.RecoverPanic(genMid.EnableCORS(genMid.RateLimit(genMid.Authenticate(router)))))
	//return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
	return chain(router, genMid.Metrics,
		genMid.RecoverPanic,
		genMid.EnableCORS,
		genMid.RateLimit,
		genMid.Authenticate)
}

func chain(route http.Handler, handlers ...func(http.Handler) http.Handler) http.Handler {
	c := route
	for i := len(handlers) - 1; i > 0; i-- {
		c = handlers[i](c)
	}
	return c
}
