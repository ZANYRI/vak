package invoices

import (
	"context"
	"errors"
	"time"

	"github.com/example/billing-service/internal/audit"
	"github.com/example/billing-service/internal/billing"
	"github.com/example/billing-service/internal/customers"
	"github.com/example/billing-service/internal/httpx"
	"github.com/example/billing-service/internal/observability"
	"github.com/example/billing-service/internal/plans"
	"github.com/example/billing-service/internal/queue"
	"github.com/example/billing-service/internal/subscriptions"
	"github.com/google/uuid"
)

// Reader interfaces decouple the integrator from concrete services.
type SubscriptionReader interface {
	Get(ctx context.Context, id uuid.UUID) (*subscriptions.Subscription, error)
}
type PlanReader interface {
	Get(ctx context.Context, id uuid.UUID) (*plans.Plan, error)
}
type CustomerReader interface {
	Get(ctx context.Context, id uuid.UUID) (*customers.Customer, error)
}
type UsageReader interface {
	UsageForPeriod(ctx context.Context, subscriptionID uuid.UUID, metric string, start, end time.Time) (int64, error)
}
type CouponReader interface {
	BillingCouponsForSubscription(ctx context.Context, subscriptionID uuid.UUID, currency string) ([]billing.Coupon, error)
}
type TaxReader interface {
	RateFor(ctx context.Context, country, region string) int
}
type Publisher interface {
	Publish(ctx context.Context, jobType string, payload any, opts queue.PublishOptions) (string, error)
}

type Service struct {
	repo     *Repository
	subs     SubscriptionReader
	plans    PlanReader
	custs    CustomerReader
	usage    UsageReader
	coupons  CouponReader
	taxes    TaxReader
	pub      Publisher
	audit    audit.Recorder
	metrics  *observability.Metrics
}

func NewService(repo *Repository, subs SubscriptionReader, planReader PlanReader, custs CustomerReader,
	usage UsageReader, coupons CouponReader, taxes TaxReader, pub Publisher,
	recorder audit.Recorder, metrics *observability.Metrics) *Service {
	return &Service{repo: repo, subs: subs, plans: planReader, custs: custs, usage: usage,
		coupons: coupons, taxes: taxes, pub: pub, audit: recorder, metrics: metrics}
}

func actorPtr(actor uuid.UUID) *uuid.UUID {
	if actor == uuid.Nil {
		return nil
	}
	return &actor
}

// Generate builds and persists an invoice for a subscription's current period.
func (s *Service) Generate(ctx context.Context, actor uuid.UUID, subscriptionID uuid.UUID) (*Invoice, error) {
	sub, err := s.subs.Get(ctx, subscriptionID)
	if err != nil {
		return nil, httpx.ErrValidation("subscription does not exist")
	}
	plan, err := s.plans.Get(ctx, sub.PlanID)
	if err != nil {
		return nil, httpx.ErrInternal("plan missing for subscription")
	}
	cust, err := s.custs.Get(ctx, sub.CustomerID)
	if err != nil {
		return nil, httpx.ErrInternal("customer missing for subscription")
	}

	currency := cust.Currency
	if currency == "" {
		currency = plan.Currency
	}

	var usageUnits int64
	if needsUsage(plan.PricingModel) && plan.UsageMetric != "" {
		usageUnits, _ = s.usage.UsageForPeriod(ctx, sub.ID, plan.UsageMetric, sub.CurrentPeriodStart, sub.CurrentPeriodEnd)
	}

	couponInputs, _ := s.coupons.BillingCouponsForSubscription(ctx, sub.ID, currency)
	taxBp := s.taxes.RateFor(ctx, cust.Country, cust.Region)

	in := billing.Input{
		Plan:               toBillingPlan(plan),
		Currency:           currency,
		Seats:              sub.Quantity,
		UsageUnits:         usageUnits,
		Coupons:            couponInputs,
		TaxRateBasisPoints: taxBp,
	}
	res := billing.Calculate(in)

	inv := &Invoice{
		CustomerID:     sub.CustomerID,
		SubscriptionID: &sub.ID,
		Status:         StatusDraft,
		Currency:       currency,
		SubtotalCents:  res.SubtotalCents,
		DiscountCents:  res.DiscountCents,
		TaxCents:       res.TaxCents,
		TotalCents:     res.TotalCents,
		AmountDueCents: res.TotalCents,
		PeriodStart:    &sub.CurrentPeriodStart,
		PeriodEnd:      &sub.CurrentPeriodEnd,
		Lines:          toInvoiceLines(res.Lines),
	}
	if err := s.repo.Create(ctx, inv); err != nil {
		return nil, err
	}

	if s.metrics != nil {
		s.metrics.InvoicesGenerated.Inc()
	}
	s.audit.Record(ctx, audit.Entry{
		ActorUserID: actorPtr(actor), Action: audit.ActionInvoiceGenerated,
		ResourceType: "invoice", ResourceID: inv.ID.String(),
		Metadata: map[string]any{"total_cents": inv.TotalCents, "currency": currency},
	})
	if s.pub != nil {
		_, _ = s.pub.Publish(ctx, queue.JobEmailInvoiceCreated,
			map[string]string{"invoice_id": inv.ID.String()}, queue.PublishOptions{})
	}
	return inv, nil
}

// CreateProration persists an invoice from pre-computed proration lines.
// It implements subscriptions.Invoicer.
func (s *Service) CreateProration(ctx context.Context, customerID, subscriptionID uuid.UUID, currency string, lines []billing.Line) (uuid.UUID, error) {
	var subtotal int64
	for _, l := range lines {
		subtotal += l.AmountCents
	}
	if subtotal < 0 {
		subtotal = 0
	}
	now := time.Now()
	due := now.AddDate(0, 0, 7)
	inv := &Invoice{
		CustomerID:     customerID,
		SubscriptionID: &subscriptionID,
		Status:         StatusOpen,
		Currency:       currency,
		SubtotalCents:  subtotal,
		TotalCents:     subtotal,
		AmountDueCents: subtotal,
		IssuedAt:       &now,
		DueAt:          &due,
		Lines:          toInvoiceLines(lines),
	}
	if err := s.repo.Create(ctx, inv); err != nil {
		return uuid.Nil, err
	}
	if s.metrics != nil {
		s.metrics.InvoicesGenerated.Inc()
	}
	return inv.ID, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	inv, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return nil, httpx.ErrNotFound("invoice not found")
	}
	return inv, err
}

func (s *Service) List(ctx context.Context, customerID *uuid.UUID, limit, offset int) ([]Invoice, error) {
	return s.repo.List(ctx, customerID, limit, offset)
}

func (s *Service) Finalize(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	inv, err := s.repo.Finalize(ctx, id, time.Now().AddDate(0, 0, 7))
	if errors.Is(err, ErrNotFound) {
		return nil, httpx.ErrConflict("invoice is not a finalizable draft")
	}
	return inv, err
}

func (s *Service) Void(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	inv, err := s.repo.Void(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return nil, httpx.ErrConflict("invoice cannot be voided")
	}
	return inv, err
}

// MarkPaid applies a successful payment to an invoice (called by payments).
func (s *Service) MarkPaid(ctx context.Context, id uuid.UUID, amountCents int64) (*Invoice, error) {
	inv, err := s.repo.MarkPaid(ctx, id, amountCents)
	if errors.Is(err, ErrNotFound) {
		return nil, httpx.ErrNotFound("invoice not found")
	}
	if err != nil {
		return nil, err
	}
	if inv.Status == StatusPaid {
		s.audit.Record(ctx, audit.Entry{
			Action: audit.ActionInvoicePaid, ResourceType: "invoice", ResourceID: inv.ID.String(),
			Metadata: map[string]any{"amount_cents": amountCents},
		})
	}
	return inv, nil
}

// OverdueInvoices returns open invoices past their due date (scheduler use).
func (s *Service) OverdueInvoices(ctx context.Context, limit int) ([]Invoice, error) {
	return s.repo.Overdue(ctx, time.Now(), limit)
}

// MarkUncollectible flags an invoice as uncollectible.
func (s *Service) MarkUncollectible(ctx context.Context, id uuid.UUID) error {
	return s.repo.SetStatus(ctx, id, StatusUncollectible)
}

func needsUsage(model string) bool {
	return model == billing.ModelUsageBased || model == billing.ModelTiered || model == billing.ModelHybrid
}

func toBillingPlan(p *plans.Plan) billing.Plan {
	bp := billing.Plan{
		PricingModel:   p.PricingModel,
		Currency:       p.Currency,
		BasePriceCents: p.BasePriceCents,
		SeatPriceCents: p.SeatPriceCents,
		IncludedSeats:  p.IncludedSeats,
		IncludedUnits:  p.IncludedUnits,
		UnitPriceCents: p.UnitPriceCents,
	}
	for _, t := range p.Tiers {
		bp.Tiers = append(bp.Tiers, billing.Tier{From: t.FromUnit, To: t.ToUnit, UnitPriceCents: t.UnitPriceCents})
	}
	return bp
}

func toInvoiceLines(lines []billing.Line) []Line {
	out := make([]Line, 0, len(lines))
	for _, l := range lines {
		out = append(out, Line{
			Type:            l.Type,
			Description:     l.Description,
			Quantity:        l.Quantity,
			UnitAmountCents: l.UnitAmountCents,
			AmountCents:     l.AmountCents,
		})
	}
	return out
}
