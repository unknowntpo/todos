package api

import (
	"net/http"
	"testing"

	"github.com/unknowntpo/todos/internal/domain/mock"
	"github.com/unknowntpo/todos/internal/logger/logrus"
	"github.com/unknowntpo/todos/internal/testutil"

	"github.com/julienschmidt/httprouter"
)

func TestGetByID(t *testing.T) {
	// newTestApplication
	// set up new logger.
	err := logrus.RegisterLog()
	if err != nil {
		t.Fatalf("Fail to set up test logger: %v", err)
	}

	router := httprouter.New()

	tu := mock.NewTaskUsecase()
	NewTaskAPI(router, tu)

	tcases := []testutil.APITestCase{
		{
			Name:       "Test Get by ID: StatusOK",
			Method:     "GET",
			URL:        "/v1/tasks/1",
			Body:       "",
			Header:     nil,
			WantStatus: http.StatusOK,
			WantBody: `{
	"task": {
		"id": 1,
		"title": "Do homework",
		"content": "Interesting",
		"done": true,
		"version": 1
	}
}
`,
		},
		{
			Name:       "Test Get by ID: StatusNotFound",
			Method:     "GET",
			URL:        "/v1/tasks/-1",
			Body:       "",
			Header:     nil,
			WantStatus: http.StatusNotFound,
			WantBody: `{
	"error": "the requested resource could not be found"
}
`,
		},
	}

	for _, tc := range tcases {
		testutil.TestEndpoint(t, router, tc)
	}
}

func TestGetByIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
