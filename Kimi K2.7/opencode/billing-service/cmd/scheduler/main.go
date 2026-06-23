package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"billing-service/internal/config"
	"billing-service/internal/database"
	"billing-service/internal/observability"
	"billing-service/internal/plans"
	"billing-service/internal/queue"
	"billing-service/internal/scheduler"
	"billing-service/internal/subscriptions"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		observability.Must(err)
	}
	logger := observability.NewLogger(cfg.LogLevel)
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("database connection failed", zap.Error(err))
	}
	defer pool.Close()

	jobStore := queue.NewJobStore(pool)
	client, err := queue.NewClient(cfg.QueueURL, cfg.QueueStream, jobStore)
	if err != nil {
		logger.Fatal("queue connection failed", zap.Error(err))
	}
	defer client.Close()

	planSvc := plans.NewService(pool)
	subSvc := subscriptions.NewService(pool, planSvc)
	s := scheduler.NewScheduler(client, subSvc, logger)
	if err := s.Start(context.Background(), cfg.SchedulerInterval); err != nil {
		logger.Fatal("scheduler failed to start", zap.Error(err))
	}

	logger.Info("scheduler started", zap.String("interval", cfg.SchedulerInterval))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.Stop()
	logger.Info("scheduler shutting down")
}
