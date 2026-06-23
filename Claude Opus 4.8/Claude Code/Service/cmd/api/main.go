// Command api runs the HTTP server for the billing service.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/billing-service/internal/app"
	"github.com/example/billing-service/internal/config"
	"github.com/example/billing-service/internal/observability"
	"github.com/example/billing-service/internal/server"
)

func main() {
	cfg, err := config.Load()
	log := observability.NewLogger(cfgEnv(cfg))
	if err != nil {
		log.Error("config error", "error", err.Error())
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container, err := app.Build(ctx, cfg, log)
	if err != nil {
		log.Error("failed to build application", "error", err.Error())
		os.Exit(1)
	}
	defer container.Close()

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           server.NewRouter(container),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Info("api listening", "addr", srv.Addr, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server error", "error", err.Error())
			stop()
		}
	}()

	<-ctx.Done()
	log.Info("shutdown signal received, draining connections")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "error", err.Error())
	}
	log.Info("api stopped")
}

func cfgEnv(cfg *config.Config) string {
	if cfg == nil {
		return "production"
	}
	return cfg.AppEnv
}
