package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
)

type userRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) domain.UserRepository {
	return &userRepo{db}
}

// Insert inserts a new record in the database for the user.
// then write user's id, created time, version back to *domain.User specified
// in method parameter.
// If user is already in the database, return domain.ErrDuplicateEmail error.
func (ur *userRepo) Insert(ctx context.Context, user *domain.User) error {
	const op errors.Op = "userRepo.Insert"

	query := `
        INSERT INTO users (name, email, password_hash, activated) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

	args := []interface{}{user.Name, user.Email, user.Password.Hash, user.Activated}

	err := ur.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			// Just send our custom error type, we don't need error message from database, because at here,
			// we don't treat err as error.
			return errors.E(op, errors.ErrDuplicateEmail)
		default:
			// ErrDatabase is a subset internal server error.
			return errors.E(op, errors.ErrDatabase, err)
		}
	}

	return nil
}

// GetByEmail gets the User details from the database based on the user's email address.
// If there's no record, domain.ErrRecordNotFound will be returned.
func (ur *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const op errors.Op = "userRepo.GetByEmail"

	query := `
        SELECT id, created_at, name, email, password_hash, activated, version
        FROM users
        WHERE email = $1`

	var user domain.User

	err := ur.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.E(op, errors.ErrRecordNotFound)
		default:
			return nil, errors.E(op, errors.ErrDatabase, err)
		}
	}

	return &user, nil
}

// Update updates the details for a specific user.
// Notice that we check against the version field to help prevent
// any race conditions during the request cycle.
// And we also check for a violation of the "users_email_key"
// constraint when performing the update.
func (ur *userRepo) Update(ctx context.Context, user *domain.User) error {
	const op errors.Op = "userRepo.Update"

	query := `
        UPDATE users 
        SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
        WHERE id = $5 AND version = $6
        RETURNING version`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	err := ur.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return errors.E(op, errors.ErrDuplicateEmail)
		case errors.Is(err, sql.ErrNoRows):
			return errors.E(op, errors.ErrEditConflict)
		default:
			return errors.E(op, errors.ErrDatabase, err)
		}
	}

	return nil
}

// GetForToken returns the user corresponding to the given token.
func (ur *userRepo) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	const op errors.Op = "userRepo.GetForToken"

	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
        SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
        FROM users
        INNER JOIN tokens
        ON users.id = tokens.user_id
        WHERE tokens.hash = $1
        AND tokens.scope = $2 
        AND tokens.expiry > $3`

	// We use the [:] operator to get a slice containing the token hash,
	// rather than passing in the array (which is not supported by the pq driver).
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user domain.User

	err := ur.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.E(op, errors.ErrRecordNotFound)
		default:
			return nil, errors.E(op, errors.ErrDatabase, err)
		}
	}

	return &user, nil
}
