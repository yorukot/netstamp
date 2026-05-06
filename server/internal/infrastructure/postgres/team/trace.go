package pgteam

import "go.opentelemetry.io/otel"

var pgteamTracer = otel.Tracer("github.com/yorukot/netstamp/internal/infrastructure/postgres")
