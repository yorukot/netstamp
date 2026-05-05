package auth

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrCredentialsInvalid = errors.New("credentials invalid")
	ErrUserNotFound       = errors.New("user not found")
)
