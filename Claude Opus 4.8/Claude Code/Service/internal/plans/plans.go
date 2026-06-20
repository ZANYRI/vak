package plans

import (
	"time"

	"github.com/google/uuid"
)

// Plan is a billing plan.
type Plan struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Currency        string    `json:"currency"`
	BillingInterval string    `json:"billing_interval"`
	PricingModel    string    `json:"pricing_model"`
	BasePriceCents  int64     `json:"base_price_cents"`
	SeatPriceCents  int64     `json:"seat_price_cents"`
	IncludedSeats   int64     `json:"included_seats"`
	IncludedUnits   int64     `json:"included_units"`
	UnitPriceCents  int64     `json:"unit_price_cents"`
	UsageMetric     string    `json:"usage_metric"`
	TrialDays       int       `json:"trial_days"`
	IsActive        bool      `json:"is_active"`
	Tiers           []Tier    `json:"tiers,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Tier is a pricing band for tiered/hybrid plans.
type Tier struct {
	ID             uuid.UUID `json:"id"`
	FromUnit       int64     `json:"from"`
	ToUnit         *int64    `json:"to"`
	UnitPriceCents int64     `json:"unit_price_cents"`
	SortOrder      int       `json:"sort_order"`
}

// CreateRequest is the create-plan payload.
type CreateRequest struct {
	Name            string       `json:"name" validate:"required"`
	Description     string       `json:"description"`
	Currency        string       `json:"currency" validate:"required,len=3"`
	BillingInterval string       `json:"billing_interval" validate:"required,oneof=monthly yearly"`
	PricingModel    string       `json:"pricing_model" validate:"required,oneof=flat per_seat usage_based tiered hybrid"`
	BasePriceCents  int64        `json:"base_price_cents" validate:"gte=0"`
	SeatPriceCents  int64        `json:"seat_price_cents" validate:"gte=0"`
	IncludedSeats   int64        `json:"included_seats" validate:"gte=0"`
	IncludedUnits   int64        `json:"included_units" validate:"gte=0"`
	UnitPriceCents  int64        `json:"unit_price_cents" validate:"gte=0"`
	UsageMetric     string       `json:"usage_metric"`
	TrialDays       int          `json:"trial_days" validate:"gte=0"`
	Tiers           []TierInput  `json:"tiers" validate:"dive"`
}

// TierInput is a tier in a create/update request.
type TierInput struct {
	From           int64  `json:"from" validate:"gte=0"`
	To             *int64 `json:"to"`
	UnitPriceCents int64  `json:"unit_price_cents" validate:"gte=0"`
}

// UpdateRequest patches mutable plan fields.
type UpdateRequest struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	BasePriceCents *int64  `json:"base_price_cents"`
	SeatPriceCents *int64  `json:"seat_price_cents"`
	IncludedSeats  *int64  `json:"included_seats"`
	IncludedUnits  *int64  `json:"included_units"`
	UnitPriceCents *int64  `json:"unit_price_cents"`
	TrialDays      *int    `json:"trial_days"`
	IsActive       *bool   `json:"is_active"`
}
