package reactor

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/unknowntpo/todos/internal/logger/zerolog"

	"github.com/julienschmidt/httprouter"
)

func helloHandler(c *Context) error {
	return c.WriteJSON(http.StatusOK, "Hello from helloHandler")
}

func Example_httprouter() {
	logBuf := bytes.NewBufferString("")

	logger := zerolog.New(logBuf)

	rc := NewReactor(logger)

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		fmt.Printf("failed to create new request: %v", err)
		os.Exit(1)
	}

	router := httprouter.New()

	// attach our handler with rc.HandlerWrapper
	router.Handler(http.MethodGet, "/", rc.HandlerWrapper(helloHandler))

	router.ServeHTTP(rr, r)

	fmt.Println(rr.Body.String())
	// Output: "Hello from helloHandler"
}
