package errors

import (
	"fmt"
	"net/http"
)

type ErrorResponseWrapper struct {
	Resp *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Status  int
	Message interface{}
	// TODO: What field should we include ?
}

type ErrResponseKind uint8

// Kinds of reponse errors.
const (
	OtherErrorResponse ErrResponseKind = iota // Unclassified error. This value is not printed in the error message.
	ServerErrorResponse
	NotFoundResponse
	MethodNotAllowedResponse
	BadRequestResponse
	FailedValidationResponse
	EditConflictResponse
	InvalidCredentialsResponse
	InvalidAuthenticationTokenResponse
	AuthenticationRequiredResponse
	InactiveAccountResponse
)

func (k ErrResponseKind) String() string {
	switch k {
	case OtherErrorResponse:
		return "other error response"
	case ServerErrorResponse:
		return "server error response"
	case NotFoundResponse:
		return "not found error response"
	case MethodNotAllowedResponse:
		return "method not allowed error response"
	case BadRequestResponse:
		return "bad request other error response"
	case FailedValidationResponse:
		return "failed validation error response"
	case EditConflictResponse:
		return "edit conflict other error response"
	case InvalidCredentialsResponse:
		return "invalid credentials error response"
	case InvalidAuthenticationTokenResponse:
		return "invalid authentication token error response"
	case AuthenticationRequiredResponse:
		return "authentication required error response"
	case InactiveAccountResponse:
		return "inactive account error response"
	}
	return "unknown error response kind"
}

func NewErrorResponse(w http.ResponseWriter, r *http.Request, kind ErrResponseKind, err error) *ErrorResponseWrapper {
	var msg interface{}
	var status int

	switch kind {
	case ServerErrorResponse:
		status = http.StatusInternalServerError
		msg = "the server encountered a problem and could not process your request"
	case NotFoundResponse:
		status = http.StatusNotFound
		msg = "the requested resource could not be found"
	case MethodNotAllowedResponse:
		status = http.StatusMethodNotAllowed
		msg = fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	case BadRequestResponse:
		status = http.StatusBadRequest
		msg = err
	case FailedValidationResponse:
		status = http.StatusUnprocessableEntity
		msg = err
	case EditConflictResponse:
		msg = "unable to update the record due to an edit conflict, please try again"
		status = http.StatusConflict
	case InvalidCredentialsResponse:
		msg = "invalid authentication credentials"
		status = http.StatusUnauthorized
	case InvalidAuthenticationTokenResponse:
		w.Header().Set("WWW-Authenticate", "Bearer")
		msg = "invalid or missing authentication token"
		status = http.StatusUnauthorized
	case AuthenticationRequiredResponse:
		msg = "you must be authenticated to access this resource"
		status = http.StatusUnauthorized
	case InactiveAccountResponse:
		msg = "your user account must be activated to access this resource"
		status = http.StatusForbidden
	default:
		panic("unknown error response kind")
	}

	return &ErrorResponseWrapper{
		&ErrorResponse{
			Status:  status,
			Message: msg,
		}}
}

/*
func SendErrorResponse(w http.ResponseWriter, r *http.Request, logger logger.Logger, err error) {
	var msg string
	var status int

	if e, ok := err.(*Error); ok {
		switch e.Kind {
		case ErrInternal:
			// log the error
			logger.PrintError(e, nil)
			status = http.StatusInternalServerError
			msg = "the server encountered a problem and could not process your request"
		case ErrRecordNotFound:
			status = http.StatusNotFound
			msg = "the requested resource could not be found"
		case ErrMethodNotAllowed:
			status = http.StatusMethodNotAllowed
			msg = fmt.Sprintf("the %s method is not supported for this resource", r.Method)
		case ErrBadRequest:
			status = http.StatusBadRequest
			msg = err.Error()
		default:
			panic("ServerErrorResponse: unknown type of error")
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
*/
