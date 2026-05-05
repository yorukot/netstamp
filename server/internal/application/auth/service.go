package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/yorukot/netstamp/internal/domain/identity"
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

func (s *Service) Register(ctx context.Context, input RegisterInput) (AuthAccessResult, error) {
	email := normalizeEmail(input.Email)

	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return AuthAccessResult{}, err
	}

	user, err := s.users.CreateUser(ctx, CreateUserInput{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return AuthAccessResult{}, err
	}

	return s.issueAccessResult(ctx, user)
}

func (s *Service) Login(ctx context.Context, input LoginInput) (AuthAccessResult, error) {
	email := normalizeEmail(input.Email)
	user, err := s.users.GetUserByEmail(ctx, email)
	if errors.Is(err, ErrUserNotFound) {
		return AuthAccessResult{}, ErrCredentialsInvalid
	}
	if err != nil {
		return AuthAccessResult{}, err
	}

	if err := s.hasher.Compare(input.Password, user.PasswordHash); err != nil {
		return AuthAccessResult{}, ErrCredentialsInvalid
	}

	return s.issueAccessResult(ctx, user)
}

func (s *Service) issueAccessResult(ctx context.Context, user identity.User) (AuthAccessResult, error) {
	token, err := s.tokens.IssueAccessToken(ctx, AccessTokenInput{
		Subject: user.ID,
		Email:   user.Email,
	})
	if err != nil {
		return AuthAccessResult{}, err
	}

	return AuthAccessResult{
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
