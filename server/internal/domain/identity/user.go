package identity

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID           string
	Email        string
	DisplayName  *string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
