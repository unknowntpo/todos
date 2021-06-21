package testutil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

// APITestCase represents the data needed to describe an API test case.
type APITestCase struct {
	Name        string
	Method, URL string
	Body        string
	Header      http.Header
	WantStatus  int
	WantBody    string
}

// Endpoint tests an HTTP endpoint using the given APITestCase spec.
func TestEndpoint(t *testing.T, router *httprouter.Router, tc APITestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		r, _ := http.NewRequest(tc.Method, tc.URL, bytes.NewBufferString(tc.Body))
		if tc.Header != nil {
			r.Header = tc.Header
		}
		if r.Header.Get("Content-Type") == "" {
			r.Header.Set("Content-Type", "application/json")
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, r)
		assert.Equal(t, tc.WantStatus, rr.Code, "Wrong status code")
		if tc.WantBody != "" {
			assert.Equal(t, tc.WantBody, rr.Body.String(), "Wrong response body")
		}
	})
}
