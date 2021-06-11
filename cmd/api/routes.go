package main

import (
	"net/http"

	_healthcheckAPI "github.com/unknowntpo/todos/internal/healthcheck/delivery/api"

	"github.com/julienschmidt/httprouter"
)

func (app *application) newRoutes() http.Handler {
	// Set up /{api_version}/healthcheck
	router := httprouter.New()
	_healthcheckAPI.NewHealthcheckAPI(router)

	// TODO: Add more api endpoints

	// TODO: Add middleware
	return router
}
