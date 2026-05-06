package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	appauth "github.com/yorukot/netstamp/internal/application/auth"
	"github.com/yorukot/netstamp/internal/domain/identity"
	"github.com/yorukot/netstamp/internal/infrastructure/postgres/sqlc"
)

var postgresTracer = otel.Tracer("github.com/yorukot/netstamp/internal/infrastructure/postgres")

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{queries: sqlc.New(pool)}
}

func (r *UserRepository) CreateUser(ctx context.Context, input appauth.CreateUserInput) (identity.User, error) {
	ctx, span := startUserDBSpan(ctx, "postgres.users.insert", "INSERT", "INSERT users")
	defer span.End()

	row, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
	})
	if err != nil {
		if isUniqueViolation(err, "uq_users_email") {
			return identity.User{}, fmt.Errorf("email already exists: %w", appauth.ErrEmailAlreadyExists)
		}
		recordDBSpanError(span, err)
		return identity.User{}, err
	}

	span.SetAttributes(attribute.Int64("db.response.returned_rows", 1))
	return identity.User{
		ID:        row.ID.String(),
		Email:     row.Email,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (identity.User, error) {
	ctx, span := startUserDBSpan(ctx, "postgres.users.select_by_email", "SELECT", "SELECT users by email")
	defer span.End()

	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.SetAttributes(attribute.Int64("db.response.returned_rows", 0))
			return identity.User{}, identity.ErrUserNotFound
		}

		recordDBSpanError(span, err)
		return identity.User{}, err
	}

	span.SetAttributes(attribute.Int64("db.response.returned_rows", 1))
	return identity.User{
		ID:           row.ID.String(),
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}, nil
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505" && pgErr.ConstraintName == constraint
}

func startUserDBSpan(ctx context.Context, name string, operation string, summary string) (context.Context, trace.Span) {
	return postgresTracer.Start(ctx, name, trace.WithAttributes(
		attribute.String("db.system.name", "postgresql"),
		attribute.String("db.operation.name", operation),
		attribute.String("db.collection.name", "users"),
		attribute.String("db.query.summary", summary),
	))
}

func recordDBSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, "database query failed")

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		span.SetAttributes(
			attribute.String("db.response.status_code", pgErr.Code),
			attribute.String("error.type", pgErr.Code),
		)
	}
}
