package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/yorukot/netstamp/internal/app"
)

func main() {
	// Graceful shutdown on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// New application with the context
	application, err := app.New(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "startup failed: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = application.Log.Sync()
	}()
	
	err = application.Run(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		application.Log.Error("startup failed", zap.Error(err))
		os.Exit(1)
	}
}
