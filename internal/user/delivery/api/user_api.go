package api

import (
	"errors"
	"net/http"
	"time"

	swagger "github.com/unknowntpo/todos/docs/go"
	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/helpers/mailer"
	"github.com/unknowntpo/todos/internal/helpers/validator"
	"github.com/unknowntpo/todos/internal/helpers/workerpool"

	"github.com/unknowntpo/todos/internal/logger"

	"github.com/julienschmidt/httprouter"
)

type userAPI struct {
	wp     *workerpool.Pool
	uu     domain.UserUsecase
	mailer mailer.Mailer
}

func NewUserAPI(router *httprouter.Router, uu domain.UserUsecase, wp *workerpool.Pool, mailer mailer.Mailer) {
	api := &userAPI{uu: uu, wp: wp, mailer: mailer}

	router.HandlerFunc(http.MethodPost, "/v1/users/registration", api.RegisterUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/activation", api.ActivateUser)
}

func (u *userAPI) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input swagger.UserRegistrationRequest

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
			helpers.ServerErrorResponse(w, r, err)
		}
		return
	}

	// After the user record has been created in the database, generate a new activation
	// token for the user.
	token, err := domain.GenerateToken(user.ID, 3*24*time.Hour, domain.ScopeActivation)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	u.wp.Schedule(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = u.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			logger.Log.PrintError(err, nil)
		}
	})

	// FIXME: cannot use user (type *domain.User) as type *swagger.User in field value.
	err = helpers.WriteJSON(w, http.StatusAccepted, &swagger.UserRegistrationResponse{User: user}, nil)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
	}
}

func (u *userAPI) ActivateUser(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"user": "ActivateUser called"}, nil)
}
