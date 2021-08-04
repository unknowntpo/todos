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

type UserAPI struct {
	UU domain.UserUsecase
}

func NewUserAPI(router *httprouter.Router, uu domain.UserUsecase) {
	handler := &UserAPI{UU: uu}

	router.HandlerFunc(http.MethodPost, "/v1/users", handler.RegisterUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", handler.ActivateUser)
}

func (u *UserAPI) RegisterUser(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"user": "RegisterUser called"}, nil)

}

func (u *UserAPI) ActivateUser(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"user": "ActivateUser called"}, nil)
}
