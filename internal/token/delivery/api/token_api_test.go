package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/mocks"
	"github.com/unknowntpo/todos/internal/logger/zerolog"
	"github.com/unknowntpo/todos/internal/reactor"
	"github.com/unknowntpo/todos/internal/testutil"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAuthenticationToken(t *testing.T) {
	// Success
	t.Run("Success", func(t *testing.T) {
		logBuf := new(bytes.Buffer)
		logger := zerolog.New(logBuf)

		userUsecase := new(mocks.UserUsecase)

		fakeUser := testutil.NewFakeUser(t, "Alice Smith", "alice@example.com", "pa55word", true)
		// TODO: Need token.New() usecase helper function
		wantToken, err := domain.GenerateToken(1, 30*time.Minute, domain.ScopeActivation)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		userUsecase.On("Login", mock.Anything, fakeUser.Email, *fakeUser.Password.Plaintext).Return(wantToken, nil)

		rc := reactor.NewReactor(logger)

		reqBody := new(bytes.Buffer)
		reqBody.WriteString(`{"email": "alice@example.com", "password": "pa55word"}`)
		r, err := http.NewRequest(http.MethodPost, "/v1/tokens/authentication", reqBody)
		if err != nil {
			t.Fatalf("failed to create new request: %v", err)
		}

		rr := httptest.NewRecorder()
		// deps:
		// tu, uu, rc
		router := httprouter.New()
		NewTokenAPI(router, userUsecase, rc)

		router.ServeHTTP(rr, r)
		t.Log(logBuf.String())
		assert.Equal(t, "", logBuf.String())
		assert.Contains(t, rr.Body.String(), wantToken.Plaintext)
		// TODO: How to specify the exact format we want ?
		assert.Contains(t, rr.Body.String(), "expiry")
	})
	t.Run("Fail on invalid credentials", func(t *testing.T) {
		// deps:
		// tu, uu, rc

	})

	// InvalidCredentials
	// rr
	// r
	//
	//	t.Skip("TODO: finish the implementation")
}
