package errors

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/logger"
)

type ErrorResponseWrapper struct {
	Resp *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Message string
	// TODO: What field should we include ?
}

func SendErrorResponse(w http.ResponseWriter, logger logger.Logger, err error) {
	var msg string
	var status int

	if e, ok := err.(*Error); ok {
		switch e.Kind {
		case ErrInternal:
			// log the error
			logger.PrintError(e, nil)
			status = http.StatusInternalServerError
			msg = "the server encountered a problem and could not process your request"
		}

		err := helpers.WriteJSON(w, status, &ErrorResponseWrapper{
			Resp: &ErrorResponse{
				Message: msg,
			},
		}, nil)
		if err != nil {
			// TODO: We need to wrap our error with message just like fmt.Errorf()
			logger.PrintError(err, nil)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		// FIXME: Not error struct we defined, what should we do ?
	}
}
