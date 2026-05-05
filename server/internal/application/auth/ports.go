package auth

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, input CreateUserInput) (User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type TokenIssuer interface {
	IssueAccessToken(ctx context.Context, input AccessTokenInput) (IssuedToken, error)
}
