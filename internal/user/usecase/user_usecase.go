package usecase

import (
	"context"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	tokenRepo      domain.TokenRepository
	contextTimeout time.Duration
}

func NewUserUsecase(ur domain.UserRepository, tr domain.TokenRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{userRepo: ur, tokenRepo: tr, contextTimeout: timeout}
}

func (uu *userUsecase) Insert(ctx context.Context, user *domain.User) error {
	const op errors.Op = "userUsecase.Insert"

	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	err := uu.userRepo.Insert(ctx, user)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// Authenticate authenticates the user with given tokenScope and tokenPlaintext.
// If succeed, returns user object and nil, if not return nil and the error.
func (uu *userUsecase) Authenticate(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	const op errors.Op = "userUsecase.Authenticate"

	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.userRepo.GetForToken(ctx, tokenScope, tokenPlaintext)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return user, nil
}

func (uu *userUsecase) Update(ctx context.Context, user *domain.User) error {
	const op errors.Op = "userUsecase.Update"

	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	err := uu.userRepo.Update(ctx, user)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// Login performs login operation and returns token if succeed,
// if failed, it returns nil and errors.ErrInvalidCredentials error,
// if some internal server error happened, returns nil and wrapped error.
func (uu *userUsecase) Login(ctx context.Context, email, password string) (*domain.Token, error) {
	const op errors.Op = "userUsecase.Login"

	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.E(op, err)
	}

	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(password)
	if err != nil {
		return nil, errors.E(op, err)
	}

	if !match {
		return nil, errors.E(op, errors.ErrInvalidCredentials, err)
	}

	// Otherwise, if the password is correct, we generate a new token with a 24-hour
	// expiry time and the scope 'authentication'.
	// At here, user is valid!
	token, err := domain.GenerateToken(user.ID, 24*time.Hour, domain.ScopeAuthentication)
	if err != nil {
		return nil, errors.E(op, err)
	}

	err = uu.tokenRepo.Insert(ctx, token)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return token, nil
}
