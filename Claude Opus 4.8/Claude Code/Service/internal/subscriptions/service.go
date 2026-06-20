package subscriptions

import (
	"context"
	"errors"
	"time"

	"github.com/example/billing-service/internal/audit"
	"github.com/example/billing-service/internal/billing"
	"github.com/example/billing-service/internal/httpx"
	"github.com/example/billing-service/internal/plans"
	"github.com/example/billing-service/internal/queue"
	"github.com/google/uuid"
)

// PlanReader fetches plans (satisfied by plans.Service).
type PlanReader interface {
	Get(ctx context.Context, id uuid.UUID) (*plans.Plan, error)
}

// Publisher enqueues background jobs (satisfied by queue.Queue).
type Publisher interface {
	Publish(ctx context.Context, jobType string, payload any, opts queue.PublishOptions) (string, error)
}

// Invoicer creates proration invoices (satisfied by invoices.Service).
type Invoicer interface {
	CreateProration(ctx context.Context, customerID, subscriptionID uuid.UUID, currency string, lines []billing.Line) (uuid.UUID, error)
}

type Service struct {
	repo     *Repository
	plans    PlanReader
	invoicer Invoicer
	pub      Publisher
	audit    audit.Recorder
}

func NewService(repo *Repository, planReader PlanReader, invoicer Invoicer, pub Publisher, recorder audit.Recorder) *Service {
	return &Service{repo: repo, plans: planReader, invoicer: invoicer, pub: pub, audit: recorder}
}

func addInterval(t time.Time, interval string) time.Time {
	if interval == "yearly" {
		return t.AddDate(1, 0, 0)
	}
	return t.AddDate(0, 1, 0)
}

func (s *Service) Create(ctx context.Context, actor uuid.UUID, req CreateRequest) (*Subscription, error) {
	plan, err := s.plans.Get(ctx, req.PlanID)
	if err != nil {
		return nil, httpx.ErrValidation("plan does not exist")
	}
	now := time.Now()
	qty := req.Quantity
	if qty < 1 {
		qty = 1
	}
	sub := &Subscription{
		CustomerID:         req.CustomerID,
		PlanID:             req.PlanID,
		Quantity:           qty,
		CurrentPeriodStart: now,
	}
	if plan.TrialDays > 0 {
		trialEnd := now.AddDate(0, 0, plan.TrialDays)
		sub.Status = StatusTrialing
		sub.TrialStart = &now
		sub.TrialEnd = &trialEnd
		sub.CurrentPeriodEnd = trialEnd
	} else {
		sub.Status = StatusActive
		sub.CurrentPeriodEnd = addInterval(now, plan.BillingInterval)
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}

	s.audit.Record(ctx, audit.Entry{
		ActorUserID: &actor, Action: audit.ActionSubscriptionCreated,
		ResourceType: "subscription", ResourceID: sub.ID.String(),
		Metadata: map[string]any{"plan_id": sub.PlanID.String(), "status": sub.Status},
	})

	// For immediately-active subscriptions, queue the first invoice.
	if sub.Status == StatusActive && s.pub != nil {
		_, _ = s.pub.Publish(ctx, queue.JobInvoiceGenerate,
			map[string]string{"subscription_id": sub.ID.String()}, queue.PublishOptions{})
	}
	return sub, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return nil, httpx.ErrNotFound("subscription not found")
	}
	return sub, err
}

func (s *Service) List(ctx context.Context, customerID *uuid.UUID, limit, offset int) ([]Subscription, error) {
	return s.repo.List(ctx, customerID, limit, offset)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Subscription, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, req)
}

func (s *Service) Cancel(ctx context.Context, actor uuid.UUID, id uuid.UUID, atPeriodEnd bool) (*Subscription, error) {
	sub, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if atPeriodEnd {
		if err := s.repo.SetCancelAtPeriodEnd(ctx, id, true); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.SetStatus(ctx, id, StatusCancelled); err != nil {
			return nil, err
		}
	}
	s.audit.Record(ctx, audit.Entry{
		ActorUserID: &actor, Action: audit.ActionSubscriptionCanceled,
		ResourceType: "subscription", ResourceID: id.String(),
		Metadata: map[string]any{"at_period_end": atPeriodEnd},
	})
	_ = sub
	return s.Get(ctx, id)
}

func (s *Service) Pause(ctx context.Context, id uuid.UUID) (*Subscription, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return nil, err
	}
	if err := s.repo.SetStatus(ctx, id, StatusPaused); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

func (s *Service) Resume(ctx context.Context, id uuid.UUID) (*Subscription, error) {
	sub, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub.Status != StatusPaused {
		return nil, httpx.ErrConflict("subscription is not paused")
	}
	if err := s.repo.SetStatus(ctx, id, StatusActive); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Renew advances the subscription to its next billing period (used by the
// renewal worker). Returns the updated subscription.
func (s *Service) Renew(ctx context.Context, id uuid.UUID) (*Subscription, error) {
	sub, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	plan, err := s.plans.Get(ctx, sub.PlanID)
	if err != nil {
		return nil, httpx.ErrInternal("plan missing for subscription")
	}
	newStart := sub.CurrentPeriodEnd
	newEnd := addInterval(newStart, plan.BillingInterval)
	if err := s.repo.AdvancePeriod(ctx, id, newStart, newEnd); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// ExpireTrial ends a trial and starts the first paid period.
func (s *Service) ExpireTrial(ctx context.Context, id uuid.UUID) (*Subscription, error) {
	sub, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	plan, err := s.plans.Get(ctx, sub.PlanID)
	if err != nil {
		return nil, httpx.ErrInternal("plan missing for subscription")
	}
	now := time.Now()
	if err := s.repo.AdvancePeriod(ctx, id, now, addInterval(now, plan.BillingInterval)); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// DueForRenewal / ExpiredTrials expose scheduler queries.
func (s *Service) DueForRenewal(ctx context.Context, now time.Time, limit int) ([]Subscription, error) {
	return s.repo.DueForRenewal(ctx, now, limit)
}
func (s *Service) ExpiredTrials(ctx context.Context, now time.Time, limit int) ([]Subscription, error) {
	return s.repo.ExpiredTrials(ctx, now, limit)
}

// MarkPastDue flags a subscription as past_due (used after a failed payment).
func (s *Service) MarkPastDue(ctx context.Context, id uuid.UUID) error {
	return s.repo.SetStatus(ctx, id, StatusPastDue)
}

// SetInvoicer wires the proration invoicer after construction to break the
// subscriptions<->invoices construction cycle.
func (s *Service) SetInvoicer(i Invoicer) { s.invoicer = i }

// ChangePlanResult bundles the updated subscription and the proration outcome.
type ChangePlanResult struct {
	Subscription *Subscription     `json:"subscription"`
	Proration    billing.Proration `json:"proration"`
	InvoiceID    *uuid.UUID        `json:"invoice_id,omitempty"`
}

func (s *Service) ChangePlan(ctx context.Context, id uuid.UUID, req ChangePlanRequest) (*ChangePlanResult, error) {
	sub, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	oldPlan, err := s.plans.Get(ctx, sub.PlanID)
	if err != nil {
		return nil, httpx.ErrInternal("current plan missing")
	}
	newPlan, err := s.plans.Get(ctx, req.PlanID)
	if err != nil {
		return nil, httpx.ErrValidation("target plan does not exist")
	}
	if oldPlan.Currency != newPlan.Currency {
		return nil, httpx.ErrValidation("cannot change to a plan with a different currency")
	}

	pror := billing.Prorate(oldPlan.BasePriceCents, newPlan.BasePriceCents,
		sub.CurrentPeriodStart, sub.CurrentPeriodEnd, time.Now())

	updated, err := s.repo.ChangePlan(ctx, id, req.PlanID)
	if err != nil {
		return nil, err
	}

	result := &ChangePlanResult{Subscription: updated, Proration: pror}
	if s.invoicer != nil && len(pror.Lines) > 0 {
		invID, err := s.invoicer.CreateProration(ctx, sub.CustomerID, id, newPlan.Currency, pror.Lines)
		if err == nil {
			result.InvoiceID = &invID
		}
	}
	return result, nil
}
