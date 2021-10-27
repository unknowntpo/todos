package api

import (
	"net/http"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
	"github.com/unknowntpo/todos/internal/reactor"
	"github.com/unknowntpo/todos/pkg/validator"

	"github.com/julienschmidt/httprouter"
)

type tokenAPI struct {
	TU domain.TokenUsecase
	UU domain.UserUsecase
	rc *reactor.Reactor
}

type AuthenticationRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticationResponse struct {
	Token *domain.Token `json:"token"`
}

func NewTokenAPI(router *httprouter.Router, tu domain.TokenUsecase, uu domain.UserUsecase, rc *reactor.Reactor) {
	api := &tokenAPI{TU: tu, UU: uu, rc: rc}
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", api.CreateAuthenticationToken)
}

// @Summary Create authentication token for user.
// @Description None.
// @Accept  json
// @Produce  json
// @Param authentication_request_body body AuthenticationRequestBody true "authentication request body"
// @Success 200 {object} AuthenticationResponse
// @Router /v1/tokens/authentication [post]
func (t *tokenAPI) CreateAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	// Parse the email and password from the request body.
	var input AuthenticationRequestBody

	err := t.rc.ReadJSON(w, r, &input)
	if err != nil {
		t.rc.BadRequestResponse(w, r, err)
		return
	}

	// Validate the email and password provided by the client.
	v := validator.New()

	domain.ValidateEmail(v, input.Email)
	domain.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		t.rc.FailedValidationResponse(w, r, v.Err())
		return
	}

	ctx := r.Context()

	user, err := t.UU.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.KindIs(err, errors.ErrRecordNotFound):
			t.rc.InvalidCredentialsResponse(w, r)
			return
		default:
			t.rc.ServerErrorResponse(w, r, err)
			return
		}
	}

	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		t.rc.ServerErrorResponse(w, r, err)
		return
	}

	if !match {
		t.rc.InvalidCredentialsResponse(w, r)
		return
	}

	// Otherwise, if the password is correct, we generate a new token with a 24-hour
	// expiry time and the scope 'authentication'.
	token, err := domain.GenerateToken(user.ID, 24*time.Hour, domain.ScopeAuthentication)
	if err != nil {
		t.rc.ServerErrorResponse(w, r, err)
		return
	}

	err = t.TU.Insert(ctx, token)
	if err != nil {
		t.rc.ServerErrorResponse(w, r, err)
		return
	}

	// Encode the token to JSON and send it in the response along with a 201 Created
	// status code.
	err = t.rc.WriteJSON(w, http.StatusCreated, &AuthenticationResponse{Token: token})
	if err != nil {
		t.rc.ServerErrorResponse(w, r, err)
		return
	}
}
