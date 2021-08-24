package domain

import (
	"context"
	"time"
)

// User represents an individual user.
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type UserUsecase interface {
	ValidatePasswordPlaintext(ctx context.Context, v Validator, password string)
}

type UserRepository interface {
	Insert(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*User, error)
}
