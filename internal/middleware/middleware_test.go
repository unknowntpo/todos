package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unknowntpo/todos/config"
	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/mocks"
	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/logger/zerolog"
	"github.com/unknowntpo/todos/internal/reactor"
	"github.com/unknowntpo/todos/internal/testutil"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/suite"
)

type MiddlewareTestSuite struct {
	suite.Suite
	mid     *Middleware
	usecase *mocks.UserUsecase
	config  *config.Config
	rc      *reactor.Reactor
	logBuf  *bytes.Buffer
}

func (suite *MiddlewareTestSuite) SetupSuite() {
}

func (suite *MiddlewareTestSuite) TearDownSuite() {
}

func (suite *MiddlewareTestSuite) SetupTest() {
	suite.logBuf = new(bytes.Buffer)
	logger := zerolog.New(suite.logBuf)

	suite.rc = reactor.NewReactor(logger)

	suite.config = new(config.Config)
	suite.usecase = new(mocks.UserUsecase)

	suite.mid = New(suite.config, suite.usecase, suite.rc)
}

func (suite *MiddlewareTestSuite) TearDownTest() {
	suite.mid = nil
	suite.usecase = nil
	suite.config = nil
	suite.rc = nil
	suite.logBuf = nil
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
	suite.Contains(rr.Body.String(), `"error": "the server encountered a problem and could not process your request"`, "response body should contain internal server error message")
	suite.Contains(suite.logBuf.String(), "deliberated panic", "logger should contain panic message")
}

func (suite *MiddlewareTestSuite) TestRequireAuthenticatedUser() {
	// When the user are not authenticated,
	// RequireAuthenticatedUser should reject the request,
	// and the user should get the error message.
	suite.Run("anonymous user should be rejected", func() {
		suite.TearDownTest()
		suite.SetupTest()
		router := httprouter.New()

		h := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		}

		router.Handler(http.MethodGet, "/", suite.mid.RequireAuthenticatedUser(http.HandlerFunc(h)))

		r, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			suite.T().Fatal("unable to create new request")
		}

		r = helpers.ContextSetUser(r, domain.AnonymousUser)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, r)

		// Because we don't treat this as an internal error.
		suite.Equal("", suite.logBuf.String())
		suite.Contains(rr.Body.String(), `"error": "you must be authenticated to access this resource"`)
		suite.TearDownTest()
	})

	suite.Run("authenticated user should be accepted", func() {
		suite.TearDownTest()
		suite.SetupTest()
		router := httprouter.New()

		h := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		}

		router.Handler(http.MethodGet, "/", suite.mid.RequireAuthenticatedUser(http.HandlerFunc(h)))

		r, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			suite.T().Fatal("unable to create new request")
		}

		fakeUser := testutil.NewFakeUser(suite.T(), "Alice Smith", "alice@example.com", "pa55word", true)

		r = helpers.ContextSetUser(r, fakeUser)

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, r)

		// Because we don't treat this as an internal error.
		suite.Equal("", suite.logBuf.String())
		suite.Equal("OK", rr.Body.String())
		suite.TearDownTest()
	})

}
