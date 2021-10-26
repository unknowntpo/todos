package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unknowntpo/todos/config"
	"github.com/unknowntpo/todos/internal/domain/mocks"
	"github.com/unknowntpo/todos/internal/logger/zerolog"
	"github.com/unknowntpo/todos/internal/reactor"

	"github.com/stretchr/testify/suite"
)

type MiddlewareTestSuite struct {
	suite.Suite
	mid *Middleware
	buf *bytes.Buffer
}

func (suite *MiddlewareTestSuite) SetupSuite() {
	// init new instance for Middleware
	// logger
	suite.buf = new(bytes.Buffer)
	logger := zerolog.New(suite.buf)

	rc := reactor.NewReactor(logger)

	config := new(config.Config)
	uu := new(mocks.UserUsecase)

	suite.mid = New(config, uu, rc)
}

func (suite *MiddlewareTestSuite) TearDownSuite() {
}

func (suite *MiddlewareTestSuite) SetupTest() {
}

func (suite *MiddlewareTestSuite) TearDownTest() {
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (suite *MiddlewareTestSuite) TestRecoverPanic() {
	// create dummy handler, with deliberated panic
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("deliberated panic")
	})

	// New request

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	// New Response recorder
	rr := httptest.NewRecorder()
	// call servehttp
	suite.mid.RecoverPanic(handler).ServeHTTP(rr, r)

	// check logger worked correctly.
	//suite.T().Log(suite.buf.String())
	// check error response is what we want.
	// check if resp header has Connection: Close.
	gotConnectionHeader := rr.Result().Header.Get("Connection")
	suite.Equal("close", gotConnectionHeader, "wrong Connection header")
	// Check if we got serverError response
	// FIXME: Why ServerError not shown in response ?
	suite.T().Log(rr.Body.String())
}
