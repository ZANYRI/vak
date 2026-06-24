package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/billing-service/internal/api"
	"github.com/example/billing-service/internal/auth"
	"github.com/example/billing-service/internal/config"
	"github.com/example/billing-service/internal/database"
	"github.com/example/billing-service/internal/queue"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg, e := config.Load()
	if e != nil {
		logger.Error("configuration error", "error", e)
		os.Exit(1)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	db, e := connectDB(ctx, cfg.DatabaseURL)
	if e != nil {
		logger.Error("database unavailable", "error", e)
		os.Exit(1)
	}
	defer db.Close()
	if e = database.Migrate(ctx, db, logger); e != nil {
		logger.Error("migration failed", "error", e)
		os.Exit(1)
	}
	q, e := connectQueue(ctx, cfg.QueueURL, logger)
	if e != nil {
		logger.Error("queue unavailable", "error", e)
		os.Exit(1)
	}
	defer q.Close()
	app := api.New(db, q, auth.New(cfg.JWTAccessSecret, cfg.JWTRefreshSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL), cfg, logger)
	if e = app.BootstrapAdmin(ctx); e != nil {
		logger.Error("bootstrap admin failed", "error", e)
		os.Exit(1)
	}
	server := &http.Server{Addr: ":" + cfg.HTTPPort, Handler: app.Router(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		logger.Info("api listening", "address", server.Addr)
		if e := server.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
			logger.Error("server failed", "error", e)
			cancel()
		}
	}()
	<-ctx.Done()
	shutdown, c := context.WithTimeout(context.Background(), 20*time.Second)
	defer c()
	if e := server.Shutdown(shutdown); e != nil {
		logger.Error("shutdown failed", "error", e)
	}
}
func connectDB(ctx context.Context, url string) (*pgxpool.Pool, error) {
	deadline := time.Now().Add(45 * time.Second)
	var last error
	for time.Now().Before(deadline) {
		db, e := database.Open(ctx, url)
		if e == nil {
			return db, nil
		}
		last = e
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
		}
	}
	return nil, last
}
func connectQueue(ctx context.Context, url string, log *slog.Logger) (*queue.Client, error) {
	deadline := time.Now().Add(45 * time.Second)
	var last error
	for time.Now().Before(deadline) {
		q, e := queue.Connect(url, log)
		if e == nil {
			return q, nil
		}
		last = e
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
		}
	}
	return nil, last
}
