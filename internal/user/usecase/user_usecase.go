package main

import (
	"context"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers/validator"
)

type userUsecase struct{}

func NewUserUsecase() domain.UserUsecase {
	return &userUsecase{}
}

func (m *userUsecase) ValidatePasswordPlaintext(ctx context.Context, v *validator.Validator, password string) {
	return
}
