# NetStamp

NetStamp is a Go backend scaffold using Chi for HTTP, gRPC for internal APIs, Zap for structured logging, PostgreSQL via pgxpool, sqlc for type-safe SQL generation, and Goose for versioned migrations.

This repository currently contains only the initial architecture and a minimal Hello World proof of setup.

## Requirements

- Go 1.25+
- PostgreSQL, only when database-backed features or migrations are being exercised

## Run Locally

```sh
go mod download
go run ./cmd/api
```

HTTP verification:

```sh
curl http://localhost:8080/v1/hello
curl http://localhost:8080/readyz
```

Expected Hello response:

```json
{"message":"Hello from NetStamp","service":"netstamp-api"}
```

gRPC verification with grpcurl:

```sh
grpcurl -plaintext localhost:9090 grpc.health.v1.Health/Check
```

## Configuration

Configuration is managed with Viper in `internal/config`.

Precedence is:

```text
environment variables > CONFIG_FILE envfile > defaults
```

Use `CONFIG_FILE` to load a local envfile. `scripts/dev.sh` sets `CONFIG_FILE=.env` automatically when `.env` exists.

```sh
cp configs/local.env.example .env
./scripts/dev.sh
```

The API starts without PostgreSQL by default. To require PostgreSQL, run the local database and set `DATABASE_REQUIRED=true` plus `DATABASE_URL`.

```sh
docker compose -f deployments/docker/docker-compose.yml up -d
```

## Test

```sh
go test ./...
```

## Migrations

Migration files belong in `db/migrations`.

```sh
DATABASE_URL=postgres://netstamp:netstamp@localhost:5432/netstamp?sslmode=disable \
  go run ./cmd/migrate -command status
```

## Project Structure

```text
cmd/
  api/                 application entrypoint
  migrate/             Goose migration command
internal/
  app/                 dependency wiring and lifecycle
  config/              environment parsing and validation
  logger/              Zap setup and request-scoped logger helpers
  domain/              pure domain concepts
  application/         use cases and application DTOs
  infrastructure/      PostgreSQL pool, transactions, and future repositories
  transport/
    http/              Chi router, handlers, middleware, response encoding
    grpc/              gRPC server and interceptors
api/proto/             protobuf source files
db/migrations/         versioned database migrations
db/query/              sqlc query files
configs/               local environment examples
deployments/           local deployment support
scripts/               development commands
tools/                 tool dependency pinning
```

## Architectural Decisions

- HTTP and gRPC are adapters. They decode requests, call application services, and encode responses.
- Domain code does not import Chi, gRPC, Zap, pgx, sqlc, or HTTP types.
- Application services own use-case orchestration and expose transport-neutral DTOs.
- Infrastructure owns PostgreSQL details and will hold sqlc-backed repositories.
- Zap is configured once at startup and request fields are attached at HTTP middleware and gRPC interceptor boundaries.
- PostgreSQL is optional for the current Hello World proof, but the pool and readiness wiring are already in place for future features.

## Where Future Features Go

- Add domain entities and ports under `internal/domain/<feature>`.
- Add use cases under `internal/application/<feature>`.
- Add REST handlers under `internal/transport/http/<feature>`.
- Add protobuf APIs under `api/proto/netstamp/<feature>/v1`.
- Add generated protobuf code under `api/gen/go` after running Buf.
- Add sqlc queries under `db/query` and generated code under `internal/infrastructure/postgres/sqlc`.
- Add repository implementations under `internal/infrastructure/postgres`.
- Add schema changes as Goose migrations under `db/migrations`.
