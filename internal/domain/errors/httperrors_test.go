package errors

import (
	"bytes"
	"database/sql"
	"fmt"
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
		// Set up bytes.Buffer to write response
		rr := httptest.NewRecorder()
		// Set up bytes.Buffer for logger to write log message to.
		logBuf := new(bytes.Buffer)
		// init logger
		logger := zerolog.New(logBuf)
		// make a error with kind = ErrInternal
		err := E(ErrInternal, "deliberated internal server error")
		// Call SendErrorResponse(httpbuf, logger, err)
		SendErrorResponse(rr, logger, err)
		// Because we don't care about timestamp, so we just check if the output of logger contains error message we want.
		assert.Contains(t, logBuf.String(), `internal server error: deliberated internal server error`, "logger should contain this message")

		wantRespBody := errJSONOutput(t, "the server encountered a problem and could not process your request")
		assert.Equal(t, wantRespBody, rr.Body.String())
	})

	t.Run("Not found Response", func(t *testing.T) {
		// Set up bytes.Buffer to write response
		rr := httptest.NewRecorder()
		// Set up bytes.Buffer for logger to write log message to.
		logBuf := new(bytes.Buffer)
		// init logger
		logger := zerolog.New(logBuf)
		// make a error with kind = ErrInternal
		err := E(ErrRecordNotFound, sql.ErrNoRows)
		// Call SendErrorResponse(httpbuf, logger, err)
		SendErrorResponse(rr, logger, err)

		assert.Equal(t, "", logBuf.String(), "logger output should be empty string")

		wantRespBody := errJSONOutput(t, "the requested resource could not be found")
		assert.Equal(t, wantRespBody, rr.Body.String())
	})

}
