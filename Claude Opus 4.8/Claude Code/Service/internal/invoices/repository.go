package invoices

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("invoice not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

const invCols = `id, customer_id, subscription_id, status, currency, subtotal_cents, discount_cents,
	tax_cents, total_cents, amount_due_cents, amount_paid_cents, period_start, period_end,
	issued_at, due_at, paid_at, created_at, updated_at`

func scanInvoice(row pgx.Row, inv *Invoice) error {
	return row.Scan(&inv.ID, &inv.CustomerID, &inv.SubscriptionID, &inv.Status, &inv.Currency,
		&inv.SubtotalCents, &inv.DiscountCents, &inv.TaxCents, &inv.TotalCents, &inv.AmountDueCents,
		&inv.AmountPaidCents, &inv.PeriodStart, &inv.PeriodEnd, &inv.IssuedAt, &inv.DueAt, &inv.PaidAt,
		&inv.CreatedAt, &inv.UpdatedAt)
}

// Create inserts an invoice and its lines in a single transaction.
func (r *Repository) Create(ctx context.Context, inv *Invoice) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO invoices (customer_id, subscription_id, status, currency, subtotal_cents,
			discount_cents, tax_cents, total_cents, amount_due_cents, amount_paid_cents,
			period_start, period_end, issued_at, due_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING id, created_at, updated_at`,
		inv.CustomerID, inv.SubscriptionID, inv.Status, inv.Currency, inv.SubtotalCents,
		inv.DiscountCents, inv.TaxCents, inv.TotalCents, inv.AmountDueCents, inv.AmountPaidCents,
		inv.PeriodStart, inv.PeriodEnd, inv.IssuedAt, inv.DueAt,
	).Scan(&inv.ID, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		return err
	}

	for i := range inv.Lines {
		l := &inv.Lines[i]
		meta := l.Metadata
		if len(meta) == 0 {
			meta = []byte("{}")
		}
		err = tx.QueryRow(ctx, `
			INSERT INTO invoice_lines (invoice_id, type, description, quantity, unit_amount_cents, amount_cents, metadata)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			RETURNING id, created_at`,
			inv.ID, l.Type, l.Description, l.Quantity, l.UnitAmountCents, l.AmountCents, meta,
		).Scan(&l.ID, &l.CreatedAt)
		if err != nil {
			return err
		}
		l.InvoiceID = inv.ID
	}
	return tx.Commit(ctx)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	inv := &Invoice{}
	err := scanInvoice(r.db.QueryRow(ctx, `SELECT `+invCols+` FROM invoices WHERE id = $1`, id), inv)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	lines, err := r.linesFor(ctx, id)
	if err != nil {
		return nil, err
	}
	inv.Lines = lines
	return inv, nil
}

func (r *Repository) linesFor(ctx context.Context, invoiceID uuid.UUID) ([]Line, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, invoice_id, type, description, quantity, unit_amount_cents, amount_cents, metadata, created_at
		FROM invoice_lines WHERE invoice_id = $1 ORDER BY created_at`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Line
	for rows.Next() {
		var l Line
		if err := rows.Scan(&l.ID, &l.InvoiceID, &l.Type, &l.Description, &l.Quantity,
			&l.UnitAmountCents, &l.AmountCents, &l.Metadata, &l.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *Repository) List(ctx context.Context, customerID *uuid.UUID, limit, offset int) ([]Invoice, error) {
	var rows pgx.Rows
	var err error
	if customerID != nil {
		rows, err = r.db.Query(ctx, `SELECT `+invCols+` FROM invoices WHERE customer_id = $1
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`, *customerID, limit, offset)
	} else {
		rows, err = r.db.Query(ctx, `SELECT `+invCols+` FROM invoices
			ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Invoice
	for rows.Next() {
		var inv Invoice
		if err := scanInvoice(rows, &inv); err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, rows.Err()
}

// Finalize moves a draft invoice to open and stamps issued_at/due_at.
func (r *Repository) Finalize(ctx context.Context, id uuid.UUID, due time.Time) (*Invoice, error) {
	tag, err := r.db.Exec(ctx, `
		UPDATE invoices SET status = 'open', issued_at = now(), due_at = $2, updated_at = now()
		WHERE id = $1 AND status = 'draft'`, id, due)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) Void(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	tag, err := r.db.Exec(ctx, `
		UPDATE invoices SET status = 'void', updated_at = now()
		WHERE id = $1 AND status IN ('draft','open')`, id)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return r.GetByID(ctx, id)
}

// MarkPaid records a payment against an invoice and flips status to paid when
// fully covered. Returns the updated invoice.
func (r *Repository) MarkPaid(ctx context.Context, id uuid.UUID, amountCents int64) (*Invoice, error) {
	tag, err := r.db.Exec(ctx, `
		UPDATE invoices SET
			amount_paid_cents = amount_paid_cents + $2,
			amount_due_cents = GREATEST(total_cents - (amount_paid_cents + $2), 0),
			status = CASE WHEN (amount_paid_cents + $2) >= total_cents THEN 'paid' ELSE status END,
			paid_at = CASE WHEN (amount_paid_cents + $2) >= total_cents THEN now() ELSE paid_at END,
			updated_at = now()
		WHERE id = $1`, id, amountCents)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return r.GetByID(ctx, id)
}

// MarkUncollectible flags an invoice that can no longer be collected.
func (r *Repository) SetStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE invoices SET status = $2, updated_at = now() WHERE id = $1`, id, status)
	return err
}

// Overdue returns open invoices past their due date.
func (r *Repository) Overdue(ctx context.Context, now time.Time, limit int) ([]Invoice, error) {
	rows, err := r.db.Query(ctx, `SELECT `+invCols+` FROM invoices
		WHERE status = 'open' AND due_at IS NOT NULL AND due_at < $1 LIMIT $2`, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Invoice
	for rows.Next() {
		var inv Invoice
		if err := scanInvoice(rows, &inv); err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, rows.Err()
}
