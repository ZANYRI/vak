package payments

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("payment not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

const cols = `id, invoice_id, amount_cents, currency, status, card_last4, failure_reason, created_at, updated_at`

func scan(row pgx.Row, p *Payment) error {
	return row.Scan(&p.ID, &p.InvoiceID, &p.AmountCents, &p.Currency, &p.Status,
		&p.CardLast4, &p.FailureReason, &p.CreatedAt, &p.UpdatedAt)
}

func (r *Repository) Create(ctx context.Context, p *Payment) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO payments (invoice_id, amount_cents, currency, status, card_last4, failure_reason)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING `+cols,
		p.InvoiceID, p.AmountCents, p.Currency, p.Status, p.CardLast4, p.FailureReason,
	).Scan(&p.ID, &p.InvoiceID, &p.AmountCents, &p.Currency, &p.Status, &p.CardLast4,
		&p.FailureReason, &p.CreatedAt, &p.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Payment, error) {
	p := &Payment{}
	err := scan(r.db.QueryRow(ctx, `SELECT `+cols+` FROM payments WHERE id = $1`, id), p)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

func (r *Repository) List(ctx context.Context, invoiceID *uuid.UUID, limit, offset int) ([]Payment, error) {
	var rows pgx.Rows
	var err error
	if invoiceID != nil {
		rows, err = r.db.Query(ctx, `SELECT `+cols+` FROM payments WHERE invoice_id = $1
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`, *invoiceID, limit, offset)
	} else {
		rows, err = r.db.Query(ctx, `SELECT `+cols+` FROM payments
			ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Payment
	for rows.Next() {
		var p Payment
		if err := scan(rows, &p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
