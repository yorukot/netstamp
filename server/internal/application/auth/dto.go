package auth

type RegisterInput struct {
	Email    string
	Password string
}

type CreateUserInput struct {
	Email        string
	PasswordHash string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthAccessResult struct {
	UserID      string
	Email       string
	AccessToken string
	TokenType   string
	ExpiresIn   int
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

type AccessTokenClaims struct {
	Subject string
	Email   string
}
