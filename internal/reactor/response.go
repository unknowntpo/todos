package reactor

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	errMsg interface{} `json: "error"`
}

func (c *Context) errorResponse(status int, message interface{}) error {
	const op errors.Op = "errResponse"
	err := c.WriteJSON(c.w, status, &ErrorResponse{
		errMsg: message,
	})
	if err != nil {
		return errors.E(op, err)
	}
}

func (c *Context) ServerErrorResponse(err error) {
	const op errors.Op = "ServerErrorResponse"
	c.logger.PrintError(err, nil)
	message := "the server encountered a problem and could not process your request"
	if e := errorResponse(http.StatusInternalServerError, message); e != nil {
		c.logger.PrintError(err, nil)
		c.w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *Context) NotFoundResponse() error {
	message := "the requested resource could not be found"
	return errorResponse(http.StatusNotFound, message)
}

func (c *Context) MethodNotAllowedResponse() error {
	message := fmt.Sprintf("the %s method is not supported for this resource", c.r.Method)
	return errorResponse(http.StatusMethodNotAllowed, message)
}

func (c *Context) BadRequestResponse(err error) error {
	return errorResponse(http.StatusBadRequest, err.Error())
}

func (c *Context) FailedValidationResponse(err validator.ValidationErrors) error {
	return errorResponse(http.StatusUnprocessableEntity, err)
}

func (c *Context) EditConflictResponse() error {
	message := "unable to update the record due to an edit conflict, please try again"
	return errorResponse(http.StatusConflict, message)
}

func (c *Context) InvalidCredentialsResponse() error {
	message := "invalid authentication credentials"
	return errorResponse(http.StatusUnauthorized, message)
}

func (c *Context) InvalidAuthenticationTokenResponse() error {
	c.w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	return errorResponse(http.StatusUnauthorized, message)
}

func (c *Context) AuthenticationRequiredResponse() error {
	message := "you must be authenticated to access this resource"
	return errorResponse(http.StatusUnauthorized, message)
}

func (c *Context) InactiveAccountResponse() error {
	message := "your user account must be activated to access this resource"
	return errorResponse(http.StatusForbidden, message)
}

func (c *Context) NotPermittedResponse() error {
	message := "your user account doesn't have the necessary permissions to access this resource"
	return errorResponse(http.StatusForbidden, message)
}

func (c *Context) RateLimitExceededResponse() error {
	message := "rate limit exceeded"
	errorResponse(http.StatusTooManyRequests, message)
}
