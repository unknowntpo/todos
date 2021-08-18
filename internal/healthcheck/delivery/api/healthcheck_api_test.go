package api

import (
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/unknowntpo/todos/internal/testutil"
)

func TestHealthcheckHandler(t *testing.T) {
	router := httprouter.New()

	// mock the dependency of HealthcheckAPI
	NewHealthcheckAPI(router, "1.0.0", "development")

	wantStatus := http.StatusOK
	wantBody := `{
	"status": "available",
	"system_info": {
		"environment": "development",
		"version": "1.0.0"
	}
}
`

	testutil.TestEndpoint(t, router, testutil.APITestCase{
		Name:       "Get healthcheck status",
		Method:     "GET",
		URL:        "/v1/healthcheck",
		Body:       "",
		Header:     nil,
		WantStatus: wantStatus,
		WantBody:   wantBody,
	})
	/*
			router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handler.HealthcheckHandler)

			// Set request to /v1/healthcheck
			r, _ := http.NewRequest("GET", "/v1/healthcheck", nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, r)

			assert.Equal(t, http.StatusOK, rr.Code, "Wrong status code")

			wantBody := `{
			"status": "available",
			"system_info": {
				"environment": "development",
				"version": "1.0.0"
			}
		}
		`
			assert.Equal(t, wantBody, rr.Body.String(), "Wrong response body")
	*/
}
