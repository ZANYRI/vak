package workers

import (
	"context"
	"fmt"
	"time"

	"billing-service/internal/invoices"
	"billing-service/internal/payments"
	"billing-service/internal/queue"
	"billing-service/internal/subscriptions"
	"billing-service/internal/usage"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Registry maps queue names to job handlers.
type Registry struct {
	handlers map[string]func(map[string]interface{}) error
}

// NewRegistry creates a worker registry wired to domain services.
func NewRegistry(
	logger *zap.Logger,
	subs *subscriptions.Service,
	inv *invoices.Service,
	pay *payments.Service,
	usg *usage.Service,
	queueClient *queue.Client,
) *Registry {
	r := &Registry{handlers: make(map[string]func(map[string]interface{}) error)}

	r.handlers[queue.QueueInvoiceGenerate] = func(p map[string]interface{}) error {
		logger.Info("processing invoice.generate", zap.Any("payload", p))
		return nil
	}
	r.handlers[queue.QueueInvoiceFinalize] = func(p map[string]interface{}) error {
		id, err := parsePayloadUUID(p, "invoice_id")
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
		_, err = inv.Finalize(ctx, id)
		return err
	}
	r.handlers[queue.QueuePaymentProcess] = func(p map[string]interface{}) error {
		id, err := parsePayloadUUID(p, "invoice_id")
		if err != nil {
			return err
		}
		card := stringValue(p["card_number"])
		if card == "" {
			return fmt.Errorf("card_number required")
		}
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
		_, err = pay.Simulate(ctx, payments.SimulateRequest{InvoiceID: id, CardNumber: card}, "")
		return err
	}
	r.handlers[queue.QueueSubscriptionRenew] = func(p map[string]interface{}) error {
		id, err := parsePayloadUUID(p, "subscription_id")
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
		_, err = subs.Resume(ctx, id)
		return err
	}
	r.handlers[queue.QueueExpireTrial] = func(p map[string]interface{}) error {
		id, err := parsePayloadUUID(p, "subscription_id")
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
		_, err = subs.Resume(ctx, id)
		return err
	}
	r.handlers[queue.QueueEmailInvoice] = noop(logger, "email.invoice_created")
	r.handlers[queue.QueueEmailPaymentFail] = noop(logger, "email.payment_failed")
	r.handlers[queue.QueueUsageAggregate] = func(p map[string]interface{}) error {
		logger.Info("aggregating usage", zap.Any("payload", p))
		return nil
	}
	_ = usg
	_ = queueClient
	return r
}

const defaultTimeout = 30 * time.Second

// Handler returns the handler for a queue.
func (r *Registry) Handler(queue string) func(map[string]interface{}) error {
	return r.handlers[queue]
}

// Queues returns registered queue names.
func (r *Registry) Queues() []string {
	queues := make([]string, 0, len(r.handlers))
	for q := range r.handlers {
		queues = append(queues, q)
	}
	return queues
}

func noop(logger *zap.Logger, name string) func(map[string]interface{}) error {
	return func(p map[string]interface{}) error {
		logger.Info("processed email job", zap.String("queue", name), zap.Any("payload", p))
		return nil
	}
}

func stringValue(v interface{}) string {
	s, _ := v.(string)
	return s
}

func parsePayloadUUID(p map[string]interface{}, key string) (uuid.UUID, error) {
	s := stringValue(p[key])
	if s == "" {
		return uuid.Nil, fmt.Errorf("%s required", key)
	}
	return uuid.Parse(s)
}
