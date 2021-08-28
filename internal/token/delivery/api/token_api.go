package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/helpers/validator"

	"github.com/julienschmidt/httprouter"
)

type tokenAPI struct {
	TU domain.TokenUsecase
	UU domain.UserUsecase
}

func NewTokenAPI(router *httprouter.Router, tu domain.TokenUsecase, uu domain.UserUsecase) {
	api := &tokenAPI{TU: tu, UU: uu}
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", api.CreateAuthenticationToken)
}

type CreateAuthTokenResponse struct {
	AuthenticationToken *domain.Token `json: authentication_token`
}

func (t *tokenAPI) CreateAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	//helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"debug": "CreateAuthenticationToken called"}, nil)

	// Parse the email and password from the request body.
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		helpers.BadRequestResponse(w, r, err)
		return
	}

	// Validate the email and password provided by the client.
	v := validator.New()

	domain.ValidateEmail(v, input.Email)
	domain.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		helpers.FailedValidationResponse(w, r, v.Errors)
		return
	}

	ctx := r.Context()

	user, err := t.UU.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			helpers.InvalidCredentialsResponse(w, r)
		default:
			helpers.ServerErrorResponse(w, r, err)
		}
		return
	}

	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	if !match {
		helpers.InvalidCredentialsResponse(w, r)
		return
	}

	// Otherwise, if the password is correct, we generate a new token with a 24-hour
	// expiry time and the scope 'authentication'.
	token, err := t.TU.New(ctx, user.ID, 24*time.Hour, domain.ScopeAuthentication)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	// Encode the token to JSON and send it in the response along with a 201 Created
	// status code.
	err = helpers.WriteJSON(w,
		http.StatusCreated,
		&CreateAuthTokenResponse{
			AuthenticationToken: token,
		}, nil)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
	}

}
