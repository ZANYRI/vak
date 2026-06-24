package payments

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"billing-service/internal/invoices"
	"billing-service/internal/models"
	"billing-service/internal/subscriptions"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service provides payment simulation.
type Service struct {
	pool          *pgxpool.Pool
	invoices      *invoices.Service
	subscriptions *subscriptions.Service
}

// NewService creates a payment service.
func NewService(pool *pgxpool.Pool, inv *invoices.Service, subs *subscriptions.Service) *Service {
	return &Service{pool: pool, invoices: inv, subscriptions: subs}
}

// SimulateRequest triggers a mock payment.
type SimulateRequest struct {
	InvoiceID uuid.UUID `json:"invoice_id" validate:"required"`
	CardNumber string   `json:"card_number" validate:"required"`
}

// Simulate processes a payment attempt.
func (s *Service) Simulate(ctx context.Context, req SimulateRequest, idempotencyKey string) (*models.Payment, error) {
	if existing, err := s.findByIdempotency(ctx, idempotencyKey); err == nil {
		return existing, nil
	}

	inv, err := s.invoices.GetByID(ctx, req.InvoiceID)
	if err != nil {
		return nil, err
	}
	status, reason := determineOutcome(req.CardNumber)

	last4 := req.CardNumber
	if len(last4) > 4 {
		last4 = last4[len(last4)-4:]
	}

	pay := &models.Payment{
		ID:             uuid.New(),
		InvoiceID:      inv.ID,
		Status:         status,
		AmountCents:    inv.TotalCents,
		Currency:       inv.Currency,
		CardLast4:      last4,
		FailureReason:  reason,
	}
	if idempotencyKey != "" {
		pay.IdempotencyKey = &idempotencyKey
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO payments (id, invoice_id, status, amount_cents, currency, card_last4, failure_reason, idempotency_key)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		pay.ID, pay.InvoiceID, pay.Status, pay.AmountCents, pay.Currency, pay.CardLast4, pay.FailureReason, pay.IdempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("insert payment: %w", err)
	}

	if pay.Status == models.PaymentSucceeded {
		now := time.Now()
		_, err = tx.Exec(ctx,
			`UPDATE invoices SET status='paid', amount_paid_cents=total_cents, amount_due_cents=0, paid_at=$1, updated_at=NOW() WHERE id=$2`,
			now, inv.ID)
		if err != nil {
			return nil, err
		}
		if inv.SubscriptionID != nil {
			_, _ = tx.Exec(ctx, `UPDATE subscriptions SET status='active', updated_at=NOW() WHERE id=$1`, *inv.SubscriptionID)
		}
	} else {
		if inv.SubscriptionID != nil {
			_, _ = tx.Exec(ctx, `UPDATE subscriptions SET status='past_due', updated_at=NOW() WHERE id=$1 AND status != 'cancelled'`, *inv.SubscriptionID)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return pay, nil
}

// List returns payments for an invoice.
func (s *Service) List(ctx context.Context, invoiceID *uuid.UUID, limit, offset int) ([]models.Payment, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	query := `SELECT id, invoice_id, status, amount_cents, currency, card_last4, failure_reason, created_at, updated_at FROM payments`
	countQuery := `SELECT COUNT(*) FROM payments`
	args := []interface{}{}
	where := ""
	if invoiceID != nil {
		args = append(args, *invoiceID)
		where = fmt.Sprintf(" WHERE invoice_id = $%d", len(args))
	}
	query += where + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.Payment
	for rows.Next() {
		var p models.Payment
		err := rows.Scan(&p.ID, &p.InvoiceID, &p.Status, &p.AmountCents, &p.Currency, &p.CardLast4, &p.FailureReason, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, p)
	}
	countArgs := []interface{}{}
	if invoiceID != nil {
		countArgs = append(countArgs, *invoiceID)
	}
	var total int64
	_ = s.pool.QueryRow(ctx, countQuery+where, countArgs...).Scan(&total)
	return list, total, rows.Err()
}

// GetByID returns a payment.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	var p models.Payment
	err := s.pool.QueryRow(ctx,
		`SELECT id, invoice_id, status, amount_cents, currency, card_last4, failure_reason, created_at, updated_at FROM payments WHERE id=$1`, id).Scan(
		&p.ID, &p.InvoiceID, &p.Status, &p.AmountCents, &p.Currency, &p.CardLast4, &p.FailureReason, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, err
	}
	return &p, nil
}

func determineOutcome(cardNumber string) (models.PaymentStatus, string) {
	trimmed := strings.TrimSpace(cardNumber)
	if strings.HasSuffix(trimmed, "0000") {
		return models.PaymentSucceeded, ""
	}
	if strings.HasSuffix(trimmed, "9999") {
		return models.PaymentFailed, "card_declined"
	}
	n, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return models.PaymentSucceeded, ""
	}
	if n.Int64() < 80 {
		return models.PaymentSucceeded, ""
	}
	return models.PaymentFailed, "card_declined"
}

func (s *Service) findByIdempotency(ctx context.Context, key string) (*models.Payment, error) {
	if key == "" {
		return nil, errors.New("no key")
	}
	var p models.Payment
	err := s.pool.QueryRow(ctx,
		`SELECT id, invoice_id, status, amount_cents, currency, card_last4, failure_reason, created_at, updated_at FROM payments WHERE idempotency_key=$1`, key).Scan(
		&p.ID, &p.InvoiceID, &p.Status, &p.AmountCents, &p.Currency, &p.CardLast4, &p.FailureReason, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
