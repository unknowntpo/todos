package main

import (
	"github.com/unknowntpo/todos/internal/domain"
)

type tokenUsecase struct {
	tr             domain.TokenRepository
	contextTimeout time.Duration
}

func (tu *tokenUsecase) ValidateTokenPlaintext(ctx context.Context, v *validator.Validator, tokenPlaintext string) {
	return
}
