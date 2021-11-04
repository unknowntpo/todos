package postgres

import (
	"context"
	"database/sql"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
)

type tokenRepo struct {
	DB *sql.DB
}

func NewTokenRepo(DB *sql.DB) domain.TokenRepository {
	return &tokenRepo{DB}
}

// Insert adds the data for a specific token to the tokens table.
func (tr *tokenRepo) Insert(ctx context.Context, token *domain.Token) error {
	const op errors.Op = "tokenRepo.Insert"

	query := `
        INSERT INTO tokens (hash, user_id, expiry, scope) 
        VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	_, err := tr.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.E(op, errors.KindDatabase, err)
	}
	return nil
}

// DeleteAllForUser deletes all tokens for a specific user and scope.
func (tr *tokenRepo) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	const op errors.Op = "tokenRepo.DeleteAllForUser"

	query := `
        DELETE FROM tokens 
        WHERE scope = $1 AND user_id = $2`

	_, err := tr.DB.ExecContext(ctx, query, scope, userID)
	if err != nil {
		return errors.E(op, errors.KindDatabase, err)
	}
	return nil
}
