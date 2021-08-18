package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
)

type userRepo struct {
	DB *sql.DB
}

func NewUserRepo() domain.UserRepository {
	return &userRepo{}
}

func (ur *userRepo) Insert(ctx context.Context, user *domain.User) error {
	return nil
}

func (ur *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (ur *userRepo) Update(ctx context.Context, user *domain.User) error {
	return nil
}

func (ur *userRepo) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	return nil, nil
}
