package coupons

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("coupon not found")
var ErrAlreadyApplied = errors.New("coupon already applied")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, c *Coupon) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO coupons (code, type, percent_off, amount_off_cents, currency,
			max_redemptions, times_redeemed, valid_from, valid_until, is_active)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, created_at, updated_at`,
		c.Code, c.Type, c.PercentOff, c.AmountOffCents, c.Currency,
		c.MaxRedemptions, c.TimesRedeemed, c.ValidFrom, c.ValidUntil, c.IsActive,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Coupon, error) {
	c := &Coupon{}
	err := r.db.QueryRow(ctx, `
		SELECT id, code, type, percent_off, amount_off_cents, currency,
			max_redemptions, times_redeemed, valid_from, valid_until, is_active, created_at, updated_at
		FROM coupons WHERE id = $1`, id,
	).Scan(&c.ID, &c.Code, &c.Type, &c.PercentOff, &c.AmountOffCents, &c.Currency,
		&c.MaxRedemptions, &c.TimesRedeemed, &c.ValidFrom, &c.ValidUntil, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*Coupon, error) {
	c := &Coupon{}
	err := r.db.QueryRow(ctx, `
		SELECT id, code, type, percent_off, amount_off_cents, currency,
			max_redemptions, times_redeemed, valid_from, valid_until, is_active, created_at, updated_at
		FROM coupons WHERE code = $1`, code,
	).Scan(&c.ID, &c.Code, &c.Type, &c.PercentOff, &c.AmountOffCents, &c.Currency,
		&c.MaxRedemptions, &c.TimesRedeemed, &c.ValidFrom, &c.ValidUntil, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Coupon, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, code, type, percent_off, amount_off_cents, currency,
			max_redemptions, times_redeemed, valid_from, valid_until, is_active, created_at, updated_at
		FROM coupons ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Coupon
	for rows.Next() {
		var c Coupon
		if err := rows.Scan(&c.ID, &c.Code, &c.Type, &c.PercentOff, &c.AmountOffCents, &c.Currency,
			&c.MaxRedemptions, &c.TimesRedeemed, &c.ValidFrom, &c.ValidUntil, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Coupon, error) {
	_, err := r.db.Exec(ctx, `
		UPDATE coupons SET
			max_redemptions = COALESCE($2, max_redemptions),
			valid_until = COALESCE($3, valid_until),
			is_active = COALESCE($4, is_active),
			updated_at = now()
		WHERE id = $1`,
		id, req.MaxRedemptions, req.ValidUntil, req.IsActive)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM coupons WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) IncrementRedemption(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE coupons SET times_redeemed = times_redeemed + 1, updated_at = now() WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ListForSubscription returns all coupons currently applied to a subscription.
func (r *Repository) ListForSubscription(ctx context.Context, subscriptionID uuid.UUID) ([]Coupon, error) {
	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.code, c.type, c.percent_off, c.amount_off_cents, c.currency,
			c.max_redemptions, c.times_redeemed, c.valid_from, c.valid_until, c.is_active, c.created_at, c.updated_at
		FROM coupons c
		JOIN subscription_coupons sc ON sc.coupon_id = c.id
		WHERE sc.subscription_id = $1
		ORDER BY sc.applied_at`, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Coupon
	for rows.Next() {
		var c Coupon
		if err := rows.Scan(&c.ID, &c.Code, &c.Type, &c.PercentOff, &c.AmountOffCents, &c.Currency,
			&c.MaxRedemptions, &c.TimesRedeemed, &c.ValidFrom, &c.ValidUntil, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *Repository) ApplyToSubscription(ctx context.Context, subscriptionID, couponID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO subscription_coupons (subscription_id, coupon_id, applied_at)
		VALUES ($1, $2, now())`,
		subscriptionID, couponID)
	if err != nil {
		var pgErr interface{ SQLState() string }
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			return ErrAlreadyApplied
		}
		return err
	}
	return nil
}
