package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
)

type tokenRepo struct {
	DB *sql.DB
}

func NewTokenRepo(DB *sql.DB) domain.TokenRepository {
	return &tokenRepo{DB}
}

// New is a shortcut which creates a new Token struct and then inserts the
// data in the tokens table.
func (tr *tokenRepo) New(ctx context.Context, userID int64, ttl time.Duration, scope string) (*domain.Token, error) {
	token, err := domain.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = tr.Insert(ctx, token)
	return token, err
}

// Insert adds the data for a specific token to the tokens table.
func (tr *tokenRepo) Insert(ctx context.Context, token *domain.Token) error {
	query := `
        INSERT INTO tokens (hash, user_id, expiry, scope) 
        VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tr.DB.ExecContext(ctx, query, args...)
	return err
}

// DeleteAllForUser deletes all tokens for a specific user and scope.
func (tr *tokenRepo) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	query := `
        DELETE FROM tokens 
        WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tr.DB.ExecContext(ctx, query, scope, userID)
	return err
}
