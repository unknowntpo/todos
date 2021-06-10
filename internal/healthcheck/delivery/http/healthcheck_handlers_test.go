package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

// TODO: Write test.
func TestHealcheckHandler(t *testing.T) {
	// Because healthcheck delivery doesn't depend on anything, so
	// we just use HealcheckDelivery{}.
	router := httprouter.New()

	handler := &HealthcheckDelivery{}
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handler.HealthcheckHandler)

	// Set request to /v1/healthcheck
	r, _ := http.NewRequest("GET", "/v1/healthcheck", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, r)
	// assert status code
	assertStatusCode(t, rr.Code, http.StatusOK)

	// assert response
	assertBody(t, rr.Body.String(), "OK")
}

// TODO: Move it to helper pkg
func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("Wrong status code, got %d, want %d", got, want)
	}
}

// TODO: Move it to helper pkg
func assertBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Wrong response body, got %s, want %s", got, want)
	}
}
