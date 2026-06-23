package plans

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"billing-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service provides plan management.
type Service struct {
	pool *pgxpool.Pool
}

// NewService creates a plan service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// CreateRequest represents a plan creation payload.
type CreateRequest struct {
	Name            string            `json:"name" validate:"required"`
	Description     string            `json:"description,omitempty"`
	Currency        string            `json:"currency" validate:"required,len=3"`
	BillingInterval string            `json:"billing_interval" validate:"required,oneof=monthly yearly"`
	BasePriceCents  int64             `json:"base_price_cents" validate:"min=0"`
	TrialDays       int               `json:"trial_days" validate:"min=0"`
	IsActive        bool              `json:"is_active"`
	Prices          []PriceRequest    `json:"prices,omitempty"`
	Tiers           []TierRequest     `json:"tiers,omitempty"`
}

// PriceRequest describes a plan price configuration.
type PriceRequest struct {
	Model           string `json:"model" validate:"required,oneof=flat per_seat usage_based tiered hybrid"`
	SeatPriceCents  *int64 `json:"seat_price_cents,omitempty"`
	IncludedSeats   *int   `json:"included_seats,omitempty"`
	UnitPriceCents  *int64 `json:"unit_price_cents,omitempty"`
	IncludedUnits   *int64 `json:"included_units,omitempty"`
}

// TierRequest describes a pricing tier band.
type TierRequest struct {
	From           int64  `json:"from" validate:"min=0"`
	To             *int64 `json:"to,omitempty"`
	UnitPriceCents int64  `json:"unit_price_cents" validate:"min=0"`
}

// Create inserts a plan, its prices and tiers atomically.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*models.Plan, error) {
	plan := &models.Plan{
		ID:              uuid.New(),
		Name:            req.Name,
		Description:     req.Description,
		Currency:        strings.ToUpper(req.Currency),
		BillingInterval: models.BillingInterval(req.BillingInterval),
		BasePriceCents:  req.BasePriceCents,
		TrialDays:       req.TrialDays,
		IsActive:        req.IsActive,
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO plans (id, name, description, currency, billing_interval, base_price_cents, trial_days, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		plan.ID, plan.Name, plan.Description, plan.Currency, plan.BillingInterval, plan.BasePriceCents, plan.TrialDays, plan.IsActive)
	if err != nil {
		return nil, fmt.Errorf("insert plan: %w", err)
	}
	for _, p := range req.Prices {
		_, err = tx.Exec(ctx,
			`INSERT INTO plan_prices (id, plan_id, model, seat_price_cents, included_seats, unit_price_cents, included_units)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			uuid.New(), plan.ID, p.Model, p.SeatPriceCents, p.IncludedSeats, p.UnitPriceCents, p.IncludedUnits)
		if err != nil {
			return nil, fmt.Errorf("insert plan price: %w", err)
		}
	}
	for _, t := range req.Tiers {
		_, err = tx.Exec(ctx,
			`INSERT INTO pricing_tiers (id, plan_id, tier_from, tier_to, unit_price_cents)
			 VALUES ($1, $2, $3, $4, $5)`,
			uuid.New(), plan.ID, t.From, t.To, t.UnitPriceCents)
		if err != nil {
			return nil, fmt.Errorf("insert pricing tier: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, plan.ID)
}

// List returns all active/inactive plans.
func (s *Service) List(ctx context.Context, limit, offset int) ([]models.Plan, int64, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, description, currency, billing_interval, base_price_cents, trial_days, is_active, created_at, updated_at
		 FROM plans ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var plans []models.Plan
	for rows.Next() {
		var p models.Plan
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Currency, &p.BillingInterval, &p.BasePriceCents, &p.TrialDays, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		plans = append(plans, p)
	}
	var total int64
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM plans`).Scan(&total)
	for i := range plans {
		_ = s.loadPricesAndTiers(ctx, &plans[i])
	}
	return plans, total, rows.Err()
}

// GetByID returns a plan with prices and tiers.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Plan, error) {
	var p models.Plan
	err := s.pool.QueryRow(ctx,
		`SELECT id, name, description, currency, billing_interval, base_price_cents, trial_days, is_active, created_at, updated_at
		 FROM plans WHERE id = $1`, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Currency, &p.BillingInterval, &p.BasePriceCents, &p.TrialDays, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("plan not found")
		}
		return nil, err
	}
	if err := s.loadPricesAndTiers(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Service) loadPricesAndTiers(ctx context.Context, p *models.Plan) error {
	prices, err := s.pool.Query(ctx,
		`SELECT id, plan_id, model, seat_price_cents, included_seats, unit_price_cents, included_units
		 FROM plan_prices WHERE plan_id = $1`, p.ID)
	if err != nil {
		return err
	}
	defer prices.Close()
	p.Prices = nil
	for prices.Next() {
		var pp models.PlanPrice
		if err := prices.Scan(&pp.ID, &pp.PlanID, &pp.Model, &pp.SeatPriceCents, &pp.IncludedSeats, &pp.UnitPriceCents, &pp.IncludedUnits); err != nil {
			return err
		}
		p.Prices = append(p.Prices, pp)
	}

	tiers, err := s.pool.Query(ctx,
		`SELECT id, plan_id, tier_from, tier_to, unit_price_cents
		 FROM pricing_tiers WHERE plan_id = $1 ORDER BY tier_from`, p.ID)
	if err != nil {
		return err
	}
	defer tiers.Close()
	p.Tiers = nil
	for tiers.Next() {
		var t models.PricingTier
		if err := tiers.Scan(&t.ID, &t.PlanID, &t.From, &t.To, &t.UnitPriceCents); err != nil {
			return err
		}
		p.Tiers = append(p.Tiers, t)
	}
	return nil
}

// UpdateRequest updates mutable plan fields.
type UpdateRequest struct {
	Name            *string        `json:"name,omitempty"`
	Description     *string        `json:"description,omitempty"`
	Currency        *string        `json:"currency,omitempty" validate:"omitempty,len=3"`
	BillingInterval *string        `json:"billing_interval,omitempty" validate:"omitempty,oneof=monthly yearly"`
	BasePriceCents  *int64         `json:"base_price_cents,omitempty"`
	TrialDays       *int           `json:"trial_days,omitempty"`
	IsActive        *bool          `json:"is_active,omitempty"`
}

// Update modifies a plan.
func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*models.Plan, error) {
	p, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.Currency != nil {
		p.Currency = strings.ToUpper(*req.Currency)
	}
	if req.BillingInterval != nil {
		p.BillingInterval = models.BillingInterval(*req.BillingInterval)
	}
	if req.BasePriceCents != nil {
		p.BasePriceCents = *req.BasePriceCents
	}
	if req.TrialDays != nil {
		p.TrialDays = *req.TrialDays
	}
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}
	_, err = s.pool.Exec(ctx,
		`UPDATE plans SET name=$1, description=$2, currency=$3, billing_interval=$4, base_price_cents=$5, trial_days=$6, is_active=$7, updated_at=NOW() WHERE id=$8`,
		p.Name, p.Description, p.Currency, p.BillingInterval, p.BasePriceCents, p.TrialDays, p.IsActive, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// Delete marks a plan inactive or removes it (here: remove prices, tiers, then plan).
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := s.pool.Exec(ctx, `DELETE FROM plans WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("plan not found")
	}
	return nil
}
