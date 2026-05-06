package auth

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrCredentialsInvalid = errors.New("credentials invalid")
	ErrUserInactive       = errors.New("user inactive")
	ErrAccessTokenInvalid = errors.New("access token invalid")
)
