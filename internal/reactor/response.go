package reactor

import (
	"fmt"
	"net/http"

	"github.com/unknowntpo/todos/internal/domain/errors"
)

type ErrorResponse struct {
	ErrMsg interface{} `json:"error"`
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) error {
	const op errors.Op = "errResponse"
	err := WriteJSON(w, status, &ErrorResponse{
		ErrMsg: message,
	})
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "ServerErrorResponse"
	message := "the server encountered a problem and could not process your request"
	if e := errorResponse(w, r, http.StatusInternalServerError, message); e != nil {
		return errors.E(op, e)
	}

	return nil
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) error {
	message := "the requested resource could not be found"
	return errorResponse(w, r, http.StatusNotFound, message)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) error {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	return errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) error {
	return errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, err error) error {
	return errorResponse(w, r, http.StatusUnprocessableEntity, err)
}

func EditConflictResponse(w http.ResponseWriter, r *http.Request) error {
	message := "unable to update the record due to an edit conflict, please try again"
	return errorResponse(w, r, http.StatusConflict, message)
}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) error {
	message := "invalid authentication credentials"
	return errorResponse(w, r, http.StatusUnauthorized, message)
}

func InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	return errorResponse(w, r, http.StatusUnauthorized, message)
}

func AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) error {
	message := "you must be authenticated to access this resource"
	return errorResponse(w, r, http.StatusUnauthorized, message)
}

func InactiveAccountResponse(w http.ResponseWriter, r *http.Request) error {
	message := "your user account must be activated to access this resource"
	return errorResponse(w, r, http.StatusForbidden, message)
}

func NotPermittedResponse(w http.ResponseWriter, r *http.Request) error {
	message := "your user account doesn't have the necessary permissions to access this resource"
	return errorResponse(w, r, http.StatusForbidden, message)
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) error {
	message := "rate limit exceeded"
	return errorResponse(w, r, http.StatusTooManyRequests, message)
}
