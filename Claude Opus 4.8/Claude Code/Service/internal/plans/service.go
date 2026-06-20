package plans

import (
	"context"
	"errors"

	"github.com/example/billing-service/internal/audit"
	"github.com/example/billing-service/internal/httpx"
	"github.com/google/uuid"
)

type Service struct {
	repo  *Repository
	audit audit.Recorder
}

func NewService(repo *Repository, recorder audit.Recorder) *Service {
	return &Service{repo: repo, audit: recorder}
}

func (s *Service) Create(ctx context.Context, actor uuid.UUID, req CreateRequest) (*Plan, error) {
	p := &Plan{
		Name:            req.Name,
		Description:     req.Description,
		Currency:        req.Currency,
		BillingInterval: req.BillingInterval,
		PricingModel:    req.PricingModel,
		BasePriceCents:  req.BasePriceCents,
		SeatPriceCents:  req.SeatPriceCents,
		IncludedSeats:   req.IncludedSeats,
		IncludedUnits:   req.IncludedUnits,
		UnitPriceCents:  req.UnitPriceCents,
		UsageMetric:     req.UsageMetric,
		TrialDays:       req.TrialDays,
		IsActive:        true,
	}
	for _, t := range req.Tiers {
		p.Tiers = append(p.Tiers, Tier{FromUnit: t.From, ToUnit: t.To, UnitPriceCents: t.UnitPriceCents})
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	s.audit.Record(ctx, audit.Entry{
		ActorUserID: &actor, Action: audit.ActionPlanCreated,
		ResourceType: "plan", ResourceID: p.ID.String(),
		Metadata: map[string]any{"name": p.Name, "model": p.PricingModel},
	})
	return p, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Plan, error) {
	p, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return nil, httpx.ErrNotFound("plan not found")
	}
	return p, err
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]Plan, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Plan, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, httpx.ErrNotFound("plan not found")
		}
		return nil, err
	}
	return s.repo.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return httpx.ErrNotFound("plan not found")
	}
	return err
}
