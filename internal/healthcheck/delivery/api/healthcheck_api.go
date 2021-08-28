package api

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/helpers"

	"github.com/julienschmidt/httprouter"
)

// healthcheckAPI represent delivery layer of healcheck http api endpoint.
type healthcheckAPI struct {
	version string
	env     string
}

// NewHealthcheckAPI registers all handlers in /v1/healcheck to the router.
func NewHealthcheckAPI(router *httprouter.Router, version, env string) {
	api := &healthcheckAPI{version: version, env: env}
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", api.Healthcheck)
}

// TODO: Display metrics.
func (h *healthcheckAPI) Healthcheck(w http.ResponseWriter, r *http.Request) {
	env := helpers.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.env,
			"version":     h.version,
		},
	}

	err := helpers.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
	}
}