package reactor

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	ErrMsg interface{} `json:"error"`
}

func (rc *Reactor) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	err := rc.WriteJSON(w, status, &ErrorResponse{
		ErrMsg: message,
	})
	if err != nil {
		// Something goes wrong during sending server error response.
		// So we write the message directly.
		rc.Logger.PrintError(err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		msg := `{"error":"the server encountered a problem and could not process your request"}`
		w.Write([]byte(msg))
		return
	}
}

func (rc *Reactor) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	rc.Logger.PrintError(err, nil)

	rc.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (rc *Reactor) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	rc.errorResponse(w, r, http.StatusNotFound, message)
}

func (rc *Reactor) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	rc.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (rc *Reactor) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	rc.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (rc *Reactor) FailedValidationResponse(w http.ResponseWriter, r *http.Request, err error) {
	rc.errorResponse(w, r, http.StatusUnprocessableEntity, err)
}

func (rc *Reactor) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	rc.errorResponse(w, r, http.StatusConflict, message)
}

func (rc *Reactor) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	rc.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (rc *Reactor) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	rc.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (rc *Reactor) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	rc.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (rc *Reactor) InactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	rc.errorResponse(w, r, http.StatusForbidden, message)
}

func (rc *Reactor) NotPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	rc.errorResponse(w, r, http.StatusForbidden, message)
}

func (rc *Reactor) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	rc.errorResponse(w, r, http.StatusTooManyRequests, message)
}
