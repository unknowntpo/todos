package api

import (
	//"errors"
	"net/http"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
	//"github.com/unknowntpo/todos/internal/helpers/validator"
	//"github.com/unknowntpo/todos/internal/logger"

	"github.com/julienschmidt/httprouter"
)

type userAPI struct {
	uu domain.UserUsecase
}

func NewUserAPI(router *httprouter.Router, uu domain.UserUsecase) {
	handler := &userAPI{uu: uu}

	router.HandlerFunc(http.MethodPost, "/v1/users", handler.RegisterUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", handler.ActivateUser)
}

func (u *userAPI) RegisterUser(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"user": "RegisterUser called"}, nil)

}

func (u *userAPI) ActivateUser(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"user": "ActivateUser called"}, nil)
}
