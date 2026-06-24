package main

import (
	"context"
	"log/slog"
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
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, e := config.Load()
	if e != nil {
		log.Error("configuration error", "error", e)
		os.Exit(1)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	db, e := database.Open(ctx, cfg.DatabaseURL)
	if e != nil {
		log.Error("database unavailable", "error", e)
		os.Exit(1)
	}
	defer db.Close()
	q, e := queue.Connect(cfg.QueueURL, log)
	if e != nil {
		log.Error("queue unavailable", "error", e)
		os.Exit(1)
	}
	defer q.Close()
	if e := waitForSchema(ctx, db); e != nil {
		log.Error("database schema unavailable", "error", e)
		return
	}
	app := api.New(db, q, auth.New(cfg.JWTAccessSecret, cfg.JWTRefreshSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL), cfg, log)
	run := func() {
		jobctx, c := context.WithTimeout(ctx, 20*time.Second)
		defer c()
		if e := app.PublishScheduled(jobctx); e != nil {
			log.Error("publish schedule", "error", e)
		}
	}
	run()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}

func waitForSchema(ctx context.Context, db *pgxpool.Pool) error {
	for {
		var ready bool
		err := db.QueryRow(ctx, `SELECT to_regclass('public.jobs') IS NOT NULL`).Scan(&ready)
		if err == nil && ready {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}
