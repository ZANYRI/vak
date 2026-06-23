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
	app := api.New(db, q, auth.New(cfg.JWTAccessSecret, cfg.JWTRefreshSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL), cfg, log)
	if e = q.Consume(ctx, app.ProcessJob); e != nil {
		log.Error("worker stopped", "error", e)
		time.Sleep(time.Second)
		os.Exit(1)
	}
}
