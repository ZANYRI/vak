// Command scheduler periodically enqueues recurring billing jobs.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/billing-service/internal/app"
	"github.com/example/billing-service/internal/config"
	"github.com/example/billing-service/internal/observability"
	"github.com/example/billing-service/internal/scheduler"
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

	sched := scheduler.New(log, container.Subs, container.Invoices, container.Queue, cfg.SchedulerInterval)
	if err := sched.Run(ctx); err != nil {
		log.Error("scheduler error", "error", err.Error())
	}
	log.Info("scheduler stopped")
}

func cfgEnv(cfg *config.Config) string {
	if cfg == nil {
		return "production"
	}
	return cfg.AppEnv
}
