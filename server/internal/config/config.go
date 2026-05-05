package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	keyAppEnv                = "APP_ENV"
	keyServiceName           = "SERVICE_NAME"
	keyAppVersion            = "APP_VERSION"
	keyLogLevel              = "LOG_LEVEL"
	keyShutdownTimeout       = "SHUTDOWN_TIMEOUT"
	keyHTTPAddr              = "HTTP_ADDR"
	keyGRPCAddr              = "GRPC_ADDR"
	keyRequestTimeout        = "REQUEST_TIMEOUT"
	keyHTTPReadHeaderTimeout = "HTTP_READ_HEADER_TIMEOUT"
	keyHTTPReadTimeout       = "HTTP_READ_TIMEOUT"
	keyHTTPWriteTimeout      = "HTTP_WRITE_TIMEOUT"
	keyHTTPIdleTimeout       = "HTTP_IDLE_TIMEOUT"
	keyDatabaseRequired      = "DATABASE_REQUIRED"
	keyDatabaseURL           = "DATABASE_URL"
	keyDBMaxConns            = "DB_MAX_CONNS"
	keyDBMinConns            = "DB_MIN_CONNS"
	keyDBMaxConnLifetime     = "DB_MAX_CONN_LIFETIME"
	keyDBMaxConnIdleTime     = "DB_MAX_CONN_IDLE_TIME"
)

var defaultSettings = map[string]any{
	keyAppEnv:                "local",
	keyServiceName:           "netstamp-api",
	keyAppVersion:            "dev",
	keyLogLevel:              "info",
	keyShutdownTimeout:       10 * time.Second,
	keyHTTPAddr:              ":8080",
	keyGRPCAddr:              ":9090",
	keyRequestTimeout:        10 * time.Second,
	keyHTTPReadHeaderTimeout: 5 * time.Second,
	keyHTTPReadTimeout:       15 * time.Second,
	keyHTTPWriteTimeout:      15 * time.Second,
	keyHTTPIdleTimeout:       60 * time.Second,
	keyDatabaseRequired:      false,
	keyDatabaseURL:           "",
	keyDBMaxConns:            int32(10),
	keyDBMinConns:            int32(0),
	keyDBMaxConnLifetime:     time.Hour,
	keyDBMaxConnIdleTime:     30 * time.Minute,
}

type Config struct {
	Env             string         `mapstructure:"APP_ENV"`
	ServiceName     string         `mapstructure:"SERVICE_NAME"`
	Version         string         `mapstructure:"APP_VERSION"`
	LogLevel        string         `mapstructure:"LOG_LEVEL"`
	ShutdownTimeout time.Duration  `mapstructure:"SHUTDOWN_TIMEOUT"`
	HTTP            HTTPConfig     `mapstructure:",squash"`
	GRPC            GRPCConfig     `mapstructure:",squash"`
	Database        DatabaseConfig `mapstructure:",squash"`
}

type HTTPConfig struct {
	Addr              string        `mapstructure:"HTTP_ADDR"`
	RequestTimeout    time.Duration `mapstructure:"REQUEST_TIMEOUT"`
	ReadHeaderTimeout time.Duration `mapstructure:"HTTP_READ_HEADER_TIMEOUT"`
	ReadTimeout       time.Duration `mapstructure:"HTTP_READ_TIMEOUT"`
	WriteTimeout      time.Duration `mapstructure:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout       time.Duration `mapstructure:"HTTP_IDLE_TIMEOUT"`
}

type GRPCConfig struct {
	Addr string `mapstructure:"GRPC_ADDR"`
}

type DatabaseConfig struct {
	URL             string        `mapstructure:"DATABASE_URL"`
	Required        bool          `mapstructure:"DATABASE_REQUIRED"`
	MaxConns        int32         `mapstructure:"DB_MAX_CONNS"`
	MinConns        int32         `mapstructure:"DB_MIN_CONNS"`
	MaxConnLifetime time.Duration `mapstructure:"DB_MAX_CONN_LIFETIME"`
	MaxConnIdleTime time.Duration `mapstructure:"DB_MAX_CONN_IDLE_TIME"`
}

func Load() (Config, error) {
	settings, err := newSettings()
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	var errs []error
	if err := settings.UnmarshalExact(&cfg); err != nil {
		errs = append(errs, fmt.Errorf("decode config: %w", err))
	}

	errs = append(errs, validate(cfg)...)
	return cfg, errors.Join(errs...)
}

func validate(cfg Config) []error {
	var errs []error
	if cfg.HTTP.Addr == "" {
		errs = append(errs, errors.New("HTTP_ADDR must not be empty"))
	}
	if cfg.GRPC.Addr == "" {
		errs = append(errs, errors.New("GRPC_ADDR must not be empty"))
	}
	if cfg.Database.MaxConns < 0 {
		errs = append(errs, errors.New("DB_MAX_CONNS must not be negative"))
	}
	if cfg.Database.MinConns < 0 {
		errs = append(errs, errors.New("DB_MIN_CONNS must not be negative"))
	}
	if cfg.Database.Required && cfg.Database.URL == "" {
		errs = append(errs, errors.New("DATABASE_URL must be set when DATABASE_REQUIRED=true"))
	}

	return errs
}

func newSettings() (*viper.Viper, error) {
	settings := viper.New()
	settings.SetConfigName(".env")
	settings.SetConfigType("env")
	settings.AddConfigPath(".")
	settings.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	settings.AutomaticEnv()

	var errs []error
	for key, value := range defaultSettings {
		settings.SetDefault(key, value)
		if err := settings.BindEnv(key); err != nil {
			errs = append(errs, fmt.Errorf("bind %s: %w", key, err))
		}
	}

	if err := loadDotEnv(settings); err != nil {
		errs = append(errs, err)
	}

	return settings, errors.Join(errs...)
}

func loadDotEnv(settings *viper.Viper) error {
	if err := settings.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return nil
		}
		return fmt.Errorf("read .env: %w", err)
	}

	return nil
}
