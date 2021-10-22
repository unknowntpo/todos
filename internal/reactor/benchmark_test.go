package reactor

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unknowntpo/todos/internal/logger/zerolog"

	"github.com/julienschmidt/httprouter"
)

func BenchmarkHandlerWrapper(b *testing.B) {
	h := func(c *Context) error {
		return nil
	}

	null := ioutil.Discard
	logger := zerolog.New(null)

	rc := NewReactor(logger)

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		b.Fatalf("failed to create new request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := httprouter.New()

	// attach our handler with rc.HandlerWrapper
	router.Handler(http.MethodGet, "/", rc.HandlerWrapper(h))

	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rr, r)
	}
	b.StopTimer()
}
