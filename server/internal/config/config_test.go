package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	clearConfigEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Env != "local" {
		t.Fatalf("expected local env, got %q", cfg.Env)
	}
	if cfg.ServiceName != "netstamp-api" {
		t.Fatalf("expected default service name, got %q", cfg.ServiceName)
	}
	if cfg.HTTP.Addr != ":8080" {
		t.Fatalf("expected default HTTP addr, got %q", cfg.HTTP.Addr)
	}
	if cfg.GRPC.Addr != ":9090" {
		t.Fatalf("expected default gRPC addr, got %q", cfg.GRPC.Addr)
	}
	if cfg.HTTP.RequestTimeout != 10*time.Second {
		t.Fatalf("expected default request timeout, got %s", cfg.HTTP.RequestTimeout)
	}
	if cfg.Database.Required {
		t.Fatal("expected database to be optional by default")
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv(keyAppEnv, "production")
	t.Setenv(keyServiceName, "netstamp-worker")
	t.Setenv(keyHTTPAddr, ":8181")
	t.Setenv(keyGRPCAddr, ":9191")
	t.Setenv(keyRequestTimeout, "250ms")
	t.Setenv(keyDatabaseRequired, "true")
	t.Setenv(keyDatabaseURL, "postgres://netstamp:netstamp@localhost:5432/netstamp?sslmode=disable")
	t.Setenv(keyDBMaxConns, "12")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Env != "production" {
		t.Fatalf("expected production env, got %q", cfg.Env)
	}
	if cfg.ServiceName != "netstamp-worker" {
		t.Fatalf("expected service override, got %q", cfg.ServiceName)
	}
	if cfg.HTTP.Addr != ":8181" {
		t.Fatalf("expected HTTP addr override, got %q", cfg.HTTP.Addr)
	}
	if cfg.GRPC.Addr != ":9191" {
		t.Fatalf("expected gRPC addr override, got %q", cfg.GRPC.Addr)
	}
	if cfg.HTTP.RequestTimeout != 250*time.Millisecond {
		t.Fatalf("expected request timeout override, got %s", cfg.HTTP.RequestTimeout)
	}
	if !cfg.Database.Required {
		t.Fatal("expected database to be required")
	}
	if cfg.Database.MaxConns != 12 {
		t.Fatalf("expected DB max conns override, got %d", cfg.Database.MaxConns)
	}
}

func TestLoadFromConfigFile(t *testing.T) {
	clearConfigEnv(t)

	dir := t.TempDir()
	configFile := filepath.Join(dir, "local.env")
	err := os.WriteFile(configFile, []byte(strings.Join([]string{
		"APP_ENV=staging",
		"SERVICE_NAME=netstamp-staging",
		"HTTP_ADDR=:8282",
		"REQUEST_TIMEOUT=2s",
		"",
	}, "\n")), 0o600)
	if err != nil {
		t.Fatalf("write config file: %v", err)
	}
	t.Setenv(keyConfigFile, configFile)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Env != "staging" {
		t.Fatalf("expected staging env, got %q", cfg.Env)
	}
	if cfg.ServiceName != "netstamp-staging" {
		t.Fatalf("expected service from config file, got %q", cfg.ServiceName)
	}
	if cfg.HTTP.Addr != ":8282" {
		t.Fatalf("expected HTTP addr from config file, got %q", cfg.HTTP.Addr)
	}
	if cfg.HTTP.RequestTimeout != 2*time.Second {
		t.Fatalf("expected request timeout from config file, got %s", cfg.HTTP.RequestTimeout)
	}
}

func TestLoadReturnsValidationErrors(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv(keyRequestTimeout, "not-a-duration")
	t.Setenv(keyDBMaxConns, "-1")
	t.Setenv(keyDatabaseRequired, "true")

	_, err := Load()
	if err == nil {
		t.Fatal("expected validation error")
	}

	message := err.Error()
	for _, want := range []string{
		"REQUEST_TIMEOUT must be a duration",
		"DB_MAX_CONNS must not be negative",
		"DATABASE_URL must be set when DATABASE_REQUIRED=true",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("expected error to contain %q, got %q", want, message)
		}
	}
}

func clearConfigEnv(t *testing.T) {
	t.Helper()

	for key := range defaultSettings {
		t.Setenv(key, "")
	}
}
