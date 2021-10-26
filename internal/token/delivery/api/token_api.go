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
	router.Handler(http.MethodPost, "/v1/tokens/authentication", rc.HandlerWrapper(api.CreateAuthenticationToken))
}

// @Summary Create authentication token for user.
// @Description None.
// @Accept  json
// @Produce  json
// @Param authentication_request_body body AuthenticationRequestBody true "authentication request body"
// @Success 200 {object} AuthenticationResponse
// @Router /v1/tokens/authentication [post]
func (t *tokenAPI) CreateAuthenticationToken(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "tokenAPI.CreateAuthenticationToken"
	// Parse the email and password from the request body.
	var input AuthenticationRequestBody

	err := reactor.ReadJSON(w, r, &input)
	if err != nil {
		return reactor.BadRequestResponse(w, r, err)
	}

	// Validate the email and password provided by the client.
	v := validator.New()

	domain.ValidateEmail(v, input.Email)
	domain.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		return reactor.FailedValidationResponse(w, r, v.Err())
	}

	ctx := r.Context()

	user, err := t.UU.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.KindIs(err, errors.ErrRecordNotFound):
			return reactor.InvalidCredentialsResponse(w, r)
		default:
			return errors.E(op, errors.Msg("failed to get user by email"), err)
		}
	}

	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		return errors.E(op, err)
	}

	if !match {
		return reactor.InvalidCredentialsResponse(w, r)
	}

	// Otherwise, if the password is correct, we generate a new token with a 24-hour
	// expiry time and the scope 'authentication'.
	token, err := domain.GenerateToken(user.ID, 24*time.Hour, domain.ScopeAuthentication)
	if err != nil {
		return errors.E(op, err)
	}

	err = t.TU.Insert(ctx, token)
	if err != nil {
		return errors.E(op, err)
	}

	// Encode the token to JSON and send it in the response along with a 201 Created
	// status code.
	return reactor.WriteJSON(w, http.StatusCreated, &AuthenticationResponse{Token: token})
}
