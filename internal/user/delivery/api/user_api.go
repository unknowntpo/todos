package api

import (
	"net/http"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/mailer"
	"github.com/unknowntpo/todos/pkg/naivepool"
	"github.com/unknowntpo/todos/pkg/validator"

	"github.com/unknowntpo/todos/internal/logger"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type userAPI struct {
	uu     domain.UserUsecase
	tu     domain.TokenUsecase
	pool   *naivepool.Pool
	mailer mailer.Mailer
	logger logger.Logger
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
	mailer mailer.Mailer,
	logger logger.Logger) {

	api := &userAPI{uu: uu, tu: tu, pool: pool, mailer: mailer, logger: logger}

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
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /v1/users/registration [post]
func (u *userAPI) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input UserRegistrationRequest

	// Parse the request body into the anonymous struct.
	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		helpers.BadRequestResponse(w, r, err)
		return
	}

	user := &domain.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if domain.ValidateUser(v, user); !v.Valid() {
		helpers.FailedValidationResponse(w, r, v.Errors)
		return
	}

	ctx := r.Context()

	// Insert the user data into the database.
	err = u.uu.Insert(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			helpers.FailedValidationResponse(w, r, v.Errors)
		default:
			u.logger.PrintError(err, nil)
			helpers.ServerErrorResponse(w, r, err)
		}
		return
	}

	// After the user record has been created in the database, generate a new activation
	// token for the user, and insert it to the database.
	token, err := domain.GenerateToken(user.ID, 3*24*time.Hour, domain.ScopeActivation)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	err = u.tu.Insert(ctx, token)
	if err != nil {
		err = errors.WithMessage(err, "failed.token.api.insert")
		u.logger.PrintError(err, nil)
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	u.pool.Schedule(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = u.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			u.logger.PrintError(err, nil)
		}
	})

	err = helpers.WriteJSON(w, http.StatusAccepted, &UserRegistrationResponse{User: user}, nil)
	if err != nil {
		err = errors.WithMessage(err, "failed to send json response. token.api")
		u.logger.PrintError(err, nil)

		helpers.ServerErrorResponse(w, r, err)
	}
}

// ActivateUser activates user based on given token.
// @Summary Activate user based on given token.
// @Description: None.
// @Produce  json
// @Param token query string true "activation token"
// @Success 200 {object} UserActivationResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /v1/users/activation [put]
func (u *userAPI) ActivateUser(w http.ResponseWriter, r *http.Request) {
	// Read token from request query string.
	tokenPlaintext := helpers.ReadString(r.URL.Query(), "token", "")

	v := validator.New()

	if domain.ValidateTokenPlaintext(v, tokenPlaintext); !v.Valid() {
		helpers.FailedValidationResponse(w, r, v.Errors)
		return
	}

	ctx := r.Context()

	user, err := u.uu.GetForToken(ctx, domain.ScopeActivation, tokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			helpers.FailedValidationResponse(w, r, v.Errors)
		default:
			helpers.ServerErrorResponse(w, r, err)
		}
		return
	}

	// Update the user's activation status.
	user.Activated = true

	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = u.uu.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEditConflict):
			helpers.EditConflictResponse(w, r)
		default:
			helpers.ServerErrorResponse(w, r, err)
		}
		return
	}

	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = u.tu.DeleteAllForUser(ctx, domain.ScopeActivation, user.ID)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	// Send the updated user details to the client in a JSON response.
	err = helpers.WriteJSON(w, http.StatusOK, &UserActivationResponse{User: user}, nil)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
	}
}
