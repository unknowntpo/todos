package errors

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/unknowntpo/todos/internal/logger/zerolog"

	"github.com/stretchr/testify/assert"
)

func TestSendErrorResponse(t *testing.T) {
	t.Run("ErrInternal", func(t *testing.T) {
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

		wantRespBody :=
			`{
	"error": {
		"Message": "the server encountered a problem and could not process your request"
	}
}
`
		assert.Equal(t, wantRespBody, rr.Body.String())
	})
}
