package main

import (
	"context"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers/validator"
)

type tokenUsecase struct {
	tr             domain.TokenRepository
	contextTimeout time.Duration
}

func (tu *tokenUsecase) ValidateTokenPlaintext(ctx context.Context, v *validator.Validator, tokenPlaintext string) {
	return
}
