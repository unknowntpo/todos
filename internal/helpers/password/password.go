package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Password is a struct containing the plaintext and hashed
// versions of the password for a user.
type Password struct {
	plaintext string
	hash      []byte
}

// Set calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// Matches checks whether the provided plaintext password matches the
// hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// GetPlainText returns the plain text password, if plain text password is not set, we return the empty string.
func (p *Password) GetPlainText() string {
	return p.plaintext
}

// GetHash returns the hash version of password, if the hash is not set, return nil.
func (p *Password) GetHash() []byte {
	return p.hash
}
