package customers

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

func (s *Service) Create(ctx context.Context, actor uuid.UUID, req CreateRequest) (*Customer, error) {
	c := &Customer{
		UserID:         req.UserID,
		Email:          req.Email,
		Name:           req.Name,
		CompanyName:    req.CompanyName,
		BillingAddress: req.BillingAddress,
		Country:        req.Country,
		Region:         req.Region,
		TaxID:          req.TaxID,
		Currency:       req.Currency,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	s.audit.Record(ctx, audit.Entry{
		ActorUserID:  &actor,
		Action:       "customer.created",
		ResourceType: "customer",
		ResourceID:   c.ID.String(),
		Metadata:     map[string]any{"email": c.Email, "name": c.Name},
	})
	return c, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Customer, error) {
	c, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return nil, httpx.ErrNotFound("customer not found")
	}
	return c, err
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]Customer, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Customer, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, httpx.ErrNotFound("customer not found")
		}
		return nil, err
	}
	return s.repo.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return httpx.ErrNotFound("customer not found")
	}
	return err
}
