package domain

import (
	"errors"
)

var (
	ErrRecordNotFound     = errors.New("record not found")    // Record not found when we request some resource in database.
	ErrDuplicateEmail     = errors.New("duplicate email")     // Duplicate Email error.
	ErrEditConflict       = errors.New("edit conflict")       // Edit conflict while manipulating database.
	ErrInvalidCredentials = errors.New("invalid credentials") // Edit conflict while manipulating database.
	ErrFailedValidation   = errors.New("failed validation")   //  Failed validation error.
)
