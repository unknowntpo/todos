package http

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// HealthcheckHTTPDelivery represent delivery layer of healcheck http api endpoint.
type HealthcheckHTTPDelivery struct{}

// NewHealthcheckHTTPDelivery registers all handlers in /v1/healcheck to the router.
func NewHealthcheckHTTPDelivery(router *httprouter.Router) {
	handler := &HealthcheckHTTPDelivery{}
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handler.HealthcheckHandler)
}

// TODO: Display metrics.
func (h *HealthcheckHTTPDelivery) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
