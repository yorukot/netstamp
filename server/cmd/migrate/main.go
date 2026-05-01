package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/yorukot/netstamp/internal/config"
)

func main() {
	databaseURL := flag.String("database-url", "", "PostgreSQL connection URL")
	dir := flag.String("dir", "db/migrations", "migration directory")
	command := flag.String("command", "status", "migration command: status, up, or down")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	resolvedDatabaseURL := *databaseURL
	if resolvedDatabaseURL == "" {
		resolvedDatabaseURL = cfg.Database.URL
	}

	if resolvedDatabaseURL == "" {
		_, _ = fmt.Fprintln(os.Stderr, "DATABASE_URL or -database-url is required")
		os.Exit(2)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "set migration dialect: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("pgx", resolvedDatabaseURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ping database: %v\n", err)
		os.Exit(1)
	}

	switch *command {
	case "status":
		err = goose.Status(db, *dir)
	case "up":
		err = goose.Up(db, *dir)
	case "down":
		err = goose.Down(db, *dir)
	default:
		err = fmt.Errorf("unsupported migration command %q", *command)
	}
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
		os.Exit(1)
	}
}
