package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/example/billing-service/internal/invoices"
	"github.com/example/billing-service/internal/payments"
	"github.com/example/billing-service/internal/queue"
	"github.com/example/billing-service/internal/subscriptions"
	"github.com/google/uuid"
)

// Worker bundles the services needed to process background jobs.
type Worker struct {
	log      *slog.Logger
	invoices *invoices.Service
	payments *payments.Service
	subs     *subscriptions.Service
}

func New(log *slog.Logger, inv *invoices.Service, pay *payments.Service, subs *subscriptions.Service) *Worker {
	return &Worker{log: log, invoices: inv, payments: pay, subs: subs}
}

// Register binds all job handlers onto the consumer.
func (w *Worker) Register(c *queue.Consumer) {
	c.Register(queue.JobInvoiceGenerate, w.handleInvoiceGenerate)
	c.Register(queue.JobInvoiceFinalize, w.handleInvoiceFinalize)
	c.Register(queue.JobPaymentProcess, w.handlePaymentProcess)
	c.Register(queue.JobSubscriptionRenew, w.handleSubscriptionRenew)
	c.Register(queue.JobExpireTrial, w.handleExpireTrial)
	c.Register(queue.JobUsageAggregate, w.handleUsageAggregate)
	c.Register(queue.JobEmailInvoiceCreated, w.handleEmailInvoiceCreated)
	c.Register(queue.JobEmailPaymentFailed, w.handleEmailPaymentFailed)
}

func decodeID(payload json.RawMessage, field string) (uuid.UUID, error) {
	var m map[string]string
	if err := json.Unmarshal(payload, &m); err != nil {
		return uuid.Nil, err
	}
	id, err := uuid.Parse(m[field])
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s in payload", field)
	}
	return id, nil
}

func (w *Worker) handleInvoiceGenerate(ctx context.Context, job queue.Job) error {
	subID, err := decodeID(job.Payload, "subscription_id")
	if err != nil {
		return err
	}
	inv, err := w.invoices.Generate(ctx, uuid.Nil, subID)
	if err != nil {
		return err
	}
	// Auto-finalize generated invoices so they become payable.
	if _, err := w.invoices.Finalize(ctx, inv.ID); err != nil {
		w.log.Warn("invoice finalize after generate failed", "invoice_id", inv.ID, "error", err.Error())
	}
	w.log.Info("invoice generated", "invoice_id", inv.ID, "subscription_id", subID)
	return nil
}

func (w *Worker) handleInvoiceFinalize(ctx context.Context, job queue.Job) error {
	invID, err := decodeID(job.Payload, "invoice_id")
	if err != nil {
		return err
	}
	_, err = w.invoices.Finalize(ctx, invID)
	return err
}

func (w *Worker) handlePaymentProcess(ctx context.Context, job queue.Job) error {
	invID, err := decodeID(job.Payload, "invoice_id")
	if err != nil {
		return err
	}
	_, err = w.payments.PayInvoice(ctx, invID)
	// A declined payment is a terminal business outcome, not a transient job error.
	if err != nil {
		w.log.Info("payment processing result", "invoice_id", invID, "result", err.Error())
	}
	return nil
}

func (w *Worker) handleSubscriptionRenew(ctx context.Context, job queue.Job) error {
	subID, err := decodeID(job.Payload, "subscription_id")
	if err != nil {
		return err
	}
	if _, err := w.subs.Renew(ctx, subID); err != nil {
		return err
	}
	inv, err := w.invoices.Generate(ctx, uuid.Nil, subID)
	if err != nil {
		return err
	}
	if _, err := w.invoices.Finalize(ctx, inv.ID); err != nil {
		w.log.Warn("finalize after renew failed", "invoice_id", inv.ID, "error", err.Error())
	}
	w.log.Info("subscription renewed", "subscription_id", subID, "invoice_id", inv.ID)
	return nil
}

func (w *Worker) handleExpireTrial(ctx context.Context, job queue.Job) error {
	subID, err := decodeID(job.Payload, "subscription_id")
	if err != nil {
		return err
	}
	if _, err := w.subs.ExpireTrial(ctx, subID); err != nil {
		return err
	}
	inv, err := w.invoices.Generate(ctx, uuid.Nil, subID)
	if err != nil {
		return err
	}
	_, _ = w.invoices.Finalize(ctx, inv.ID)
	w.log.Info("trial expired, first invoice issued", "subscription_id", subID, "invoice_id", inv.ID)
	return nil
}

func (w *Worker) handleUsageAggregate(ctx context.Context, job queue.Job) error {
	w.log.Info("usage aggregation tick", "payload", string(job.Payload))
	return nil
}

func (w *Worker) handleEmailInvoiceCreated(ctx context.Context, job queue.Job) error {
	w.log.Info("📧 invoice created email (mock)", "payload", string(job.Payload))
	return nil
}

func (w *Worker) handleEmailPaymentFailed(ctx context.Context, job queue.Job) error {
	w.log.Info("📧 payment failed email (mock)", "payload", string(job.Payload))
	return nil
}
