package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes setup routing of our router.
func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/tasks", app.requirePermission("movies:read", app.listTasksHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tasks", app.requirePermission("movies:write", app.createTaskHandler))
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", app.requirePermission("movies:read", app.showTaskHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/tasks/:id", app.requirePermission("movies:write", app.updateTaskHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/tasks/:id", app.requirePermission("movies:write", app.deleteTaskHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
