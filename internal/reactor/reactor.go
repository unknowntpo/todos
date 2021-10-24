package reactor

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/domain/errors"
	"github.com/unknowntpo/todos/internal/logger"
)

type HandlerFunc func(c *Context) error

// TODO: What about middleware ?

// What we Want
/*

httprouter.Handler(method, path, r.HandlerWrapper(api.Healthcheck))
*/
type Reactor struct {
	logger logger.Logger
}

func NewReactor(logger logger.Logger) *Reactor {
	return &Reactor{logger: logger}
}

// HandlerFunc allow us to use a function with signature HandlerFunc as our actual handler,
// which can simplify error handling.
// Now, we can do error handling inside the http.Handler that HandlerWrapper returns, and put
// some dependencies inside Reactor.
// Example usage:
/*
func (h *healthcheckAPI) Healthcheck(c *reactor.Context) error {
	err :=c.WriteJSON(w, http.StatusOK, &HealthcheckResponse{
		Status:      "available",
		Version:     h.version,
		Environment: h.env,
	})
	if err != nil {
		// Maybe using c.WriteJSON ?
		c.ServerErrorResponse(w, r, err)
	}
}
*/
func (rc *Reactor) HandlerWrapper(h HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use ctxPool to reduce allocation and gc.
		c := ctxPool.Get().(*Context)
		defer ctxPool.Put(c)
		c.W = w
		c.R = r

		err := h(c)
		// Only Internal Server Error will come to here.
		// TODO: Use kindIs to check , if failed, panic
		if err != nil {
			c.ServerErrorResponse(err)
			return
		}
	})
}
