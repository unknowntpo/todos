package postgres

import "fmt"

type tokenRepo struct {
	DB *sql.DB
}

func NewTokenRepo(DB *sql.DB) domain.TokenRepository {
	return &tokenRepo{DB}
}

func (tr *tokenRepo) New(ctx context.Context, userID int64, ttl time.Duration, scope string) (*Token, error) {
	return nil, nil
}
func (tr *tokenRepo) Insert(ctx context.Context, token *Token) error {
	return nil
}
func (tr *tokenRepo) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	return nil
}
