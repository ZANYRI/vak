package models

import (
	"time"

	"github.com/google/uuid"
)

// Customer represents a billing customer.
type Customer struct {
	ID             uuid.UUID  `json:"id"`
	UserID         *uuid.UUID `json:"user_id,omitempty"`
	Email          string     `json:"email"`
	Name           string     `json:"name"`
	CompanyName    string     `json:"company_name,omitempty"`
	BillingAddress string     `json:"billing_address,omitempty"`
	Country        string     `json:"country"`
	Region         string     `json:"region,omitempty"`
	TaxID          string     `json:"tax_id,omitempty"`
	Currency       string     `json:"currency"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// Plan represents a subscription plan.
type Plan struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description,omitempty"`
	Currency        string          `json:"currency"`
	BillingInterval BillingInterval `json:"billing_interval"`
	BasePriceCents  int64           `json:"base_price_cents"`
	TrialDays       int             `json:"trial_days"`
	IsActive        bool            `json:"is_active"`
	Prices          []PlanPrice     `json:"prices,omitempty"`
	Tiers           []PricingTier   `json:"tiers,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// PlanPrice holds additional pricing configuration for a plan.
type PlanPrice struct {
	ID              uuid.UUID `json:"id"`
	PlanID          uuid.UUID `json:"plan_id"`
	Model           PlanModel `json:"model"`
	SeatPriceCents  *int64    `json:"seat_price_cents,omitempty"`
	IncludedSeats   *int      `json:"included_seats,omitempty"`
	UnitPriceCents  *int64    `json:"unit_price_cents,omitempty"`
	IncludedUnits   *int64    `json:"included_units,omitempty"`
}

// PricingTier represents a tiered price band.
type PricingTier struct {
	ID             uuid.UUID `json:"id"`
	PlanID         uuid.UUID `json:"plan_id"`
	From           int64     `json:"from"`
	To             *int64    `json:"to,omitempty"`
	UnitPriceCents int64     `json:"unit_price_cents"`
}

// Subscription links a customer to a plan.
type Subscription struct {
	ID                   uuid.UUID          `json:"id"`
	CustomerID           uuid.UUID          `json:"customer_id"`
	PlanID               uuid.UUID          `json:"plan_id"`
	Plan                 *Plan              `json:"plan,omitempty"`
	Status               SubscriptionStatus `json:"status"`
	Quantity             int                `json:"quantity"`
	Seats                int                `json:"seats"`
	CurrentPeriodStart   time.Time          `json:"current_period_start"`
	CurrentPeriodEnd     time.Time          `json:"current_period_end"`
	TrialStart           *time.Time         `json:"trial_start,omitempty"`
	TrialEnd             *time.Time         `json:"trial_end,omitempty"`
	CancelAtPeriodEnd    bool               `json:"cancel_at_period_end"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

// UsageEvent records a consumption event.
type UsageEvent struct {
	ID             uuid.UUID `json:"id"`
	CustomerID     uuid.UUID `json:"customer_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Metric         string    `json:"metric"`
	Quantity       int64     `json:"quantity"`
	IdempotencyKey string    `json:"idempotency_key"`
	RecordedAt     time.Time `json:"recorded_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// UsageSummary aggregates usage for a billing period.
type UsageSummary struct {
	ID            uuid.UUID `json:"id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Metric        string    `json:"metric"`
	PeriodStart   time.Time `json:"period_start"`
	PeriodEnd     time.Time `json:"period_end"`
	TotalQuantity int64     `json:"total_quantity"`
}

// Coupon represents a discount coupon.
type Coupon struct {
	ID              uuid.UUID  `json:"id"`
	Code            string     `json:"code"`
	Type            CouponType `json:"type"`
	PercentOff      *int       `json:"percent_off,omitempty"`
	AmountOffCents  *int64     `json:"amount_off_cents,omitempty"`
	Currency        *string    `json:"currency,omitempty"`
	MaxRedemptions  *int       `json:"max_redemptions,omitempty"`
	TimesRedeemed   int        `json:"times_redeemed"`
	ValidFrom       time.Time  `json:"valid_from"`
	ValidUntil      *time.Time `json:"valid_until,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TaxRule defines a simple tax rule.
type TaxRule struct {
	ID                   uuid.UUID `json:"id"`
	Country              string    `json:"country"`
	Region               string    `json:"region,omitempty"`
	TaxName              string    `json:"tax_name"`
	TaxRateBasisPoints   int       `json:"tax_rate_basis_points"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Invoice represents a billing invoice.
type Invoice struct {
	ID             uuid.UUID     `json:"id"`
	CustomerID     uuid.UUID     `json:"customer_id"`
	SubscriptionID *uuid.UUID    `json:"subscription_id,omitempty"`
	Status         InvoiceStatus `json:"status"`
	Currency       string        `json:"currency"`
	SubtotalCents  int64         `json:"subtotal_cents"`
	DiscountCents  int64         `json:"discount_cents"`
	TaxCents       int64         `json:"tax_cents"`
	TotalCents     int64         `json:"total_cents"`
	AmountDueCents int64         `json:"amount_due_cents"`
	AmountPaidCents int64        `json:"amount_paid_cents"`
	PeriodStart    *time.Time    `json:"period_start,omitempty"`
	PeriodEnd      *time.Time    `json:"period_end,omitempty"`
	IssuedAt       *time.Time    `json:"issued_at,omitempty"`
	DueAt          *time.Time    `json:"due_at,omitempty"`
	PaidAt         *time.Time    `json:"paid_at,omitempty"`
	Lines          []InvoiceLine `json:"lines,omitempty"`
	IdempotencyKey *string       `json:"-"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// InvoiceLine is a line item of an invoice.
type InvoiceLine struct {
	ID              uuid.UUID              `json:"id"`
	InvoiceID       uuid.UUID              `json:"invoice_id"`
	Type            string                 `json:"type"`
	Description     string                 `json:"description"`
	Quantity        int64                  `json:"quantity"`
	UnitAmountCents int64                  `json:"unit_amount_cents"`
	AmountCents     int64                  `json:"amount_cents"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

// Payment records a payment attempt.
type Payment struct {
	ID            uuid.UUID     `json:"id"`
	InvoiceID     uuid.UUID     `json:"invoice_id"`
	Status        PaymentStatus `json:"status"`
	AmountCents   int64         `json:"amount_cents"`
	Currency      string        `json:"currency"`
	CardLast4     string        `json:"card_last4,omitempty"`
	FailureReason string        `json:"failure_reason,omitempty"`
	IdempotencyKey *string      `json:"-"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// Job is a persisted background job.
type Job struct {
	ID           uuid.UUID `json:"id"`
	Queue        string    `json:"queue"`
	Payload      map[string]interface{} `json:"payload"`
	Status       JobStatus `json:"status"`
	Attempts     int       `json:"attempts"`
	MaxAttempts  int       `json:"max_attempts"`
	RunAfter     time.Time `json:"run_after"`
	ErrorMessage string    `json:"error_message,omitempty"`
	DeadAt       *time.Time `json:"dead_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AuditLog records important actions.
type AuditLog struct {
	ID           uuid.UUID              `json:"id"`
	ActorUserID  *uuid.UUID             `json:"actor_user_id,omitempty"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   *uuid.UUID             `json:"resource_id,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	IPAddress    string                 `json:"ip_address,omitempty"`
	UserAgent    string                 `json:"user_agent,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}
