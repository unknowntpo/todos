package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HealthcheckAPI represent delivery layer of healcheck http api endpoint.
type HealthcheckAPI struct{}

// NewHealthcheckAPI registers all handlers in /v1/healcheck to the router.
func NewHealthcheckAPI(router *httprouter.Router) {
	handler := &HealthcheckAPI{}
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handler.HealthcheckHandler)
}

// TODO: Display metrics.
func (h *HealthcheckAPI) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
