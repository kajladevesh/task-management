package errors

import (
	"errors"
)

var (
	// ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmailOrUsernameTaken = errors.New("email or username is already taken")

	ErrUserNotFound       = errors.New("user not found")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrMissingUsername    = errors.New("missing username")
	ErrMissingPassword    = errors.New("missing password")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
