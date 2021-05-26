package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes setup routing of our router.
func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/tasks", app.listTasksHandler)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tasks", app.createTaskHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", app.showTaskHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/tasks/:id", app.updateTaskHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/tasks/:id", app.deleteTaskHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
