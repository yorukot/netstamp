package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	keyConfigFile            = "CONFIG_FILE"
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

var defaultSettings = map[string]string{
	keyConfigFile:            "",
	keyAppEnv:                "local",
	keyServiceName:           "netstamp-api",
	keyAppVersion:            "dev",
	keyLogLevel:              "info",
	keyShutdownTimeout:       "10s",
	keyHTTPAddr:              ":8080",
	keyGRPCAddr:              ":9090",
	keyRequestTimeout:        "10s",
	keyHTTPReadHeaderTimeout: "5s",
	keyHTTPReadTimeout:       "15s",
	keyHTTPWriteTimeout:      "15s",
	keyHTTPIdleTimeout:       "60s",
	keyDatabaseRequired:      "false",
	keyDatabaseURL:           "",
	keyDBMaxConns:            "10",
	keyDBMinConns:            "0",
	keyDBMaxConnLifetime:     "1h",
	keyDBMaxConnIdleTime:     "30m",
}

type Config struct {
	Env             string
	ServiceName     string
	Version         string
	LogLevel        string
	ShutdownTimeout time.Duration
	HTTP            HTTPConfig
	GRPC            GRPCConfig
	Database        DatabaseConfig
}

type HTTPConfig struct {
	Addr              string
	RequestTimeout    time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

type GRPCConfig struct {
	Addr string
}

type DatabaseConfig struct {
	URL             string
	Required        bool
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func Load() (Config, error) {
	settings, err := newSettings()
	if err != nil {
		return Config{}, err
	}

	var errs []error

	shutdownTimeout, err := durationSetting(settings, keyShutdownTimeout)
	errs = appendIfErr(errs, err)

	requestTimeout, err := durationSetting(settings, keyRequestTimeout)
	errs = appendIfErr(errs, err)

	readHeaderTimeout, err := durationSetting(settings, keyHTTPReadHeaderTimeout)
	errs = appendIfErr(errs, err)

	readTimeout, err := durationSetting(settings, keyHTTPReadTimeout)
	errs = appendIfErr(errs, err)

	writeTimeout, err := durationSetting(settings, keyHTTPWriteTimeout)
	errs = appendIfErr(errs, err)

	idleTimeout, err := durationSetting(settings, keyHTTPIdleTimeout)
	errs = appendIfErr(errs, err)

	dbRequired, err := boolSetting(settings, keyDatabaseRequired)
	errs = appendIfErr(errs, err)

	maxConns, err := int32Setting(settings, keyDBMaxConns)
	errs = appendIfErr(errs, err)

	minConns, err := int32Setting(settings, keyDBMinConns)
	errs = appendIfErr(errs, err)

	maxConnLifetime, err := durationSetting(settings, keyDBMaxConnLifetime)
	errs = appendIfErr(errs, err)

	maxConnIdleTime, err := durationSetting(settings, keyDBMaxConnIdleTime)
	errs = appendIfErr(errs, err)

	cfg := Config{
		Env:             stringSetting(settings, keyAppEnv),
		ServiceName:     stringSetting(settings, keyServiceName),
		Version:         stringSetting(settings, keyAppVersion),
		LogLevel:        stringSetting(settings, keyLogLevel),
		ShutdownTimeout: shutdownTimeout,
		HTTP: HTTPConfig{
			Addr:              stringSetting(settings, keyHTTPAddr),
			RequestTimeout:    requestTimeout,
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
		},
		GRPC: GRPCConfig{
			Addr: stringSetting(settings, keyGRPCAddr),
		},
		Database: DatabaseConfig{
			URL:             stringSetting(settings, keyDatabaseURL),
			Required:        dbRequired,
			MaxConns:        maxConns,
			MinConns:        minConns,
			MaxConnLifetime: maxConnLifetime,
			MaxConnIdleTime: maxConnIdleTime,
		},
	}

	if cfg.HTTP.Addr == "" {
		errs = append(errs, errors.New("HTTP_ADDR must not be empty"))
	}
	if cfg.GRPC.Addr == "" {
		errs = append(errs, errors.New("GRPC_ADDR must not be empty"))
	}
	if cfg.Database.Required && cfg.Database.URL == "" {
		errs = append(errs, errors.New("DATABASE_URL must be set when DATABASE_REQUIRED=true"))
	}

	return cfg, errors.Join(errs...)
}

func newSettings() (*viper.Viper, error) {
	settings := viper.New()
	settings.SetConfigType("env")
	settings.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	settings.AutomaticEnv()

	var errs []error
	for key, value := range defaultSettings {
		settings.SetDefault(key, value)
		if err := settings.BindEnv(key); err != nil {
			errs = append(errs, fmt.Errorf("bind %s: %w", key, err))
		}
	}

	if err := readConfigFile(settings); err != nil {
		errs = append(errs, err)
	}

	return settings, errors.Join(errs...)
}

func readConfigFile(settings *viper.Viper) error {
	configFile := stringSetting(settings, keyConfigFile)
	if configFile == "" {
		return nil
	}

	settings.SetConfigFile(configFile)
	if err := settings.ReadInConfig(); err != nil {
		return fmt.Errorf("read config file %q: %w", configFile, err)
	}

	return nil
}

func stringSetting(settings *viper.Viper, key string) string {
	return strings.TrimSpace(settings.GetString(key))
}

func durationSetting(settings *viper.Viper, key string) (time.Duration, error) {
	value := stringSetting(settings, key)
	if value == "" {
		return 0, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a duration: %w", key, err)
	}
	return parsed, nil
}

func int32Setting(settings *viper.Viper, key string) (int32, error) {
	value := stringSetting(settings, key)
	if value == "" {
		return 0, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%s must be an int32: %w", key, err)
	}
	if parsed < 0 {
		return 0, fmt.Errorf("%s must not be negative", key)
	}
	return int32(parsed), nil
}

func boolSetting(settings *viper.Viper, key string) (bool, error) {
	value := stringSetting(settings, key)
	if value == "" {
		return false, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", key, err)
	}
	return parsed, nil
}

func appendIfErr(errs []error, err error) []error {
	if err != nil {
		return append(errs, err)
	}
	return errs
}
