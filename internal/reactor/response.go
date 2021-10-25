package reactor

import (
	"fmt"
	"net/http"

	"github.com/unknowntpo/todos/internal/domain/errors"
	"github.com/unknowntpo/todos/pkg/validator"
)

type ErrorResponse struct {
	ErrMsg interface{} `json:"error"`
}

func (c *Context) errorResponse(status int, message interface{}) error {
	const op errors.Op = "errResponse"
	err := c.WriteJSON(status, &ErrorResponse{
		ErrMsg: message,
	})
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

func (c *Context) ServerErrorResponse() error {
	const op errors.Op = "ServerErrorResponse"
	message := "the server encountered a problem and could not process your request"
	if e := c.errorResponse(http.StatusInternalServerError, message); e != nil {
		return errors.E(op, e)
	}

	return nil
}

func (c *Context) NotFoundResponse() error {
	message := "the requested resource could not be found"
	return c.errorResponse(http.StatusNotFound, message)
}

func (c *Context) MethodNotAllowedResponse() error {
	message := fmt.Sprintf("the %s method is not supported for this resource", c.r.Method)
	return c.errorResponse(http.StatusMethodNotAllowed, message)
}

func (c *Context) BadRequestResponse(err error) error {
	return c.errorResponse(http.StatusBadRequest, err.Error())
}

func (c *Context) FailedValidationResponse(err validator.ValidationErrors) error {
	return c.errorResponse(http.StatusUnprocessableEntity, err)
}

func (c *Context) EditConflictResponse() error {
	message := "unable to update the record due to an edit conflict, please try again"
	return c.errorResponse(http.StatusConflict, message)
}

func (c *Context) InvalidCredentialsResponse() error {
	message := "invalid authentication credentials"
	return c.errorResponse(http.StatusUnauthorized, message)
}

func (c *Context) InvalidAuthenticationTokenResponse() error {
	c.w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	return c.errorResponse(http.StatusUnauthorized, message)
}

func (c *Context) AuthenticationRequiredResponse() error {
	message := "you must be authenticated to access this resource"
	return c.errorResponse(http.StatusUnauthorized, message)
}

func (c *Context) InactiveAccountResponse() error {
	message := "your user account must be activated to access this resource"
	return c.errorResponse(http.StatusForbidden, message)
}

func (c *Context) NotPermittedResponse() error {
	message := "your user account doesn't have the necessary permissions to access this resource"
	return c.errorResponse(http.StatusForbidden, message)
}

func (c *Context) RateLimitExceededResponse() error {
	message := "rate limit exceeded"
	return c.errorResponse(http.StatusTooManyRequests, message)
}
