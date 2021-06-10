package http

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// HealthcheckDelivery represent delivery layer of healcheck api endpoint.
type HealthcheckDelivery struct{}

// NewHealthcheckDelivery registers all handlers in /v1/healcheck to the router.
func NewHealthcheckDelivery(router *httprouter.Router) {
	handler := &HealthcheckDelivery{}
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handler.HealthcheckHandler)
}

// TODO: Display metrics.
func (h *HealthcheckDelivery) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
