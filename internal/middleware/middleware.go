package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/logger"
)

type generalMiddleware struct {
}

func New() *generalMiddleware {
	return nil

}

func (mid *generalMiddleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				// TODO: How to handle panic?
				// Try to log stack trace ?
				logger.Log.PrintError(errors.New(fmt.Sprintf("%s", err)), nil)
				helpers.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
