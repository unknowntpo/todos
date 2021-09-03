package api

import (
	//"errors"
	"net/http"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/helpers/background"

	//"github.com/unknowntpo/todos/internal/helpers/validator"
	//"github.com/unknowntpo/todos/internal/logger"

	"github.com/julienschmidt/httprouter"
)

type userAPI struct {
	bg *background.Background
	uu domain.UserUsecase
}

func NewUserAPI(router *httprouter.Router, uu domain.UserUsecase, bg *background.Background) {
	api := &userAPI{uu: uu, bg: bg}

	router.HandlerFunc(http.MethodPost, "/v1/users/registration", api.RegisterUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/activation", api.ActivateUser)
}

func (u *userAPI) RegisterUser(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"user": "RegisterUser called"}, nil)
}

func (u *userAPI) ActivateUser(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"user": "ActivateUser called"}, nil)
}
