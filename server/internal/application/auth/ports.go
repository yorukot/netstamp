package auth

import (
	"context"

	"github.com/yorukot/netstamp/internal/domain/identity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, input CreateUserInput) (identity.User, error)
	GetUserByEmail(ctx context.Context, email string) (identity.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password string, passwordHash string) error
}

type TokenIssuer interface {
	IssueAccessToken(ctx context.Context, input AccessTokenInput) (IssuedToken, error)
}
