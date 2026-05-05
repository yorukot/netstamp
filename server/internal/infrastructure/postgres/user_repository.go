package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	appauth "github.com/yorukot/netstamp/internal/application/auth"
	"github.com/yorukot/netstamp/internal/infrastructure/postgres/sqlc"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{queries: sqlc.New(pool)}
}

func (r *UserRepository) CreateUser(ctx context.Context, input appauth.CreateUserInput) (appauth.User, error) {
	row, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
	})
	if err != nil {
		if isUniqueViolation(err, "uq_users_email") {
			return appauth.User{}, appauth.ErrEmailAlreadyExists
		}
		return appauth.User{}, err
	}

	return appauth.User{
		ID:        row.ID.String(),
		Email:     row.Email,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505" && pgErr.ConstraintName == constraint
}
