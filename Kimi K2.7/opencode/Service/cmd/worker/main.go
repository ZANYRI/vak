package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"billing-service/internal/config"
	"billing-service/internal/coupons"
	"billing-service/internal/customers"
	"billing-service/internal/database"
	"billing-service/internal/invoices"
	"billing-service/internal/observability"
	"billing-service/internal/payments"
	"billing-service/internal/plans"
	"billing-service/internal/queue"
	"billing-service/internal/subscriptions"
	"billing-service/internal/taxes"
	"billing-service/internal/usage"
	"billing-service/internal/workers"

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
	customerSvc := customers.NewService(pool)
	usageSvc := usage.NewService(pool)
	couponSvc := coupons.NewService(pool)
	taxSvc := taxes.NewService(pool)
	invoiceSvc := invoices.NewService(pool, customerSvc, subSvc, planSvc, usageSvc, couponSvc, taxSvc)
	paymentSvc := payments.NewService(pool, invoiceSvc, subSvc)

	registry := workers.NewRegistry(logger, subSvc, invoiceSvc, paymentSvc, usageSvc, client)
	w := workers.NewWorker(registry, client, logger)
	if err := w.Run(); err != nil {
		logger.Fatal("worker failed to start", zap.Error(err))
	}

	logger.Info("worker started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("worker shutting down")
}
