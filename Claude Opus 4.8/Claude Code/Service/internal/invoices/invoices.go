package invoices

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Invoice statuses.
const (
	StatusDraft         = "draft"
	StatusOpen          = "open"
	StatusPaid          = "paid"
	StatusVoid          = "void"
	StatusUncollectible = "uncollectible"
)

// Invoice is a bill issued to a customer.
type Invoice struct {
	ID              uuid.UUID  `json:"id"`
	CustomerID      uuid.UUID  `json:"customer_id"`
	SubscriptionID  *uuid.UUID `json:"subscription_id"`
	Status          string     `json:"status"`
	Currency        string     `json:"currency"`
	SubtotalCents   int64      `json:"subtotal_cents"`
	DiscountCents   int64      `json:"discount_cents"`
	TaxCents        int64      `json:"tax_cents"`
	TotalCents      int64      `json:"total_cents"`
	AmountDueCents  int64      `json:"amount_due_cents"`
	AmountPaidCents int64      `json:"amount_paid_cents"`
	PeriodStart     *time.Time `json:"period_start"`
	PeriodEnd       *time.Time `json:"period_end"`
	IssuedAt        *time.Time `json:"issued_at"`
	DueAt           *time.Time `json:"due_at"`
	PaidAt          *time.Time `json:"paid_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Lines           []Line     `json:"lines"`
}

// Line is a single charge on an invoice.
type Line struct {
	ID              uuid.UUID       `json:"id"`
	InvoiceID       uuid.UUID       `json:"invoice_id"`
	Type            string          `json:"type"`
	Description     string          `json:"description"`
	Quantity        int64           `json:"quantity"`
	UnitAmountCents int64           `json:"unit_amount_cents"`
	AmountCents     int64           `json:"amount_cents"`
	Metadata        json.RawMessage `json:"metadata"`
	CreatedAt       time.Time       `json:"created_at"`
}

// GenerateRequest triggers invoice generation for a subscription.
type GenerateRequest struct {
	SubscriptionID uuid.UUID `json:"subscription_id" validate:"required"`
}
