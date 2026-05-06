package auth

import "errors"

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrDisplayNameRequired = errors.New("display name required")
	ErrDisplayNameTooLong  = errors.New("display name too long")
	ErrCredentialsInvalid  = errors.New("credentials invalid")
	ErrUserInactive        = errors.New("user inactive")
	ErrAccessTokenInvalid  = errors.New("access token invalid")
)
