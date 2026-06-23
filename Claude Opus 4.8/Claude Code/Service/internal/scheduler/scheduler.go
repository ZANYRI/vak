package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/example/billing-service/internal/invoices"
	"github.com/example/billing-service/internal/queue"
	"github.com/example/billing-service/internal/subscriptions"
)

// Scheduler periodically finds work and publishes jobs to the queue. It does
// not perform heavy work itself (see service.md "Scheduler").
type Scheduler struct {
	log      *slog.Logger
	subs     *subscriptions.Service
	invoices *invoices.Service
	pub      *queue.Queue
	interval time.Duration
	batch    int
}

func New(log *slog.Logger, subs *subscriptions.Service, inv *invoices.Service, pub *queue.Queue, interval time.Duration) *Scheduler {
	return &Scheduler{log: log, subs: subs, invoices: inv, pub: pub, interval: interval, batch: 100}
}

// Run ticks until the context is cancelled (graceful shutdown).
func (s *Scheduler) Run(ctx context.Context) error {
	s.log.Info("scheduler started", "interval", s.interval.String())
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.tick(ctx) // run once immediately
	for {
		select {
		case <-ctx.Done():
			s.log.Info("scheduler shutting down")
			return nil
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	now := time.Now()
	s.enqueueRenewals(ctx, now)
	s.enqueueTrialExpirations(ctx, now)
	s.retryOverdue(ctx)
	s.aggregateUsage(ctx)
}

func (s *Scheduler) enqueueRenewals(ctx context.Context, now time.Time) {
	subs, err := s.subs.DueForRenewal(ctx, now, s.batch)
	if err != nil {
		s.log.Error("query renewals failed", "error", err.Error())
		return
	}
	for _, sub := range subs {
		_, err := s.pub.Publish(ctx, queue.JobSubscriptionRenew,
			map[string]string{"subscription_id": sub.ID.String()},
			queue.PublishOptions{IdempotencyKey: "renew:" + sub.ID.String() + ":" + sub.CurrentPeriodEnd.Format(time.RFC3339)})
		if err != nil && err != queue.ErrDuplicate {
			s.log.Error("publish renew failed", "subscription_id", sub.ID, "error", err.Error())
		}
	}
	if len(subs) > 0 {
		s.log.Info("queued renewals", "count", len(subs))
	}
}

func (s *Scheduler) enqueueTrialExpirations(ctx context.Context, now time.Time) {
	subs, err := s.subs.ExpiredTrials(ctx, now, s.batch)
	if err != nil {
		s.log.Error("query expired trials failed", "error", err.Error())
		return
	}
	for _, sub := range subs {
		_, err := s.pub.Publish(ctx, queue.JobExpireTrial,
			map[string]string{"subscription_id": sub.ID.String()},
			queue.PublishOptions{IdempotencyKey: "trial:" + sub.ID.String()})
		if err != nil && err != queue.ErrDuplicate {
			s.log.Error("publish expire_trial failed", "subscription_id", sub.ID, "error", err.Error())
		}
	}
	if len(subs) > 0 {
		s.log.Info("queued trial expirations", "count", len(subs))
	}
}

func (s *Scheduler) retryOverdue(ctx context.Context) {
	invs, err := s.invoices.OverdueInvoices(ctx, s.batch)
	if err != nil {
		s.log.Error("query overdue invoices failed", "error", err.Error())
		return
	}
	for _, inv := range invs {
		_, err := s.pub.Publish(ctx, queue.JobPaymentProcess,
			map[string]string{"invoice_id": inv.ID.String()}, queue.PublishOptions{})
		if err != nil && err != queue.ErrDuplicate {
			s.log.Error("publish payment retry failed", "invoice_id", inv.ID, "error", err.Error())
		}
	}
	if len(invs) > 0 {
		s.log.Info("queued overdue payment retries", "count", len(invs))
	}
}

func (s *Scheduler) aggregateUsage(ctx context.Context) {
	_, err := s.pub.Publish(ctx, queue.JobUsageAggregate, map[string]string{"at": time.Now().Format(time.RFC3339)}, queue.PublishOptions{})
	if err != nil && err != queue.ErrDuplicate {
		s.log.Error("publish usage aggregate failed", "error", err.Error())
	}
}
