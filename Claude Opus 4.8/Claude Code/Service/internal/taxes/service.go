package taxes

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

func (s *Service) Create(ctx context.Context, actor uuid.UUID, req CreateRequest) (*TaxRule, error) {
	t := &TaxRule{
		Country:            req.Country,
		Region:             req.Region,
		TaxName:            req.TaxName,
		TaxRateBasisPoints: req.TaxRateBasisPoints,
		IsActive:           true,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	s.audit.Record(ctx, audit.Entry{
		ActorUserID: &actor, Action: audit.ActionTaxRuleUpdated,
		ResourceType: "tax_rule", ResourceID: t.ID.String(),
		Metadata: map[string]any{"country": t.Country, "region": t.Region},
	})
	return t, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*TaxRule, error) {
	t, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return nil, httpx.ErrNotFound("tax rule not found")
	}
	return t, err
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]TaxRule, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) Update(ctx context.Context, actor uuid.UUID, id uuid.UUID, req UpdateRequest) (*TaxRule, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, httpx.ErrNotFound("tax rule not found")
		}
		return nil, err
	}
	t, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	s.audit.Record(ctx, audit.Entry{
		ActorUserID: &actor, Action: audit.ActionTaxRuleUpdated,
		ResourceType: "tax_rule", ResourceID: id.String(),
	})
	return t, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return httpx.ErrNotFound("tax rule not found")
	}
	return err
}

// RateFor returns the active tax rate in basis points for (country, region),
// or 0 when no active rule is found.
func (s *Service) RateFor(ctx context.Context, country, region string) int {
	t, err := s.repo.FindActiveFor(ctx, country, region)
	if err != nil {
		return 0
	}
	return t.TaxRateBasisPoints
}
