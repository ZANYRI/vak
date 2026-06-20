package usage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

// Record stores a usage event idempotently. Re-submitting the same
// idempotency_key returns the original event without double counting.
func (s *Service) Record(ctx context.Context, req RecordRequest) (*Event, bool, error) {
	recordedAt := time.Now()
	if req.RecordedAt != nil {
		recordedAt = *req.RecordedAt
	}
	e := &Event{
		CustomerID:     req.CustomerID,
		SubscriptionID: req.SubscriptionID,
		Metric:         req.Metric,
		Quantity:       req.Quantity,
		IdempotencyKey: req.IdempotencyKey,
		RecordedAt:     recordedAt,
	}
	return s.repo.Insert(ctx, e)
}

func (s *Service) List(ctx context.Context, subscriptionID *uuid.UUID, limit, offset int) ([]Event, error) {
	return s.repo.List(ctx, subscriptionID, limit, offset)
}

func (s *Service) Summary(ctx context.Context, subscriptionID uuid.UUID) ([]Summary, error) {
	return s.repo.SummaryBySubscription(ctx, subscriptionID)
}

// UsageForPeriod is used by the billing engine to total metered usage.
func (s *Service) UsageForPeriod(ctx context.Context, subscriptionID uuid.UUID, metric string, start, end time.Time) (int64, error) {
	return s.repo.SumForPeriod(ctx, subscriptionID, metric, start, end)
}
