package errors

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unknowntpo/todos/internal/logger/zerolog"

	"github.com/stretchr/testify/assert"
)

func errJSONOutput(t *testing.T, msg string) string {
	t.Helper()
	s := `{
	"error": {
		"Message": %q
	}
}
`
	return fmt.Sprintf(s, msg)
}

func TestSendErrorResponse(t *testing.T) {
	t.Run("Internal Server Error Response", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		logBuf := new(bytes.Buffer)
		logger := zerolog.New(logBuf)

		err = E(ErrInternal, New("deliberated internal server error"))

		SendErrorResponse(rr, r, logger, err)
		// Because we don't care about timestamp, so we just check if the output of logger contains error message we want.
		assert.Contains(t, logBuf.String(), `internal server error: deliberated internal server error`, "logger should contain this message")

		wantRespBody := errJSONOutput(t, "the server encountered a problem and could not process your request")
		assert.Equal(t, wantRespBody, rr.Body.String())
	})

	t.Run("Not found Response", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		logBuf := new(bytes.Buffer)
		logger := zerolog.New(logBuf)
		err = E(ErrRecordNotFound, sql.ErrNoRows)
		SendErrorResponse(rr, r, logger, err)

		assert.Equal(t, "", logBuf.String(), "logger output should be empty string")

		wantRespBody := errJSONOutput(t, "the requested resource could not be found")
		assert.Equal(t, wantRespBody, rr.Body.String())
	})

	t.Run("Method Not allowed Response", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodPut, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		logBuf := new(bytes.Buffer)
		logger := zerolog.New(logBuf)
		err = E(ErrMethodNotAllowed)

		SendErrorResponse(rr, r, logger, err)

		assert.Equal(t, "", logBuf.String(), "logger output should be empty string")

		wantRespBody := errJSONOutput(t, fmt.Sprintf("the %s method is not supported for this resource", r.Method))
		assert.Equal(t, wantRespBody, rr.Body.String())
	})

	t.Run("Bad Request Response", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		logBuf := new(bytes.Buffer)
		logger := zerolog.New(logBuf)
		err = E(ErrBadRequest, Op("ReadJSON"))

		SendErrorResponse(rr, r, logger, err)

		assert.Equal(t, "", logBuf.String(), "logger output should be empty string")

		wantRespBody := errJSONOutput(t, err.Error())
		assert.Equal(t, wantRespBody, rr.Body.String())
	})
}
