package usage

import (
	"time"

	"github.com/google/uuid"
)

// Event is a reported usage measurement.
type Event struct {
	ID             uuid.UUID `json:"id"`
	CustomerID     uuid.UUID `json:"customer_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Metric         string    `json:"metric"`
	Quantity       int64     `json:"quantity"`
	IdempotencyKey string    `json:"idempotency_key"`
	RecordedAt     time.Time `json:"recorded_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// RecordRequest reports a usage event.
type RecordRequest struct {
	CustomerID     uuid.UUID  `json:"customer_id" validate:"required"`
	SubscriptionID uuid.UUID  `json:"subscription_id" validate:"required"`
	Metric         string     `json:"metric" validate:"required"`
	Quantity       int64      `json:"quantity" validate:"gt=0"`
	IdempotencyKey string     `json:"idempotency_key" validate:"required"`
	RecordedAt     *time.Time `json:"recorded_at"`
}

// Summary aggregates usage for a subscription/metric.
type Summary struct {
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Metric         string    `json:"metric"`
	TotalQuantity  int64     `json:"total_quantity"`
}
