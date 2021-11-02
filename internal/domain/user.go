package domain

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
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

// AnonymousUser represents an anonymous user.
var AnonymousUser = &User{}

// IsAnonymous checks if a User instance is the AnonymousUser.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type UserUsecase interface {
	Insert(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Register(ctx context.Context, user *User) error
	Activate(ctx context.Context, tokenPlaintext string) (*User, error)
	Login(ctx context.Context, email, password string) (*Token, error)
	Authenticate(ctx context.Context, tokenScope, tokenPlaintext string) (*User, error)
}

type UserRepository interface {
	Insert(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*User, error)
}
