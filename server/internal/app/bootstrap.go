package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	apphello "github.com/yorukot/netstamp/internal/application/hello"
	"github.com/yorukot/netstamp/internal/config"
	"github.com/yorukot/netstamp/internal/infrastructure/postgres"
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
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log, _, err := logger.New(logger.Config{
		Env:     cfg.Env,
		Service: cfg.ServiceName,
		Version: cfg.Version,
		Level:   cfg.LogLevel,
	})
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	dbPool, err := openDatabase(ctx, cfg.Database)
	if err != nil {
		_ = log.Sync()
		return nil, err
	}

	helloSvc := apphello.NewService(cfg.ServiceName)
	readiness := postgres.NewReadinessCheck(dbPool, cfg.Database.Required)

	httpHandler := httpserver.NewRouter(httpserver.Dependencies{
		Log:            log,
		HelloService:   helloSvc,
		ReadinessCheck: readiness,
		RequestTimeout: cfg.HTTP.RequestTimeout,
	})

	return &Application{
		Config: cfg,
		Log:    log,
		HTTPServer: httpserver.NewServer(httpserver.Config{
			Addr:              cfg.HTTP.Addr,
			ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
			ReadTimeout:       cfg.HTTP.ReadTimeout,
			WriteTimeout:      cfg.HTTP.WriteTimeout,
			IdleTimeout:       cfg.HTTP.IdleTimeout,
		}, httpHandler),
		GRPCServer: grpcserver.New(grpcserver.Dependencies{
			Log:         log,
			ServiceName: cfg.ServiceName,
		}),
		DBPool: dbPool,
	}, nil
}

func openDatabase(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	if cfg.URL == "" {
		if cfg.Required {
			return nil, fmt.Errorf("DATABASE_URL is required when DATABASE_REQUIRED=true")
		}
		return nil, nil
	}

	pool, err := postgres.NewPool(ctx, postgres.PoolConfig{
		URL:             cfg.URL,
		MaxConns:        cfg.MaxConns,
		MinConns:        cfg.MinConns,
		MaxConnLifetime: cfg.MaxConnLifetime,
		MaxConnIdleTime: cfg.MaxConnIdleTime,
	})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	return pool, nil
}
