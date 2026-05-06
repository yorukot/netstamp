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
	events SecurityEventRecorder
}

func NewService(users UserRepository, hasher PasswordHasher, tokens TokenIssuer, events SecurityEventRecorder) *Service {
	return &Service{
		users:  users,
		hasher: hasher,
		tokens: tokens,
		events: events,
	}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (AuthAccessResult, error) {
	email := normalizeEmail(input.Email)

	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		s.events.RecordAuthEvent(ctx, AuthEvent{
			Name:    AuthEventRegisterFailure,
			Action:  AuthActionRegister,
			Outcome: AuthOutcomeFailure,
			Reason:  AuthReasonPasswordHashFailed,
			Email:   email,
			Err:     err,
		})
		return AuthAccessResult{}, err
	}

	user, err := s.users.CreateUser(ctx, CreateUserInput{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		event := AuthEvent{
			Name:    AuthEventRegisterFailure,
			Action:  AuthActionRegister,
			Outcome: AuthOutcomeFailure,
			Reason:  AuthReasonUserCreateFailed,
			Email:   email,
			Err:     err,
		}
		if errors.Is(err, ErrEmailAlreadyExists) {
			event.Reason = AuthReasonEmailAlreadyExists
			event.Err = nil
		}
		s.events.RecordAuthEvent(ctx, event)
		return AuthAccessResult{}, err
	}

	result, err := s.issueAccessResult(ctx, user)
	if err != nil {
		s.events.RecordAuthEvent(ctx, AuthEvent{
			Name:    AuthEventTokenIssueFailure,
			Action:  AuthActionRegister,
			Outcome: AuthOutcomeFailure,
			Reason:  AuthReasonAccessTokenIssueFail,
			UserID:  user.ID,
			Email:   user.Email,
			Err:     err,
		})
		return AuthAccessResult{}, err
	}

	s.events.RecordAuthEvent(ctx, AuthEvent{
		Name:    AuthEventRegisterSuccess,
		Action:  AuthActionRegister,
		Outcome: AuthOutcomeSuccess,
		UserID:  user.ID,
		Email:   user.Email,
	})

	return result, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (AuthAccessResult, error) {
	email := normalizeEmail(input.Email)
	user, err := s.users.GetUserByEmail(ctx, email)
	if errors.Is(err, identity.ErrUserNotFound) {
		s.events.RecordAuthEvent(ctx, AuthEvent{
			Name:    AuthEventLoginFailure,
			Action:  AuthActionLogin,
			Outcome: AuthOutcomeFailure,
			Reason:  AuthReasonCredentialsInvalid,
			Email:   email,
		})
		return AuthAccessResult{}, ErrCredentialsInvalid
	}
	if err != nil {
		s.events.RecordAuthEvent(ctx, AuthEvent{
			Name:    AuthEventLoginFailure,
			Action:  AuthActionLogin,
			Outcome: AuthOutcomeFailure,
			Reason:  AuthReasonUserLookupFailed,
			Email:   email,
			Err:     err,
		})
		return AuthAccessResult{}, err
	}

	if err := s.hasher.Compare(input.Password, user.PasswordHash); err != nil {
		s.events.RecordAuthEvent(ctx, AuthEvent{
			Name:    AuthEventLoginFailure,
			Action:  AuthActionLogin,
			Outcome: AuthOutcomeFailure,
			Reason:  AuthReasonCredentialsInvalid,
			UserID:  user.ID,
			Email:   user.Email,
		})
		return AuthAccessResult{}, ErrCredentialsInvalid
	}

	result, err := s.issueAccessResult(ctx, user)
	if err != nil {
		s.events.RecordAuthEvent(ctx, AuthEvent{
			Name:    AuthEventTokenIssueFailure,
			Action:  AuthActionLogin,
			Outcome: AuthOutcomeFailure,
			Reason:  AuthReasonAccessTokenIssueFail,
			UserID:  user.ID,
			Email:   user.Email,
			Err:     err,
		})
		return AuthAccessResult{}, err
	}

	s.events.RecordAuthEvent(ctx, AuthEvent{
		Name:    AuthEventLoginSuccess,
		Action:  AuthActionLogin,
		Outcome: AuthOutcomeSuccess,
		UserID:  user.ID,
		Email:   user.Email,
	})

	return result, nil
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
