package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	appauth "github.com/yorukot/netstamp/internal/application/auth"
	"github.com/yorukot/netstamp/internal/config"
	"github.com/yorukot/netstamp/internal/infrastructure/postgres"
	"github.com/yorukot/netstamp/internal/infrastructure/security"
	"github.com/yorukot/netstamp/internal/logger"
	grpcserver "github.com/yorukot/netstamp/internal/transport/grpc"
	httpserver "github.com/yorukot/netstamp/internal/transport/http"
)

type Application struct {
	Config     config.Config
	Log        *zap.Logger
	HTTPServer *http.Server
	GRPCServer *grpc.Server
	DBPool     *pgxpool.Pool
}

func New(ctx context.Context) (*Application, error) {
	// Load configuration from environment variables and .env file
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// Creating logger before database connection to ensure we can log any errors that occur during startup
	log, _, err := logger.New(logger.Config{
		Env:     cfg.Env,
		Service: cfg.ServiceName,
		Version: cfg.Version,
		Level:   cfg.LogLevel,
	})
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	// Open database connection pool.
	dbPool, err := postgres.NewPool(ctx, postgres.PoolConfig{
		ConnectionString: cfg.Database.ConnectionString(),
		MaxConns:         cfg.Database.MaxConns,
		MinConns:         cfg.Database.MinConns,
		MaxConnLifetime:  cfg.Database.MaxConnLifetime,
		MaxConnIdleTime:  cfg.Database.MaxConnIdleTime,
	})
	if err != nil {
		return nil, err
	}

	// Initialize application services and handlers
	userRepo := postgres.NewUserRepository(dbPool)
	passwordHasher := security.NewArgon2idPasswordHasher(security.Argon2idConfig{
		MemoryKiB:   cfg.Auth.Argon2idMemoryKiB,
		Iterations:  cfg.Auth.Argon2idIterations,
		Parallelism: cfg.Auth.Argon2idParallelism,
	})
	tokenIssuer := security.NewJWTIssuer(cfg.Auth.JWTSecret, cfg.Auth.AccessTokenTTL)
	authEvents := logger.NewAuthEventRecorder(log, cfg.LogPseudonymKey)

	authSvc := appauth.NewService(userRepo, passwordHasher, tokenIssuer, authEvents)
	readiness := postgres.NewReadinessCheck(dbPool)

	httpHandler := httpserver.NewRouter(httpserver.Dependencies{
		Log:            log,
		APIVersion:     cfg.Version,
		AuthService:    authSvc,
		ReadinessCheck: readiness,
		RequestTimeout: cfg.HTTP.RequestTimeout,
	})

	// GRPC server setup will be done in the Run method, as it requires the application context and dependencies to be fully initialized.
	return &Application{
		Config:     cfg,
		Log:        log,
		HTTPServer: httpserver.NewServer(cfg.HTTP, httpHandler),
		GRPCServer: grpcserver.New(grpcserver.Dependencies{
			Log:         log,
			ServiceName: cfg.ServiceName,
		}),
		DBPool: dbPool,
	}, nil
}
