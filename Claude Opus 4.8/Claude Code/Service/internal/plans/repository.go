package plans

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("plan not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, p *Plan) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO plans (name, description, currency, billing_interval, pricing_model,
			base_price_cents, seat_price_cents, included_seats, included_units, unit_price_cents,
			usage_metric, trial_days, is_active)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, created_at, updated_at`,
		p.Name, p.Description, p.Currency, p.BillingInterval, p.PricingModel,
		p.BasePriceCents, p.SeatPriceCents, p.IncludedSeats, p.IncludedUnits, p.UnitPriceCents,
		p.UsageMetric, p.TrialDays, p.IsActive,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return err
	}

	for i, t := range p.Tiers {
		var id uuid.UUID
		err = tx.QueryRow(ctx, `
			INSERT INTO pricing_tiers (plan_id, from_unit, to_unit, unit_price_cents, sort_order)
			VALUES ($1,$2,$3,$4,$5) RETURNING id`,
			p.ID, t.FromUnit, t.ToUnit, t.UnitPriceCents, i,
		).Scan(&id)
		if err != nil {
			return err
		}
		p.Tiers[i].ID = id
		p.Tiers[i].SortOrder = i
	}
	return tx.Commit(ctx)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Plan, error) {
	p := &Plan{}
	err := r.db.QueryRow(ctx, `
		SELECT id, name, description, currency, billing_interval, pricing_model,
			base_price_cents, seat_price_cents, included_seats, included_units, unit_price_cents,
			usage_metric, trial_days, is_active, created_at, updated_at
		FROM plans WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Currency, &p.BillingInterval, &p.PricingModel,
		&p.BasePriceCents, &p.SeatPriceCents, &p.IncludedSeats, &p.IncludedUnits, &p.UnitPriceCents,
		&p.UsageMetric, &p.TrialDays, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	tiers, err := r.tiersFor(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Tiers = tiers
	return p, nil
}

func (r *Repository) tiersFor(ctx context.Context, planID uuid.UUID) ([]Tier, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, from_unit, to_unit, unit_price_cents, sort_order
		FROM pricing_tiers WHERE plan_id = $1 ORDER BY sort_order`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tiers []Tier
	for rows.Next() {
		var t Tier
		if err := rows.Scan(&t.ID, &t.FromUnit, &t.ToUnit, &t.UnitPriceCents, &t.SortOrder); err != nil {
			return nil, err
		}
		tiers = append(tiers, t)
	}
	return tiers, rows.Err()
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Plan, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, description, currency, billing_interval, pricing_model,
			base_price_cents, seat_price_cents, included_seats, included_units, unit_price_cents,
			usage_metric, trial_days, is_active, created_at, updated_at
		FROM plans ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Plan
	for rows.Next() {
		var p Plan
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Currency, &p.BillingInterval, &p.PricingModel,
			&p.BasePriceCents, &p.SeatPriceCents, &p.IncludedSeats, &p.IncludedUnits, &p.UnitPriceCents,
			&p.UsageMetric, &p.TrialDays, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Plan, error) {
	_, err := r.db.Exec(ctx, `
		UPDATE plans SET
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			base_price_cents = COALESCE($4, base_price_cents),
			seat_price_cents = COALESCE($5, seat_price_cents),
			included_seats = COALESCE($6, included_seats),
			included_units = COALESCE($7, included_units),
			unit_price_cents = COALESCE($8, unit_price_cents),
			trial_days = COALESCE($9, trial_days),
			is_active = COALESCE($10, is_active),
			updated_at = now()
		WHERE id = $1`,
		id, req.Name, req.Description, req.BasePriceCents, req.SeatPriceCents,
		req.IncludedSeats, req.IncludedUnits, req.UnitPriceCents, req.TrialDays, req.IsActive)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM plans WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
