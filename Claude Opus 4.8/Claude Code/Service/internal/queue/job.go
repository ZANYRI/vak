package queue

import "encoding/json"

// Job type constants (see service.md "Queue Broker and Workers").
const (
	JobInvoiceGenerate     = "invoice.generate"
	JobInvoiceFinalize     = "invoice.finalize"
	JobPaymentProcess      = "payment.process"
	JobSubscriptionRenew   = "subscription.renew"
	JobExpireTrial         = "subscription.expire_trial"
	JobEmailInvoiceCreated = "email.invoice_created"
	JobEmailPaymentFailed  = "email.payment_failed"
	JobUsageAggregate      = "usage.aggregate"
)

// Job is a unit of background work.
type Job struct {
	ID             string          `json:"id"` // jobs table id
	Type           string          `json:"type"`
	Payload        json.RawMessage `json:"payload"`
	Attempts       int             `json:"attempts"`
	MaxAttempts    int             `json:"max_attempts"`
	IdempotencyKey string          `json:"idempotency_key,omitempty"`
}

// PublishOptions tune how a job is enqueued.
type PublishOptions struct {
	IdempotencyKey string
	MaxAttempts    int
}
