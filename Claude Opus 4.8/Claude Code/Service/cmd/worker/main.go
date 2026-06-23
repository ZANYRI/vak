// Command worker consumes background jobs from the queue.
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/example/billing-service/internal/app"
	"github.com/example/billing-service/internal/config"
	"github.com/example/billing-service/internal/observability"
	"github.com/example/billing-service/internal/queue"
	"github.com/example/billing-service/internal/workers"
	"github.com/google/uuid"
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

	worker := workers.New(log, container.Invoices, container.Payments, container.Subs)

	var wg sync.WaitGroup
	for i := 0; i < cfg.WorkerConcurrency; i++ {
		consumerID := "worker-" + uuid.NewString()[:8]
		consumer := queue.NewConsumer(container.Queue, log, container.Metrics, consumerID)
		worker.Register(consumer)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := consumer.Run(ctx); err != nil {
				log.Error("consumer stopped with error", "error", err.Error())
			}
		}()
	}

	// Periodically publish queue depth as a gauge.
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if depth, err := container.Queue.QueueDepth(ctx); err == nil {
					container.Metrics.QueueDepth.WithLabelValues("billing:jobs").Set(float64(depth))
				}
			}
		}
	}()

	log.Info("worker pool started", "concurrency", cfg.WorkerConcurrency)
	<-ctx.Done()
	log.Info("worker shutdown signal received")
	wg.Wait()
	log.Info("worker stopped")
}

func cfgEnv(cfg *config.Config) string {
	if cfg == nil {
		return "production"
	}
	return cfg.AppEnv
}
