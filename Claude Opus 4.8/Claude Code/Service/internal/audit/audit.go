package audit

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Entry describes an auditable action.
type Entry struct {
	ActorUserID  *uuid.UUID
	Action       string
	ResourceType string
	ResourceID   string
	Metadata     map[string]any
	IPAddress    string
	UserAgent    string
}

// Recorder records audit entries.
type Recorder interface {
	Record(ctx context.Context, e Entry)
}

// Logger writes audit entries to the audit_logs table. Failures are logged,
// never propagated, so auditing cannot break the main request flow.
type Logger struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewLogger(db *pgxpool.Pool, log *slog.Logger) *Logger {
	return &Logger{db: db, log: log}
}

func (l *Logger) Record(ctx context.Context, e Entry) {
	meta := e.Metadata
	if meta == nil {
		meta = map[string]any{}
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		metaJSON = []byte("{}")
	}
	_, err = l.db.Exec(ctx, `
		INSERT INTO audit_logs (actor_user_id, action, resource_type, resource_id, metadata, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		e.ActorUserID, e.Action, e.ResourceType, e.ResourceID, metaJSON, e.IPAddress, e.UserAgent)
	if err != nil {
		l.log.Error("failed to write audit log", "action", e.Action, "error", err.Error())
	}
}

// Common action names.
const (
	ActionUserLogin            = "user.login"
	ActionPlanCreated          = "plan.created"
	ActionSubscriptionCreated  = "subscription.created"
	ActionSubscriptionCanceled = "subscription.cancelled"
	ActionInvoiceGenerated     = "invoice.generated"
	ActionInvoicePaid          = "invoice.paid"
	ActionPaymentFailed        = "payment.failed"
	ActionCouponApplied        = "coupon.applied"
	ActionTaxRuleUpdated       = "tax_rule.updated"
)
