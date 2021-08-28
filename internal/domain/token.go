package domain

import (
	"context"
	"time"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenUsecase interface {
	New(ctx context.Context, userID int64, ttl time.Duration, scope string) (*Token, error)
	Insert(ctx context.Context, token *Token) error
	DeleteAllForUser(ctx context.Context, scope string, userID int64) error
}

type TokenRepository interface {
	New(ctx context.Context, userID int64, ttl time.Duration, scope string) (*Token, error)
	Insert(ctx context.Context, token *Token) error
	DeleteAllForUser(ctx context.Context, scope string, userID int64) error
}