package auth

import "time"

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterResult struct {
	UserID      string
	Email       string
	AccessToken string
	TokenType   string
	ExpiresIn   int
}

type CreateUserInput struct {
	Email        string
	PasswordHash string
}

type User struct {
	ID        string
	Email     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccessTokenInput struct {
	Subject string
	Email   string
}

type IssuedToken struct {
	Value     string
	TokenType string
	ExpiresIn int
}
