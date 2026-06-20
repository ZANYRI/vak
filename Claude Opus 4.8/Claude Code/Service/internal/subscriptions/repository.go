package subscriptions

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("subscription not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

const cols = `id, customer_id, plan_id, status, quantity, current_period_start, current_period_end,
	trial_start, trial_end, cancel_at_period_end, created_at, updated_at`

func scan(row pgx.Row, s *Subscription) error {
	return row.Scan(&s.ID, &s.CustomerID, &s.PlanID, &s.Status, &s.Quantity,
		&s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.TrialStart, &s.TrialEnd,
		&s.CancelAtPeriodEnd, &s.CreatedAt, &s.UpdatedAt)
}

func (r *Repository) Create(ctx context.Context, s *Subscription) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO subscriptions (customer_id, plan_id, status, quantity,
			current_period_start, current_period_end, trial_start, trial_end, cancel_at_period_end)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING `+cols,
		s.CustomerID, s.PlanID, s.Status, s.Quantity, s.CurrentPeriodStart, s.CurrentPeriodEnd,
		s.TrialStart, s.TrialEnd, s.CancelAtPeriodEnd,
	).Scan(&s.ID, &s.CustomerID, &s.PlanID, &s.Status, &s.Quantity, &s.CurrentPeriodStart,
		&s.CurrentPeriodEnd, &s.TrialStart, &s.TrialEnd, &s.CancelAtPeriodEnd, &s.CreatedAt, &s.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Subscription, error) {
	s := &Subscription{}
	err := scan(r.db.QueryRow(ctx, `SELECT `+cols+` FROM subscriptions WHERE id = $1`, id), s)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return s, err
}

func (r *Repository) List(ctx context.Context, customerID *uuid.UUID, limit, offset int) ([]Subscription, error) {
	var rows pgx.Rows
	var err error
	if customerID != nil {
		rows, err = r.db.Query(ctx, `SELECT `+cols+` FROM subscriptions WHERE customer_id = $1
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`, *customerID, limit, offset)
	} else {
		rows, err = r.db.Query(ctx, `SELECT `+cols+` FROM subscriptions
			ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Subscription
	for rows.Next() {
		var s Subscription
		if err := scan(rows, &s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Subscription, error) {
	_, err := r.db.Exec(ctx, `
		UPDATE subscriptions SET
			quantity = COALESCE($2, quantity),
			cancel_at_period_end = COALESCE($3, cancel_at_period_end),
			updated_at = now()
		WHERE id = $1`, id, req.Quantity, req.CancelAtPeriodEnd)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) SetStatus(ctx context.Context, id uuid.UUID, status string) error {
	tag, err := r.db.Exec(ctx, `UPDATE subscriptions SET status = $2, updated_at = now() WHERE id = $1`, id, status)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) SetCancelAtPeriodEnd(ctx context.Context, id uuid.UUID, v bool) error {
	_, err := r.db.Exec(ctx, `UPDATE subscriptions SET cancel_at_period_end = $2, updated_at = now() WHERE id = $1`, id, v)
	return err
}

func (r *Repository) ChangePlan(ctx context.Context, id, planID uuid.UUID) (*Subscription, error) {
	_, err := r.db.Exec(ctx, `UPDATE subscriptions SET plan_id = $2, updated_at = now() WHERE id = $1`, id, planID)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

// AdvancePeriod moves the billing period forward and reactivates the subscription.
func (r *Repository) AdvancePeriod(ctx context.Context, id uuid.UUID, newStart, newEnd time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE subscriptions SET current_period_start = $2, current_period_end = $3,
			status = CASE WHEN status IN ('trialing','past_due') THEN 'active' ELSE status END,
			updated_at = now()
		WHERE id = $1`, id, newStart, newEnd)
	return err
}

// DueForRenewal returns active subscriptions whose period has ended.
func (r *Repository) DueForRenewal(ctx context.Context, now time.Time, limit int) ([]Subscription, error) {
	rows, err := r.db.Query(ctx, `SELECT `+cols+` FROM subscriptions
		WHERE status IN ('active','past_due') AND current_period_end <= $1 LIMIT $2`, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Subscription
	for rows.Next() {
		var s Subscription
		if err := scan(rows, &s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// ExpiredTrials returns trialing subscriptions whose trial has ended.
func (r *Repository) ExpiredTrials(ctx context.Context, now time.Time, limit int) ([]Subscription, error) {
	rows, err := r.db.Query(ctx, `SELECT `+cols+` FROM subscriptions
		WHERE status = 'trialing' AND trial_end IS NOT NULL AND trial_end <= $1 LIMIT $2`, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Subscription
	for rows.Next() {
		var s Subscription
		if err := scan(rows, &s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
