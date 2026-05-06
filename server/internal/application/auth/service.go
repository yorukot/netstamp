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

// Register is the service entry for the register action
func (s *Service) Register(ctx context.Context, input RegisterInput) (AuthAccessResult, error) {
	email := normalizeEmail(input.Email)
	ctx, flow := s.startAuthFlow(ctx, "auth.register", AuthActionRegister, email)
	defer flow.End()

	passwordHash, err := s.hashPassword(ctx, input.Password)
	if err != nil {
		return AuthAccessResult{}, flow.TechnicalFailure(AuthEventRegisterFailure, AuthReasonPasswordHashFailed, err)
	}

	user, err := s.createUser(ctx, CreateUserInput{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if errors.Is(err, ErrEmailAlreadyExists) {
		return AuthAccessResult{}, flow.BusinessFailure(AuthEventRegisterFailure, AuthReasonEmailAlreadyExists, err)
	}
	if err != nil {
		return AuthAccessResult{}, flow.TechnicalFailure(AuthEventRegisterFailure, AuthReasonUserCreateFailed, err)
	}
	flow.SetUser(user)

	result, err := s.issueAccessResult(ctx, user)
	if err != nil {
		return AuthAccessResult{}, flow.TechnicalFailure(AuthEventTokenIssueFailure, AuthReasonAccessTokenIssueFail, err)
	}

	flow.Success(AuthEventRegisterSuccess)

	return result, nil
}

// Login is the enrty for the login action
func (s *Service) Login(ctx context.Context, input LoginInput) (AuthAccessResult, error) {
	email := normalizeEmail(input.Email)
	ctx, flow := s.startAuthFlow(ctx, "auth.login", AuthActionLogin, email)
	defer flow.End()

	user, err := s.getUserByEmail(ctx, email)
	if errors.Is(err, identity.ErrUserNotFound) {
		return AuthAccessResult{}, flow.BusinessFailure(AuthEventLoginFailure, AuthReasonCredentialsInvalid, ErrCredentialsInvalid)
	}
	if err != nil {
		return AuthAccessResult{}, flow.TechnicalFailure(AuthEventLoginFailure, AuthReasonUserLookupFailed, err)
	}
	flow.SetUser(user)

	if err := s.comparePassword(ctx, input.Password, user.PasswordHash); err != nil {
		return AuthAccessResult{}, flow.BusinessFailure(AuthEventLoginFailure, AuthReasonCredentialsInvalid, ErrCredentialsInvalid)
	}

	if !user.IsActive {
		return AuthAccessResult{}, flow.BusinessFailure(AuthEventLoginFailure, AuthReasonUserInactive, ErrUserInactive)
	}

	result, err := s.issueAccessResult(ctx, user)
	if err != nil {
		return AuthAccessResult{}, flow.TechnicalFailure(AuthEventTokenIssueFailure, AuthReasonAccessTokenIssueFail, err)
	}

	flow.Success(AuthEventLoginSuccess)

	return result, nil
}

func (s *Service) issueAccessResult(ctx context.Context, user identity.User) (AuthAccessResult, error) {
	ctx, span := authTracer.Start(ctx, "auth.issue_access_token")
	defer span.End()

	token, err := s.tokens.IssueAccessToken(ctx, AccessTokenInput{
		Subject: user.ID,
		Email:   user.Email,
	})
	if err != nil {
		recordSpanError(span, err, AuthReasonAccessTokenIssueFail)
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

func (s *Service) hashPassword(ctx context.Context, password string) (string, error) {
	_, span := authTracer.Start(ctx, "auth.password_hash")
	defer span.End()

	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		recordSpanError(span, err, AuthReasonPasswordHashFailed)
		return "", err
	}

	return passwordHash, nil
}

func (s *Service) comparePassword(ctx context.Context, password string, passwordHash string) error {
	_, span := authTracer.Start(ctx, "auth.password_compare")
	defer span.End()

	return s.hasher.Compare(password, passwordHash)
}

func (s *Service) createUser(ctx context.Context, input CreateUserInput) (identity.User, error) {
	ctx, span := authTracer.Start(ctx, "auth.create_user")
	defer span.End()

	user, err := s.users.CreateUser(ctx, input)
	if err != nil {
		if !errors.Is(err, ErrEmailAlreadyExists) {
			recordSpanError(span, err, AuthReasonUserCreateFailed)
		}
		return identity.User{}, err
	}

	return user, nil
}

func (s *Service) getUserByEmail(ctx context.Context, email string) (identity.User, error) {
	ctx, span := authTracer.Start(ctx, "auth.get_user_by_email")
	defer span.End()

	user, err := s.users.GetUserByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, identity.ErrUserNotFound) {
			recordSpanError(span, err, AuthReasonUserLookupFailed)
		}
		return identity.User{}, err
	}

	return user, nil
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
