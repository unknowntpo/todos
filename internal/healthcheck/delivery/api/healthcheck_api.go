package api

import (
	"net/http"

	swagger "github.com/unknowntpo/todos/docs/go"
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

func (h *healthcheckAPI) Healthcheck(w http.ResponseWriter, r *http.Request) {
	panic("delibrated panic at healthcheck api")

	err := helpers.WriteJSON(w, http.StatusOK, &swagger.HealthcheckResponse{
		Status:      "available",
		Version:     h.version,
		Environment: h.env,
	}, nil)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
	}
}
