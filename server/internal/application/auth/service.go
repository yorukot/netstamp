package auth

import (
	"context"
	"errors"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/yorukot/netstamp/internal/domain/identity"
)

var authTracer = otel.Tracer("github.com/yorukot/netstamp/internal/application/auth")

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
	ctx, span := authTracer.Start(ctx, "auth.register", trace.WithAttributes(
		attribute.String("auth.action", "register"),
	))
	defer span.End()

	email := normalizeEmail(input.Email)

	passwordHash, err := s.hashPassword(ctx, input.Password)
	if err != nil {
		recordSpanError(span, err, string(AuthReasonPasswordHashFailed))
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

	user, err := s.createUser(ctx, CreateUserInput{
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
		span.SetAttributes(attribute.String("auth.failure.reason", string(event.Reason)))
		if event.Err != nil {
			recordSpanError(span, event.Err, string(event.Reason))
		}
		return AuthAccessResult{}, err
	}
	span.SetAttributes(attribute.String("user.id", user.ID))

	result, err := s.issueAccessResult(ctx, user)
	if err != nil {
		recordSpanError(span, err, string(AuthReasonAccessTokenIssueFail))
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

	span.SetAttributes(attribute.String("auth.outcome", "success"))
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
	ctx, span := authTracer.Start(ctx, "auth.login", trace.WithAttributes(
		attribute.String("auth.action", "login"),
	))
	defer span.End()

	email := normalizeEmail(input.Email)
	user, err := s.getUserByEmail(ctx, email)
	if errors.Is(err, identity.ErrUserNotFound) {
		span.SetAttributes(
			attribute.String("auth.outcome", "failure"),
			attribute.String("auth.failure.reason", string(AuthReasonCredentialsInvalid)),
		)
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
		recordSpanError(span, err, string(AuthReasonUserLookupFailed))
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
	span.SetAttributes(attribute.String("user.id", user.ID))

	if err := s.comparePassword(ctx, input.Password, user.PasswordHash); err != nil {
		span.SetAttributes(
			attribute.String("auth.outcome", "failure"),
			attribute.String("auth.failure.reason", string(AuthReasonCredentialsInvalid)),
		)
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
		recordSpanError(span, err, string(AuthReasonAccessTokenIssueFail))
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

	span.SetAttributes(attribute.String("auth.outcome", "success"))
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
	ctx, span := authTracer.Start(ctx, "auth.issue_access_token")
	defer span.End()

	token, err := s.tokens.IssueAccessToken(ctx, AccessTokenInput{
		Subject: user.ID,
		Email:   user.Email,
	})
	if err != nil {
		recordSpanError(span, err, string(AuthReasonAccessTokenIssueFail))
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
		recordSpanError(span, err, string(AuthReasonPasswordHashFailed))
		return "", err
	}

	return passwordHash, nil
}

func (s *Service) comparePassword(ctx context.Context, password string, passwordHash string) error {
	_, span := authTracer.Start(ctx, "auth.password_compare")
	defer span.End()

	err := s.hasher.Compare(password, passwordHash)
	if err != nil {
		span.SetAttributes(
			attribute.String("auth.outcome", "failure"),
			attribute.String("auth.failure.reason", string(AuthReasonCredentialsInvalid)),
		)
		return err
	}

	return nil
}

func (s *Service) createUser(ctx context.Context, input CreateUserInput) (identity.User, error) {
	ctx, span := authTracer.Start(ctx, "auth.create_user")
	defer span.End()

	user, err := s.users.CreateUser(ctx, input)
	if err != nil {
		reason := AuthReasonUserCreateFailed
		if errors.Is(err, ErrEmailAlreadyExists) {
			reason = AuthReasonEmailAlreadyExists
		}
		span.SetAttributes(attribute.String("auth.failure.reason", string(reason)))
		if reason != AuthReasonEmailAlreadyExists {
			recordSpanError(span, err, string(reason))
		}
		return identity.User{}, err
	}

	span.SetAttributes(attribute.String("user.id", user.ID))
	return user, nil
}

func (s *Service) getUserByEmail(ctx context.Context, email string) (identity.User, error) {
	ctx, span := authTracer.Start(ctx, "auth.get_user_by_email")
	defer span.End()

	user, err := s.users.GetUserByEmail(ctx, email)
	if err != nil {
		reason := AuthReasonUserLookupFailed
		if errors.Is(err, identity.ErrUserNotFound) {
			reason = AuthReasonCredentialsInvalid
		}
		span.SetAttributes(attribute.String("auth.failure.reason", string(reason)))
		if reason != AuthReasonCredentialsInvalid {
			recordSpanError(span, err, string(reason))
		}
		return identity.User{}, err
	}

	span.SetAttributes(attribute.String("user.id", user.ID))
	return user, nil
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func recordSpanError(span trace.Span, err error, description string) {
	span.RecordError(err)
	span.SetStatus(codes.Error, description)
	span.SetAttributes(
		attribute.String("error.reason", description),
	)
}
