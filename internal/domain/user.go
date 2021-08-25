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
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type UserUsecase interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	Insert(ctx context.Context, user *User) error
	GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type UserRepository interface {
	Insert(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*User, error)
}
