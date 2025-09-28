package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"playgoround/internal/app"
	"playgoround/internal/config"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	application, err := app.New(cfg)
	if err != nil {
		slog.Error("failed to init application", slog.Any("err", err))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err = application.Start(ctx)
	if err != nil {
		slog.Error("failed to start application", slog.Any("err", err))
		os.Exit(1)
	}

	<-ctx.Done()

	slog.Info("GracefulShutdown started")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := application.Stop(shutdownCtx); err != nil {
		slog.Error("graceful shutdown error", slog.Any("err", err))
	}
	slog.Info("GracefulShutdown completed")
}
