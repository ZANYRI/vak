package payments

import (
	"context"
	"errors"
	"math/rand"
	"strings"

	"github.com/example/billing-service/internal/audit"
	"github.com/example/billing-service/internal/httpx"
	"github.com/example/billing-service/internal/invoices"
	"github.com/example/billing-service/internal/observability"
	"github.com/example/billing-service/internal/queue"
	"github.com/google/uuid"
)

// InvoiceStore is the subset of the invoices service used by payments.
type InvoiceStore interface {
	Get(ctx context.Context, id uuid.UUID) (*invoices.Invoice, error)
	MarkPaid(ctx context.Context, id uuid.UUID, amountCents int64) (*invoices.Invoice, error)
}

// SubscriptionMarker flags a subscription past_due after a failed payment.
type SubscriptionMarker interface {
	MarkPastDue(ctx context.Context, id uuid.UUID) error
}

// Publisher enqueues notification jobs.
type Publisher interface {
	Publish(ctx context.Context, jobType string, payload any, opts queue.PublishOptions) (string, error)
}

type Service struct {
	repo     *Repository
	invoices InvoiceStore
	subs     SubscriptionMarker
	pub      Publisher
	audit    audit.Recorder
	metrics  *observability.Metrics
	rng      *rand.Rand
}

func NewService(repo *Repository, inv InvoiceStore, subs SubscriptionMarker, pub Publisher,
	recorder audit.Recorder, metrics *observability.Metrics) *Service {
	return &Service{repo: repo, invoices: inv, subs: subs, pub: pub, audit: recorder,
		metrics: metrics, rng: rand.New(rand.NewSource(1))}
}

// decideOutcome applies the mock provider rules from service.md.
func (s *Service) decideOutcome(cardNumber string) (success bool, reason string) {
	switch {
	case strings.HasSuffix(cardNumber, "0000"):
		return true, ""
	case strings.HasSuffix(cardNumber, "9999"):
		return false, "card declined"
	default:
		if s.rng.Intn(2) == 0 {
			return false, "card declined"
		}
		return true, ""
	}
}

func last4(card string) string {
	if len(card) <= 4 {
		return card
	}
	return card[len(card)-4:]
}

// Simulate processes a mock payment for an invoice's outstanding amount.
func (s *Service) Simulate(ctx context.Context, req SimulateRequest) (*Payment, error) {
	inv, err := s.invoices.Get(ctx, req.InvoiceID)
	if err != nil {
		return nil, httpx.ErrNotFound("invoice not found")
	}
	if inv.Status == invoices.StatusPaid {
		return nil, httpx.ErrConflict("invoice is already paid")
	}
	if inv.Status == invoices.StatusVoid {
		return nil, httpx.ErrConflict("invoice is void")
	}

	amount := inv.AmountDueCents
	if amount <= 0 {
		amount = inv.TotalCents
	}

	success, reason := s.decideOutcome(req.CardNumber)
	p := &Payment{
		InvoiceID:   inv.ID,
		AmountCents: amount,
		Currency:    inv.Currency,
		CardLast4:   last4(req.CardNumber),
	}
	if success {
		p.Status = StatusSucceeded
	} else {
		p.Status = StatusFailed
		p.FailureReason = reason
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	if success {
		if _, err := s.invoices.MarkPaid(ctx, inv.ID, amount); err != nil {
			return nil, err
		}
		if s.metrics != nil {
			s.metrics.PaymentsSucceeded.Inc()
		}
		return p, nil
	}

	// Failure path: mark subscription past_due, notify, audit.
	if inv.SubscriptionID != nil && s.subs != nil {
		_ = s.subs.MarkPastDue(ctx, *inv.SubscriptionID)
	}
	if s.metrics != nil {
		s.metrics.PaymentsFailed.Inc()
	}
	s.audit.Record(ctx, audit.Entry{
		Action: audit.ActionPaymentFailed, ResourceType: "invoice", ResourceID: inv.ID.String(),
		Metadata: map[string]any{"reason": reason, "amount_cents": amount},
	})
	if s.pub != nil {
		_, _ = s.pub.Publish(ctx, queue.JobEmailPaymentFailed,
			map[string]string{"invoice_id": inv.ID.String()}, queue.PublishOptions{})
	}
	return p, httpx.ErrPaymentFailed(reason)
}

// PayInvoice simulates a guaranteed-success payment for the invoice (used by
// the POST /invoices/{id}/pay convenience endpoint).
func (s *Service) PayInvoice(ctx context.Context, invoiceID uuid.UUID) (*Payment, error) {
	return s.Simulate(ctx, SimulateRequest{InvoiceID: invoiceID, CardNumber: "4242424242420000"})
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Payment, error) {
	p, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return nil, httpx.ErrNotFound("payment not found")
	}
	return p, err
}

func (s *Service) List(ctx context.Context, invoiceID *uuid.UUID, limit, offset int) ([]Payment, error) {
	return s.repo.List(ctx, invoiceID, limit, offset)
}
