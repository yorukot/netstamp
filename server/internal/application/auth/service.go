package auth

import (
	"context"
	"strings"
)

type Service struct {
	users  UserRepository
	hasher PasswordHasher
	tokens TokenIssuer
}

func NewService(users UserRepository, hasher PasswordHasher, tokens TokenIssuer) *Service {
	return &Service{
		users:  users,
		hasher: hasher,
		tokens: tokens,
	}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (RegisterResult, error) {
	email := normalizeEmail(input.Email)
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return RegisterResult{}, err
	}

	user, err := s.users.CreateUser(ctx, CreateUserInput{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return RegisterResult{}, err
	}

	token, err := s.tokens.IssueAccessToken(ctx, AccessTokenInput{
		Subject: user.ID,
		Email:   user.Email,
	})
	if err != nil {
		return RegisterResult{}, err
	}

	return RegisterResult{
		UserID:      user.ID,
		Email:       user.Email,
		AccessToken: token.Value,
		TokenType:   token.TokenType,
		ExpiresIn:   token.ExpiresIn,
	}, nil
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
