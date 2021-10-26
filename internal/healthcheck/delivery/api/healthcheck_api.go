package api

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/reactor"

	"github.com/julienschmidt/httprouter"
)

// healthcheckAPI represent delivery layer of healcheck http api endpoint.
type healthcheckAPI struct {
	version string
	env     string
	rc      *reactor.Reactor
}

type HealthcheckResponse struct {
	Status      string `json: "status"`
	Environment string `json: "environment"`
	Version     string `json: "version"`
}

// NewHealthcheckAPI registers all handlers in /v1/healcheck to the router.
func NewHealthcheckAPI(router *httprouter.Router, version, env string, rc *reactor.Reactor) {
	api := &healthcheckAPI{version: version, env: env, rc: rc}
	router.Handler(http.MethodGet, "/v1/healthcheck", rc.HandlerWrapper(api.Healthcheck))
}

// Healthcheck shows status of service.
// @Summary Show status of service.
// @Description None.
// @Produce json
// @Success 200 {object} HealthcheckResponse
// @Router /v1/healthcheck [get]
func (h *healthcheckAPI) Healthcheck(w http.ResponseWriter, r *http.Request) error {
	return reactor.WriteJSON(w, http.StatusOK, &HealthcheckResponse{
		Status:      "available",
		Version:     h.version,
		Environment: h.env,
	})
}
