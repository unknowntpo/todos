package reactor

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/logger"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

type Reactor struct {
	Logger logger.Logger
}

func NewReactor(logger logger.Logger) *Reactor {
	return &Reactor{Logger: logger}
}

// HandlerFunc allow us to use a function with signature HandlerFunc as our actual handler,
// which can simplify error handling.
// Now, we can do error handling inside the http.Handler that HandlerWrapper returns, and put
// some dependencies inside Reactor, e.g. logger.
func (rc *Reactor) HandlerWrapper(h HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		// Only Internal Server Error will come to here.
		// TODO: Use kindIs to check , if failed, panic
		if err != nil {
			rc.Logger.PrintError(err, nil)
			if e := ServerErrorResponse(w, r); e != nil {
				// Something goes wrong during sending server error response.
				// So we write the message directly.
				rc.Logger.PrintError(e, nil)
				w.WriteHeader(http.StatusInternalServerError)
				msg := `{"error":"the server encountered a problem and could not process your request"}`
				w.Write([]byte(msg))
				return
			}
			return
		}
	})
}
