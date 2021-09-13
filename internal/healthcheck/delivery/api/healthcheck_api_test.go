package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheckHandler(t *testing.T) {
	// build route
	router := httprouter.New()
	NewHealthcheckAPI(router, "v3.14", "development")

	// build request: GET /v1/healthcheck
	r, err := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	if err != nil {
		t.Fatalf("failed to make new request: %v", err)
	}

	// make responserecorder
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	wantBody :=
		`{
	"status": "available",
	"environment": "development",
	"version": "v3.14"
}
`
	assert.Equal(t, wantBody, w.Body.String(), "response body must be equal")
}
