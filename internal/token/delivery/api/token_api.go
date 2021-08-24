package api

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/helpers/validator"

	"github.com/julienschmidt/httprouter"
)

type tokenAPI struct {
	TU domain.TokenUsecase
}

func NewTokenAPI(router *httprouter.Router, tu domain.TokenUsecase) {
	api := &tokenAPI{TU: tu}
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", api.CreateAuthenticationToken)
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

	helpers.ValidateEmail(v, input.Email)
	helpers.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		helpers.FailedValidationResponse(w, r, v.Errors)
		return
	}

	// TODO: Move code below to proper position.
	/*
		// Lookup the user record based on the email address. If no matching user was
		// found, then we call the app.invalidCredentialsResponse() helper to send a 401
		// Unauthorized response to the client (we will create this helper in a moment).
		user, err := app.models.Users.GetByEmail(input.Email)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidCredentialsResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		// Check if the provided password matches the actual password for the user.
		match, err := user.Password.Matches(input.Password)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// If the passwords don't match, then we call the app.invalidCredentialsResponse()
		// helper again and return.
		if !match {
			app.invalidCredentialsResponse(w, r)
			return
		}

		// Otherwise, if the password is correct, we generate a new token with a 24-hour
		// expiry time and the scope 'authentication'.
		token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Encode the token to JSON and send it in the response along with a 201 Created
		// status code.
		err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}

	*/
}
