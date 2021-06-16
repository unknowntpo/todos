package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func newMockHealthcheckAPI() *HealthcheckAPI {
	return &HealthcheckAPI{version: "1.0.0", env: "development"}
}

func TestHealcheckHandler(t *testing.T) {
	router := httprouter.New()

	handler := newMockHealthcheckAPI()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handler.HealthcheckHandler)

	// Set request to /v1/healthcheck
	r, _ := http.NewRequest("GET", "/v1/healthcheck", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusOK, rr.Code, "Wrong status code")

	wantBody := `{
	"status": "available",
	"system_info": {
		"environment": "development",
		"version": "1.0.0"
	}
}
`
	assert.Equal(t, wantBody, rr.Body.String(), "Wrong response body")
}
