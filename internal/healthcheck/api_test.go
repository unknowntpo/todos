package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	handler := http.HandlerFunc(healthcheckHandler)
	r := httptest.NewRequest("GET", "/healthcheck", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, r)
	got := rr.Body.String()
	want := "OK"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
