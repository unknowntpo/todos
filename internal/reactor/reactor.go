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
		// TODO: Use sync.Pool to avoid repeatly allocation.
		c := &Context{w: w, r: r}
		err := h(c)
		if err != nil {
			rc.logger.PrintError(err, nil)
			/*
				// TODO: Use some variable to store this message.
				message := "the server encountered a problem and could not process your request"
				c.WriteJSON(http.StatusInternalServerError, map[string]interface{}{
					"error": message,
				})
			*/
			if e := c.WriteJSON(
				http.StatusInternalServerError,
				errors.NewErrorResponse(w, r, errors.InternalServerErrorResponse, err),
			); e != nil {
				rc.logger.PrintError(e, nil)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})
}
