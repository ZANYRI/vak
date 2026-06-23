package usage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"billing-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service provides usage event management.
type Service struct {
	pool *pgxpool.Pool
}

// NewService creates a usage service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// ReportRequest submits a usage event.
type ReportRequest struct {
	CustomerID     uuid.UUID `json:"customer_id" validate:"required"`
	SubscriptionID uuid.UUID `json:"subscription_id" validate:"required"`
	Metric         string    `json:"metric" validate:"required"`
	Quantity       int64     `json:"quantity" validate:"min=1"`
	IdempotencyKey string    `json:"idempotency_key" validate:"required"`
	RecordedAt     *time.Time `json:"recorded_at,omitempty"`
}

// Report inserts a usage event idempotently.
func (s *Service) Report(ctx context.Context, req ReportRequest) (*models.UsageEvent, error) {
	existing, err := s.findByIdempotency(ctx, req.IdempotencyKey)
	if err == nil {
		return existing, nil
	}
	recorded := time.Now().UTC()
	if req.RecordedAt != nil {
		recorded = *req.RecordedAt
	}
	event := &models.UsageEvent{
		ID:             uuid.New(),
		CustomerID:     req.CustomerID,
		SubscriptionID: req.SubscriptionID,
		Metric:         req.Metric,
		Quantity:       req.Quantity,
		IdempotencyKey: req.IdempotencyKey,
		RecordedAt:     recorded,
	}
	_, err = s.pool.Exec(ctx,
		`INSERT INTO usage_events (id, customer_id, subscription_id, metric, quantity, idempotency_key, recorded_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		event.ID, event.CustomerID, event.SubscriptionID, event.Metric, event.Quantity, event.IdempotencyKey, event.RecordedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return s.findByIdempotency(ctx, req.IdempotencyKey)
		}
		return nil, fmt.Errorf("insert usage event: %w", err)
	}
	return event, nil
}

// List returns usage events for a subscription.
func (s *Service) List(ctx context.Context, subscriptionID uuid.UUID, limit, offset int) ([]models.UsageEvent, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := s.pool.Query(ctx,
		`SELECT id, customer_id, subscription_id, metric, quantity, idempotency_key, recorded_at, created_at
		 FROM usage_events WHERE subscription_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		subscriptionID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.UsageEvent
	for rows.Next() {
		var e models.UsageEvent
		if err := rows.Scan(&e.ID, &e.CustomerID, &e.SubscriptionID, &e.Metric, &e.Quantity, &e.IdempotencyKey, &e.RecordedAt, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, e)
	}
	var total int64
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM usage_events WHERE subscription_id = $1`, subscriptionID).Scan(&total)
	return list, total, rows.Err()
}

// Summary aggregates usage for a subscription and metric within a period.
func (s *Service) Summary(ctx context.Context, subscriptionID uuid.UUID, metric string, from, to time.Time) (*models.UsageSummary, error) {
	var total int64
	err := s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(quantity), 0) FROM usage_events
		 WHERE subscription_id = $1 AND metric = $2 AND recorded_at >= $3 AND recorded_at < $4`,
		subscriptionID, metric, from, to).Scan(&total)
	if err != nil {
		return nil, err
	}
	return &models.UsageSummary{
		SubscriptionID: subscriptionID,
		Metric:         metric,
		PeriodStart:    from,
		PeriodEnd:      to,
		TotalQuantity:  total,
	}, nil
}

func (s *Service) findByIdempotency(ctx context.Context, key string) (*models.UsageEvent, error) {
	var e models.UsageEvent
	err := s.pool.QueryRow(ctx,
		`SELECT id, customer_id, subscription_id, metric, quantity, idempotency_key, recorded_at, created_at
		 FROM usage_events WHERE idempotency_key = $1`, key).Scan(
		&e.ID, &e.CustomerID, &e.SubscriptionID, &e.Metric, &e.Quantity, &e.IdempotencyKey, &e.RecordedAt, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
