package coupons

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"billing-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service provides coupon management.
type Service struct {
	pool *pgxpool.Pool
}

// NewService creates a coupon service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// CreateRequest represents a coupon creation payload.
type CreateRequest struct {
	Code           string     `json:"code" validate:"required"`
	Type           string     `json:"type" validate:"required,oneof=percentage fixed_amount"`
	PercentOff     *int       `json:"percent_off,omitempty"`
	AmountOffCents *int64     `json:"amount_off_cents,omitempty"`
	Currency       *string    `json:"currency,omitempty"`
	MaxRedemptions *int       `json:"max_redemptions,omitempty"`
	ValidFrom      *time.Time `json:"valid_from,omitempty"`
	ValidUntil     *time.Time `json:"valid_until,omitempty"`
	IsActive       bool       `json:"is_active"`
}

// ValidateCoupon validates the coupon payload.
func ValidateCoupon(req CreateRequest) error {
	t := models.CouponType(req.Type)
	if t == models.CouponPercentage {
		if req.PercentOff == nil || *req.PercentOff < 1 || *req.PercentOff > 100 {
			return fmt.Errorf("percentage coupons must have percent_off between 1 and 100")
		}
	}
	if t == models.CouponFixedAmount {
		if req.AmountOffCents == nil || *req.AmountOffCents <= 0 {
			return fmt.Errorf("fixed_amount coupons must have amount_off_cents > 0")
		}
	}
	return nil
}

// Create inserts a coupon.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*models.Coupon, error) {
	if err := ValidateCoupon(req); err != nil {
		return nil, err
	}
	c := &models.Coupon{
		ID:        uuid.New(),
		Code:      strings.ToUpper(req.Code),
		Type:      models.CouponType(req.Type),
		PercentOff: req.PercentOff,
		AmountOffCents: req.AmountOffCents,
		MaxRedemptions: req.MaxRedemptions,
		IsActive:  req.IsActive,
	}
	validFrom := time.Now().UTC()
	if req.ValidFrom != nil {
		validFrom = *req.ValidFrom
	}
	c.ValidFrom = validFrom
	c.ValidUntil = req.ValidUntil
	if req.Currency != nil {
		cur := strings.ToUpper(*req.Currency)
		c.Currency = &cur
	}

	_, err := s.pool.Exec(ctx,
		`INSERT INTO coupons (id, code, type, percent_off, amount_off_cents, currency, max_redemptions, times_redeemed, valid_from, valid_until, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, 0, $8, $9, $10)`,
		c.ID, c.Code, c.Type, c.PercentOff, c.AmountOffCents, c.Currency, c.MaxRedemptions, c.ValidFrom, c.ValidUntil, c.IsActive)
	if err != nil {
		return nil, fmt.Errorf("insert coupon: %w", err)
	}
	return c, nil
}

// List returns coupons.
func (s *Service) List(ctx context.Context, limit, offset int) ([]models.Coupon, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := s.pool.Query(ctx,
		`SELECT id, code, type, percent_off, amount_off_cents, currency, max_redemptions, times_redeemed, valid_from, valid_until, is_active, created_at, updated_at
		 FROM coupons ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.Coupon
	for rows.Next() {
		var c models.Coupon
		err := rows.Scan(&c.ID, &c.Code, &c.Type, &c.PercentOff, &c.AmountOffCents, &c.Currency, &c.MaxRedemptions, &c.TimesRedeemed, &c.ValidFrom, &c.ValidUntil, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, c)
	}
	var total int64
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM coupons`).Scan(&total)
	return list, total, rows.Err()
}

// GetByID returns a coupon by id.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Coupon, error) {
	var c models.Coupon
	err := s.pool.QueryRow(ctx,
		`SELECT id, code, type, percent_off, amount_off_cents, currency, max_redemptions, times_redeemed, valid_from, valid_until, is_active, created_at, updated_at
		 FROM coupons WHERE id = $1`, id).Scan(
		&c.ID, &c.Code, &c.Type, &c.PercentOff, &c.AmountOffCents, &c.Currency, &c.MaxRedemptions, &c.TimesRedeemed, &c.ValidFrom, &c.ValidUntil, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("coupon not found")
		}
		return nil, err
	}
	return &c, nil
}

// GetByCode returns a coupon by code.
func (s *Service) GetByCode(ctx context.Context, code string) (*models.Coupon, error) {
	var c models.Coupon
	err := s.pool.QueryRow(ctx,
		`SELECT id, code, type, percent_off, amount_off_cents, currency, max_redemptions, times_redeemed, valid_from, valid_until, is_active, created_at, updated_at
		 FROM coupons WHERE code = $1`, strings.ToUpper(code)).Scan(
		&c.ID, &c.Code, &c.Type, &c.PercentOff, &c.AmountOffCents, &c.Currency, &c.MaxRedemptions, &c.TimesRedeemed, &c.ValidFrom, &c.ValidUntil, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("coupon not found")
		}
		return nil, err
	}
	return &c, nil
}

// UpdateRequest updates mutable fields of a coupon.
type UpdateRequest struct {
	IsActive *bool `json:"is_active,omitempty"`
}

// Update modifies a coupon.
func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*models.Coupon, error) {
	if req.IsActive == nil {
		return s.GetByID(ctx, id)
	}
	_, err := s.pool.Exec(ctx, `UPDATE coupons SET is_active=$1, updated_at=NOW() WHERE id=$2`, *req.IsActive, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// Delete removes a coupon.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM coupons WHERE id = $1`, id)
	return err
}

// ApplyCouponRequest applies a coupon to a subscription.
type ApplyCouponRequest struct {
	CouponID uuid.UUID `json:"coupon_id" validate:"required"`
}

// ApplyCoupon links a coupon to a subscription and increments redemption count.
func (s *Service) ApplyCoupon(ctx context.Context, subscriptionID uuid.UUID, couponID uuid.UUID) error {
	coupon, err := s.GetByID(ctx, couponID)
	if err != nil {
		return err
	}
	if !coupon.IsActive {
		return fmt.Errorf("coupon is inactive")
	}
	if coupon.ValidFrom.After(time.Now()) || (coupon.ValidUntil != nil && coupon.ValidUntil.Before(time.Now())) {
		return fmt.Errorf("coupon is expired or not yet valid")
	}
	if coupon.MaxRedemptions != nil && coupon.TimesRedeemed >= *coupon.MaxRedemptions {
		return fmt.Errorf("coupon max redemptions reached")
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		`INSERT INTO subscription_coupons (id, subscription_id, coupon_id) VALUES ($1, $2, $3)`,
		uuid.New(), subscriptionID, couponID)
	if err != nil {
		return fmt.Errorf("apply coupon: %w", err)
	}
	_, err = tx.Exec(ctx,
		`UPDATE coupons SET times_redeemed = times_redeemed + 1, updated_at=NOW() WHERE id=$1`, couponID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// GetAppliedCoupons returns coupons applied to a subscription.
func (s *Service) GetAppliedCoupons(ctx context.Context, subscriptionID uuid.UUID) ([]models.Coupon, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT c.id, c.code, c.type, c.percent_off, c.amount_off_cents, c.currency, c.max_redemptions, c.times_redeemed, c.valid_from, c.valid_until, c.is_active, c.created_at, c.updated_at
		 FROM coupons c JOIN subscription_coupons sc ON c.id = sc.coupon_id WHERE sc.subscription_id = $1`, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Coupon
	for rows.Next() {
		var c models.Coupon
		err := rows.Scan(&c.ID, &c.Code, &c.Type, &c.PercentOff, &c.AmountOffCents, &c.Currency, &c.MaxRedemptions, &c.TimesRedeemed, &c.ValidFrom, &c.ValidUntil, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}
