package domain

import (
	"time"

	"github.com/unknowntpo/todos/internal/helpers/password"
)

// User represents an individual user.
type User struct {
	ID        int64             `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Password  password.Password `json:"-"`
	Activated bool              `json:"activated"`
	Version   int               `json:"-"`
}

type UserUsecase interface{}

type UserRepository interface{}
