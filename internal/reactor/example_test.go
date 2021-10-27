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

type hello struct {
	rc *Reactor
}

func NewHello(router *httprouter.Router, rc *Reactor) {
	hello := &hello{rc}
	router.HandlerFunc(http.MethodGet, "/", hello.helloHandler)
}

func (h *hello) helloHandler(w http.ResponseWriter, r *http.Request) {
	h.rc.WriteJSON(w, http.StatusOK, "Hello from helloHandler")
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

	NewHello(router, rc)

	router.ServeHTTP(rr, r)

	fmt.Println(rr.Body.String())
	// Output: "Hello from helloHandler"
}
