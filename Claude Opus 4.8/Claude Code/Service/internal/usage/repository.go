package usage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("usage event not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

const eventCols = `id, customer_id, subscription_id, metric, quantity, idempotency_key, recorded_at, created_at`

func scanEvent(row pgx.Row, e *Event) error {
	return row.Scan(&e.ID, &e.CustomerID, &e.SubscriptionID, &e.Metric, &e.Quantity,
		&e.IdempotencyKey, &e.RecordedAt, &e.CreatedAt)
}

// Insert stores a usage event. Returns (event, false, nil) on insert, or the
// existing event with (event, true, nil) if the idempotency key already exists.
func (r *Repository) Insert(ctx context.Context, e *Event) (*Event, bool, error) {
	err := scanEvent(r.db.QueryRow(ctx, `
		INSERT INTO usage_events (customer_id, subscription_id, metric, quantity, idempotency_key, recorded_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING `+eventCols,
		e.CustomerID, e.SubscriptionID, e.Metric, e.Quantity, e.IdempotencyKey, e.RecordedAt), e)
	if err != nil {
		if isUnique(err) {
			existing, gerr := r.GetByIdempotencyKey(ctx, e.IdempotencyKey)
			if gerr != nil {
				return nil, false, gerr
			}
			return existing, true, nil
		}
		return nil, false, err
	}
	return e, false, nil
}

func (r *Repository) GetByIdempotencyKey(ctx context.Context, key string) (*Event, error) {
	e := &Event{}
	err := scanEvent(r.db.QueryRow(ctx, `SELECT `+eventCols+` FROM usage_events WHERE idempotency_key = $1`, key), e)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return e, err
}

func (r *Repository) List(ctx context.Context, subscriptionID *uuid.UUID, limit, offset int) ([]Event, error) {
	var rows pgx.Rows
	var err error
	if subscriptionID != nil {
		rows, err = r.db.Query(ctx, `SELECT `+eventCols+` FROM usage_events WHERE subscription_id = $1
			ORDER BY recorded_at DESC LIMIT $2 OFFSET $3`, *subscriptionID, limit, offset)
	} else {
		rows, err = r.db.Query(ctx, `SELECT `+eventCols+` FROM usage_events
			ORDER BY recorded_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Event
	for rows.Next() {
		var e Event
		if err := scanEvent(rows, &e); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// SumForPeriod totals quantity for a subscription/metric within [start, end].
func (r *Repository) SumForPeriod(ctx context.Context, subscriptionID uuid.UUID, metric string, start, end time.Time) (int64, error) {
	var total int64
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(quantity), 0) FROM usage_events
		WHERE subscription_id = $1 AND metric = $2 AND recorded_at >= $3 AND recorded_at < $4`,
		subscriptionID, metric, start, end).Scan(&total)
	return total, err
}

// SummaryBySubscription totals quantity grouped by metric for a subscription.
func (r *Repository) SummaryBySubscription(ctx context.Context, subscriptionID uuid.UUID) ([]Summary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT subscription_id, metric, COALESCE(SUM(quantity),0)
		FROM usage_events WHERE subscription_id = $1
		GROUP BY subscription_id, metric ORDER BY metric`, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Summary
	for rows.Next() {
		var s Summary
		if err := rows.Scan(&s.SubscriptionID, &s.Metric, &s.TotalQuantity); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func isUnique(err error) bool {
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
