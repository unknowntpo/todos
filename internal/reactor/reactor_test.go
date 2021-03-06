package reactor

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unknowntpo/todos/internal/domain/errors"
	"github.com/unknowntpo/todos/internal/logger/zerolog"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestHandlerWrapper(t *testing.T) {
	// Set up logger.
	//logBuf := new(bytes.Buffer)
	logBuf := bytes.NewBufferString("")

	logger := zerolog.New(logBuf)

	rc := NewReactor(logger)

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("failed to create new request: %v", err)
	}

	router := httprouter.New()

	h := func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("something goes wrong")

		rc.ServerErrorResponse(w, r, err)
	}

	// attach our handler with rc.HandlerWrapper
	router.HandlerFunc(http.MethodGet, "/", h)

	router.ServeHTTP(rr, r)

	wantMsg := "something goes wrong"

	wantBodyMsg := "the server encountered a problem and could not process your request"
	assert.Contains(t, rr.Body.String(), wantBodyMsg, "response should contain these message")
	assert.Contains(t, logBuf.String(), wantMsg, "output of logger should contain these message")
}
