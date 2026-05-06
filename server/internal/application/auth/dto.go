package auth

type RegisterInput struct {
	Email       string
	DisplayName string
	Password    string
}

type CreateUserInput struct {
	Email        string
	DisplayName  string
	PasswordHash string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthAccessResult struct {
	UserID      string
	Email       string
	DisplayName *string
	AccessToken string
	TokenType   string
	ExpiresIn   int
}

type AccessTokenInput struct {
	Subject     string
	Email       string
	DisplayName *string
}

type IssuedToken struct {
	Value     string
	TokenType string
	ExpiresIn int
}

type AccessTokenClaims struct {
	Subject     string
	Email       string
	DisplayName *string
}
