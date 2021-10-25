package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unknowntpo/todos/internal/logger/zerolog"
	"github.com/unknowntpo/todos/internal/reactor"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestHealthcheckHandler(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := zerolog.New(buf)
	router := httprouter.New()

	rc := reactor.NewReactor(logger)

	NewHealthcheckAPI(router, "v3.14", "development", rc)

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
	"Status": "available",
	"Environment": "development",
	"Version": "v3.14"
}
`
	assert.Equal(t, wantBody, w.Body.String(), "response body must be equal")
	assert.Equal(t, "", buf.String(), "buf should contain nothing")
}
