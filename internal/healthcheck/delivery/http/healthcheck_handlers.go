package http

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type HealthcheckDelivery struct{}

func NewHealthcheckDelivery(router *httprouter.Router) {
	handler := &HealthcheckDelivery{}
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handler.HealthcheckHandler)
}

func (h *HealthcheckDelivery) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
