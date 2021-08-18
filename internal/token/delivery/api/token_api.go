package api

import (
	"net/http"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"

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
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"debug": "CreateAuthenticationToken called"}, nil)
}
