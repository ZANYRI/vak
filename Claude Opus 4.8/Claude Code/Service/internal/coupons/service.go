package coupons

import (
	"context"
	"errors"
	"time"

	"github.com/example/billing-service/internal/audit"
	"github.com/example/billing-service/internal/billing"
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

// Validate checks whether a coupon may be applied to an invoice in the given
// currency at time now. It returns httpx errors describing the failure.
func (s *Service) Validate(c *Coupon, currency string, now time.Time) error {
	if !c.IsActive {
		return httpx.ErrConflict("coupon is not active")
	}
	if (c.ValidFrom != nil && now.Before(*c.ValidFrom)) || (c.ValidUntil != nil && now.After(*c.ValidUntil)) {
		return httpx.ErrConflict("coupon is expired or not yet valid")
	}
	if c.MaxRedemptions != nil && c.TimesRedeemed >= *c.MaxRedemptions {
		return httpx.ErrConflict("coupon redemption limit reached")
	}
	switch c.Type {
	case "percentage":
		if c.PercentOff == nil || *c.PercentOff < 1 || *c.PercentOff > 100 {
			return httpx.ErrValidation("coupon percent_off must be between 1 and 100")
		}
	case "fixed_amount":
		if c.Currency == nil || *c.Currency != currency {
			return httpx.ErrValidation("coupon currency must match invoice currency")
		}
	}
	return nil
}

func (s *Service) Create(ctx context.Context, actor uuid.UUID, req CreateRequest) (*Coupon, error) {
	c := &Coupon{
		Code:           req.Code,
		Type:           req.Type,
		PercentOff:     req.PercentOff,
		AmountOffCents: req.AmountOffCents,
		Currency:       req.Currency,
		MaxRedemptions: req.MaxRedemptions,
		TimesRedeemed:  0,
		ValidFrom:      req.ValidFrom,
		ValidUntil:     req.ValidUntil,
		IsActive:       true,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	s.audit.Record(ctx, audit.Entry{
		ActorUserID: &actor, Action: audit.ActionCouponApplied,
		ResourceType: "coupon", ResourceID: c.ID.String(),
		Metadata: map[string]any{"code": c.Code, "type": c.Type},
	})
	return c, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Coupon, error) {
	c, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return nil, httpx.ErrNotFound("coupon not found")
	}
	return c, err
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]Coupon, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Coupon, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, httpx.ErrNotFound("coupon not found")
		}
		return nil, err
	}
	return s.repo.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return httpx.ErrNotFound("coupon not found")
	}
	return err
}

// BillingCouponsForSubscription returns the valid coupons applied to a
// subscription, converted to billing.Coupon inputs for invoice calculation.
// Invalid/expired coupons are silently skipped.
func (s *Service) BillingCouponsForSubscription(ctx context.Context, subscriptionID uuid.UUID, currency string) ([]billing.Coupon, error) {
	applied, err := s.repo.ListForSubscription(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	var out []billing.Coupon
	for i := range applied {
		c := &applied[i]
		if s.Validate(c, currency, now) != nil {
			continue
		}
		bc := billing.Coupon{Type: c.Type}
		if c.PercentOff != nil {
			bc.PercentOff = *c.PercentOff
		}
		if c.AmountOffCents != nil {
			bc.AmountOffCents = *c.AmountOffCents
		}
		if c.Currency != nil {
			bc.Currency = *c.Currency
		}
		out = append(out, bc)
	}
	return out, nil
}

// ApplyCoupon validates a coupon by code and applies it to a subscription,
// incrementing its redemption count and recording an audit entry.
func (s *Service) ApplyCoupon(ctx context.Context, actor uuid.UUID, subscriptionID uuid.UUID, code string, invoiceCurrency string) (*Coupon, error) {
	c, err := s.repo.GetByCode(ctx, code)
	if errors.Is(err, ErrNotFound) {
		return nil, httpx.ErrNotFound("coupon not found")
	}
	if err != nil {
		return nil, err
	}
	if err := s.Validate(c, invoiceCurrency, time.Now()); err != nil {
		return nil, err
	}
	if err := s.repo.ApplyToSubscription(ctx, subscriptionID, c.ID); err != nil {
		if errors.Is(err, ErrAlreadyApplied) {
			return nil, httpx.ErrConflict("coupon already applied")
		}
		return nil, err
	}
	if err := s.repo.IncrementRedemption(ctx, c.ID); err != nil {
		return nil, err
	}
	c.TimesRedeemed++
	s.audit.Record(ctx, audit.Entry{
		ActorUserID: &actor, Action: audit.ActionCouponApplied,
		ResourceType: "subscription", ResourceID: subscriptionID.String(),
		Metadata: map[string]any{"coupon_id": c.ID.String(), "code": c.Code},
	})
	return c, nil
}
