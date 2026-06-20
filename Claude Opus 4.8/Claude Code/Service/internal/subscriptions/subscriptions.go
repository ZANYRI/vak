package subscriptions

import (
	"time"

	"github.com/google/uuid"
)

// Subscription statuses.
const (
	StatusTrialing  = "trialing"
	StatusActive    = "active"
	StatusPastDue   = "past_due"
	StatusPaused    = "paused"
	StatusCancelled = "cancelled"
	StatusExpired   = "expired"
)

// Subscription is a customer's enrollment in a plan.
type Subscription struct {
	ID                 uuid.UUID  `json:"id"`
	CustomerID         uuid.UUID  `json:"customer_id"`
	PlanID             uuid.UUID  `json:"plan_id"`
	Status             string     `json:"status"`
	Quantity           int64      `json:"quantity"`
	CurrentPeriodStart time.Time  `json:"current_period_start"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end"`
	TrialStart         *time.Time `json:"trial_start"`
	TrialEnd           *time.Time `json:"trial_end"`
	CancelAtPeriodEnd  bool       `json:"cancel_at_period_end"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// CreateRequest is the create-subscription payload.
type CreateRequest struct {
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	PlanID     uuid.UUID `json:"plan_id" validate:"required"`
	Quantity   int64     `json:"quantity" validate:"gte=1"`
}

// UpdateRequest patches mutable fields.
type UpdateRequest struct {
	Quantity          *int64 `json:"quantity"`
	CancelAtPeriodEnd *bool  `json:"cancel_at_period_end"`
}

// ChangePlanRequest switches a subscription to a different plan with proration.
type ChangePlanRequest struct {
	PlanID uuid.UUID `json:"plan_id" validate:"required"`
}
