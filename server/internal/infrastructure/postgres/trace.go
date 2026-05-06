package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func StartUserDBSpan(ctx context.Context, tracer trace.Tracer, name string, operation string, summary string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, trace.WithAttributes(
		attribute.String("db.system.name", "postgresql"),
		attribute.String("db.operation.name", operation),
		attribute.String("db.collection.name", "users"),
		attribute.String("db.query.summary", summary),
	))
}

func RecordDBSpanError(span trace.Span, err error) {
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
