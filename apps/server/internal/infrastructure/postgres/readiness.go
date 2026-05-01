package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewReadinessCheck(pool *pgxpool.Pool, required bool) func(context.Context) error {
	return func(ctx context.Context) error {
		if pool == nil {
			if required {
				return errors.New("database is required but not configured")
			}
			return nil
		}
		return pool.Ping(ctx)
	}
}
