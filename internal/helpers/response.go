package helpers

import (
	"fmt"
	"net/http"
)

// Dependency: logger
func logError(r *http.Request, err error) {
	// TODO: Decide which logger to use
	// jsonlog?
	// Where to store jsonlog?

	/*
		app.logger.PrintError(err, map[string]string{
			"request_method": r.Method,
			"request_url":    r.URL.String(),
		})
	*/
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := Envelope{"error": message}

	err := WriteJSON(w, status, env, nil)
	if err != nil {
		// TODO: logger interface
		logError(r, err)
		w.WriteHeader(500)
	}
}

// TODO: Logger interface
func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// Note that the errors parameter here has the type map[string]string, which is exactly
// the same as the errors map contained in our Validator type.
func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	errorResponse(w, r, http.StatusConflict, message)
}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	errorResponse(w, r, http.StatusUnauthorized, message)
}

func InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	errorResponse(w, r, http.StatusUnauthorized, message)
}

func AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	errorResponse(w, r, http.StatusUnauthorized, message)
}

func InactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	errorResponse(w, r, http.StatusForbidden, message)
}

func NotPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	errorResponse(w, r, http.StatusForbidden, message)
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	errorResponse(w, r, http.StatusTooManyRequests, message)
}
