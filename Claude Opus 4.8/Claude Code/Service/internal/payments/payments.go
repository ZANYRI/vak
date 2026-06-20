package payments

import (
	"time"

	"github.com/google/uuid"
)

// Payment statuses.
const (
	StatusPending   = "pending"
	StatusSucceeded = "succeeded"
	StatusFailed    = "failed"
	StatusRefunded  = "refunded"
)

// Payment is a (simulated) payment attempt against an invoice.
type Payment struct {
	ID            uuid.UUID `json:"id"`
	InvoiceID     uuid.UUID `json:"invoice_id"`
	AmountCents   int64     `json:"amount_cents"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	CardLast4     string    `json:"card_last4"`
	FailureReason string    `json:"failure_reason"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SimulateRequest triggers a mock payment for an invoice.
type SimulateRequest struct {
	InvoiceID  uuid.UUID `json:"invoice_id" validate:"required"`
	CardNumber string    `json:"card_number" validate:"required,min=4"`
}
