package scheduler

import (
	"context"
	"fmt"
	"time"

	"billing-service/internal/observability"
	"billing-service/internal/queue"
	"billing-service/internal/subscriptions"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Scheduler periodically publishes background jobs.
type Scheduler struct {
	cron    *cron.Cron
	queue   *queue.Client
	subs    *subscriptions.Service
	logger  *zap.Logger
}

// NewScheduler creates a scheduler.
func NewScheduler(queue *queue.Client, subs *subscriptions.Service, logger *zap.Logger) *Scheduler {
	return &Scheduler{queue: queue, subs: subs, logger: logger}
}

// Start begins the cron scheduler.
func (s *Scheduler) Start(ctx context.Context, schedule string) error {
	if schedule == "" {
		schedule = "*/1 * * * *"
	}
	s.cron = cron.New()
	_, err := s.cron.AddFunc(schedule, func() {
		s.logger.Info("scheduler tick")
		s.enqueueRenewals(ctx)
		s.enqueueTrialExpirations(ctx)
		s.enqueueUsageAggregations(ctx)
	})
	if err != nil {
		return fmt.Errorf("invalid scheduler expression: %w", err)
	}
	s.cron.Start()
	return nil
}

// Stop halts the scheduler.
func (s *Scheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
	}
}

// Wait blocks until context cancellation.
func (s *Scheduler) Wait(ctx context.Context) {
	<-ctx.Done()
}

func (s *Scheduler) enqueueRenewals(ctx context.Context) {
	list, _, err := s.subs.List(ctx, nil, 1000, 0)
	if err != nil {
		s.logger.Warn("failed to list subscriptions for renewal", zap.Error(err))
		return
	}
	now := time.Now()
	for _, sub := range list {
		if sub.Status == "active" && sub.CurrentPeriodEnd.Before(now) {
			_, _ = s.queue.Publish(ctx, queue.QueueSubscriptionRenew, map[string]interface{}{"subscription_id": sub.ID.String()}, 3)
			observability.IncJobProcessed(queue.QueueSubscriptionRenew, "queued")
		}
	}
}

func (s *Scheduler) enqueueTrialExpirations(ctx context.Context) {
	list, _, err := s.subs.List(ctx, nil, 1000, 0)
	if err != nil {
		s.logger.Warn("failed to list subscriptions for trial expiry", zap.Error(err))
		return
	}
	now := time.Now()
	for _, sub := range list {
		if sub.Status == "trialing" && sub.TrialEnd != nil && sub.TrialEnd.Before(now) {
			_, _ = s.queue.Publish(ctx, queue.QueueExpireTrial, map[string]interface{}{"subscription_id": sub.ID.String()}, 3)
			observability.IncJobProcessed(queue.QueueExpireTrial, "queued")
		}
	}
}

func (s *Scheduler) enqueueUsageAggregations(ctx context.Context) {
	list, _, err := s.subs.List(ctx, nil, 1000, 0)
	if err != nil {
		s.logger.Warn("failed to list subscriptions for usage aggregation", zap.Error(err))
		return
	}
	for _, sub := range list {
		_, _ = s.queue.Publish(ctx, queue.QueueUsageAggregate, map[string]interface{}{"subscription_id": sub.ID.String()}, 3)
		observability.IncJobProcessed(queue.QueueUsageAggregate, "queued")
	}
}
