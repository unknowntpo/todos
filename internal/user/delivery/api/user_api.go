package api

import (
	"net/http"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
	"github.com/unknowntpo/todos/internal/mailer"
	"github.com/unknowntpo/todos/pkg/naivepool"
	"github.com/unknowntpo/todos/pkg/validator"

	"github.com/unknowntpo/todos/internal/reactor"

	"github.com/julienschmidt/httprouter"
)

type userAPI struct {
	uu     domain.UserUsecase
	tu     domain.TokenUsecase
	pool   *naivepool.Pool
	mailer *mailer.Mailer
	rc     *reactor.Reactor
}

type UserRegistrationResponse struct {
	User *domain.User `json:"user"`
}

type UserRegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserActivationResponse struct {
	User *domain.User `json:"user"`
}

func NewUserAPI(router *httprouter.Router,
	uu domain.UserUsecase,
	tu domain.TokenUsecase,
	pool *naivepool.Pool,
	mailer *mailer.Mailer, rc *reactor.Reactor) {

	api := &userAPI{uu: uu, tu: tu, pool: pool, mailer: mailer, rc: rc}

	router.HandlerFunc(http.MethodPost, "/v1/users/registration", api.RegisterUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/activation", api.ActivateUser)
}

// RegisterUser registers user based on given information.
// @Summary Register user based on given information.
// @Description: None.
// @Accept  json
// @Produce  json
// @Param reqBody body UserRegistrationRequest true "request body"
// @Success 202 {object} UserRegistrationResponse
// @Failure 400 {object} reactor.ErrorResponse
// @Failure 404 {object} reactor.ErrorResponse
// @Failure 500 {object} reactor.ErrorResponse
// @Router /v1/users/registration [post]
func (u *userAPI) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input UserRegistrationRequest

	// Parse the request body into the anonymous struct.
	err := u.rc.ReadJSON(w, r, &input)
	if err != nil {
		u.rc.BadRequestResponse(w, r, err)
		return
	}

	user := &domain.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		u.rc.ServerErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if domain.ValidateUser(v, user); !v.Valid() {
		u.rc.FailedValidationResponse(w, r, v.Err())
		return
	}

	ctx := r.Context()

	// Insert the user data into the database.
	err = u.uu.Insert(ctx, user)
	if err != nil {
		switch {
		case errors.KindIs(err, errors.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			u.rc.FailedValidationResponse(w, r, v.Err())
			return
		default:
			u.rc.ServerErrorResponse(w, r, err)
			return
		}
	}

	// After the user record has been created in the database, generate a new activation
	// token for the user, and insert it to the database.
	token, err := domain.GenerateToken(user.ID, 3*24*time.Hour, domain.ScopeActivation)
	if err != nil {
		u.rc.ServerErrorResponse(w, r, err)
		return
	}

	err = u.tu.Insert(ctx, token)
	if err != nil {
		u.rc.ServerErrorResponse(w, r, err)
		return
	}

	u.pool.Schedule(func() {
		const op errors.Op = "userAPI.RegisterUser"

		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = u.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			u.rc.Logger.PrintError(
				errors.E(
					op,
					errors.UserEmail(user.Email),
					errors.ErrInternal,
					errors.Msg("failed to send welcome email"),
					err,
				),
				nil,
			)
			return
		}
	})

	err = u.rc.WriteJSON(w, http.StatusAccepted, &UserRegistrationResponse{User: user})
	if err != nil {
		u.rc.ServerErrorResponse(w, r, err)
	}
}

// ActivateUser activates user based on given token.
// @Summary Activate user based on given token.
// @Description: None.
// @Produce  json
// @Param token query string true "activation token"
// @Success 200 {object} UserActivationResponse
// @Failure 400 {object} reactor.ErrorResponse
// @Failure 404 {object} reactor.ErrorResponse
// @Failure 500 {object} reactor.ErrorResponse
// @Router /v1/users/activation [put]
func (u *userAPI) ActivateUser(w http.ResponseWriter, r *http.Request) {
	// Read token from request query string.
	qs := r.URL.Query()
	tokenPlaintext := u.rc.ReadString(qs, "token", "")

	v := validator.New()

	if domain.ValidateTokenPlaintext(v, tokenPlaintext); !v.Valid() {
		u.rc.FailedValidationResponse(w, r, v.Err())
		return
	}

	ctx := r.Context()

	user, err := u.uu.Authenticate(ctx, domain.ScopeActivation, tokenPlaintext)
	if err != nil {
		switch {
		case errors.KindIs(err, errors.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			u.rc.FailedValidationResponse(w, r, v.Err())
			return
		default:
			u.rc.ServerErrorResponse(w, r, err)
			return
		}
	}

	// Update the user's activation status.
	user.Activated = true

	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = u.uu.Update(ctx, user)
	if err != nil {
		switch {
		case errors.KindIs(err, errors.ErrEditConflict):
			u.rc.EditConflictResponse(w, r)
			return
		default:
			u.rc.ServerErrorResponse(w, r, err)
			return
		}
	}

	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = u.tu.DeleteAllForUser(ctx, domain.ScopeActivation, user.ID)
	if err != nil {
		u.rc.ServerErrorResponse(w, r, err)
		return
	}

	// Send the updated user details to the client in a JSON response.
	err = u.rc.WriteJSON(w, http.StatusOK, &UserActivationResponse{User: user})
	if err != nil {
		u.rc.ServerErrorResponse(w, r, err)
	}
}
