package main

import (
	"net/http"

	_healthcheckHttpDelivery "github.com/unknowntpo/todos/internal/healthcheck/delivery/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) newRoutes() http.Handler {
	// Set up /{api_version}/healthcheck
	router := httprouter.New()
	_healthcheckHttpDelivery.NewHealthcheckHTTPDelivery(router)

	// TODO: Add more api endpoints

	// TODO: Add middleware
	return router
}
