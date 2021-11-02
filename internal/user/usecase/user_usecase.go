package usecase

import (
	"context"
	"time"

	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/domain/errors"
	"github.com/unknowntpo/todos/internal/logger"
	"github.com/unknowntpo/todos/internal/mailer"
	"github.com/unknowntpo/todos/pkg/naivepool"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	tokenUsecase   domain.TokenUsecase
	pool           *naivepool.Pool
	mailer         *mailer.Mailer
	logger         logger.Logger
	contextTimeout time.Duration
}

func NewUserUsecase(
	ur domain.UserRepository,
	tu domain.TokenUsecase,
	p *naivepool.Pool,
	mailer *mailer.Mailer,
	logger logger.Logger,
	timeout time.Duration,
) domain.UserUsecase {
	return &userUsecase{
		userRepo:       ur,
		tokenUsecase:   tu,
		pool:           p,
		mailer:         mailer,
		logger:         logger,
		contextTimeout: timeout,
	}
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

	err = uu.tokenUsecase.Insert(ctx, token)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return token, nil
}

// Activate performs user activation and returns user, nil if succeed,
// if failed, it returns nil and errors.ErrInvalidCredentials error,
// if some internal server error happened, returns nil and wrapped error.
func (uu *userUsecase) Activate(ctx context.Context, tokenPlaintext string) (*domain.User, error) {
	const op errors.Op = "userUsecase.Activate"

	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.Authenticate(ctx, domain.ScopeActivation, tokenPlaintext)
	if err != nil {
		return nil, errors.E(op, err)
	}

	// Update the user's activation status.
	user.Activated = true

	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = uu.Update(ctx, user)
	if err != nil {
		return nil, errors.E(op, err)
	}

	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = uu.tokenUsecase.DeleteAllForUser(ctx, domain.ScopeActivation, user.ID)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return user, nil
}

// Register registers a user and send welcome email with activation token, then insert the user into
// database by calling userRepo.Insert.
// It returns error if exists.
func (uu *userUsecase) Register(ctx context.Context, user *domain.User) error {
	const op errors.Op = "userUsecase.Register"

	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	// Insert the user data into the database.
	err := uu.Insert(ctx, user)
	if err != nil {
		return errors.E(op, err)
	}

	// After the user record has been created in the database, generate a new activation
	// token for the user, and insert it to the database.
	token, err := domain.GenerateToken(user.ID, 3*24*time.Hour, domain.ScopeActivation)
	if err != nil {
		return errors.E(op, err)
	}

	err = uu.tokenUsecase.Insert(ctx, token)
	if err != nil {
		return errors.E(op, err)
	}

	uu.pool.Schedule(func() {
		const op errors.Op = "userAPI.RegisterUser"

		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = uu.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			uu.logger.PrintError(
				errors.E(
					op,
					errors.UserEmail(user.Email),
					errors.ErrInternal,
					errors.Msg("failed to send welcome email"),
					err,
				),
				nil,
			)
			return
		}
	})

	return nil
}
