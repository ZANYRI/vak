package invoices

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"billing-service/internal/billing"
	"billing-service/internal/coupons"
	"billing-service/internal/customers"
	"billing-service/internal/models"
	"billing-service/internal/plans"
	"billing-service/internal/subscriptions"
	"billing-service/internal/taxes"
	"billing-service/internal/usage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service provides invoice generation and management.
type Service struct {
	pool          *pgxpool.Pool
	customers     *customers.Service
	subscriptions *subscriptions.Service
	plans         *plans.Service
	usage         *usage.Service
	coupons       *coupons.Service
	taxes         *taxes.Service
}

// NewService creates an invoice service.
func NewService(pool *pgxpool.Pool, c *customers.Service, subs *subscriptions.Service, p *plans.Service, u *usage.Service, cp *coupons.Service, t *taxes.Service) *Service {
	return &Service{pool: pool, customers: c, subscriptions: subs, plans: p, usage: u, coupons: cp, taxes: t}
}

// GenerateRequest triggers invoice generation.
type GenerateRequest struct {
	SubscriptionID uuid.UUID `json:"subscription_id" validate:"required"`
	CustomerID     uuid.UUID `json:"customer_id" validate:"required"`
}

// Generate creates an invoice for a subscription for the current period.
func (s *Service) Generate(ctx context.Context, req GenerateRequest, idempotencyKey string) (*models.Invoice, error) {
	if existing, err := s.findByIdempotency(ctx, idempotencyKey); err == nil {
		return existing, nil
	}
	sub, err := s.subscriptions.GetByID(ctx, req.SubscriptionID)
	if err != nil {
		return nil, err
	}
	if sub.CustomerID != req.CustomerID {
		return nil, fmt.Errorf("subscription does not belong to customer")
	}
	if sub.Status == models.StatusCancelled || sub.Status == models.StatusExpired {
		return nil, fmt.Errorf("subscription is not active")
	}
	customer, err := s.customers.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}
	plan, err := s.plans.GetByID(ctx, sub.PlanID)
	if err != nil {
		return nil, err
	}

	// Period: for simplicity use current period start/end from subscription.
	periodStart := sub.CurrentPeriodStart
	periodEnd := sub.CurrentPeriodEnd
	usageQty, err := s.usageUsage(ctx, req.SubscriptionID, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}

	var coupon *models.Coupon
	applied, err := s.coupons.GetAppliedCoupons(ctx, sub.ID)
	if err == nil && len(applied) > 0 {
		for _, c := range applied {
			if s.couponValid(&c) {
				coupon = &c
				break
			}
		}
	}

	taxBasisPoints := int64(0)
	if taxRule, err := s.taxes.GetForRegion(ctx, customer.Country, customer.Region); err == nil && taxRule != nil {
		taxBasisPoints = int64(taxRule.TaxRateBasisPoints)
	}

	calcInput := billing.CalculationInput{
		Customer:           customer,
		Subscription:       sub,
		Plan:               plan,
		UsageQuantity:      usageQty,
		Coupon:             coupon,
		TaxRateBasisPoints: taxBasisPoints,
	}
	result, err := billing.CalculateInvoice(calcInput, time.Now())
	if err != nil {
		return nil, err
	}

	inv := &models.Invoice{
		ID:              uuid.New(),
		CustomerID:      customer.ID,
		SubscriptionID:  &sub.ID,
		Status:          models.InvoiceDraft,
		Currency:        plan.Currency,
		SubtotalCents:   result.SubtotalCents,
		DiscountCents:   result.DiscountCents,
		TaxCents:        result.TaxCents,
		TotalCents:      result.TotalCents,
		AmountDueCents:  result.TotalCents,
		AmountPaidCents: 0,
		PeriodStart:     &periodStart,
		PeriodEnd:       &periodEnd,
		Lines:           result.Lines,
	}
	if idempotencyKey != "" {
		inv.IdempotencyKey = &idempotencyKey
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		`INSERT INTO invoices (id, customer_id, subscription_id, status, currency, subtotal_cents, discount_cents, tax_cents, total_cents, amount_due_cents, amount_paid_cents, period_start, period_end, idempotency_key)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		inv.ID, inv.CustomerID, inv.SubscriptionID, inv.Status, inv.Currency, inv.SubtotalCents, inv.DiscountCents, inv.TaxCents, inv.TotalCents, inv.AmountDueCents, inv.AmountPaidCents, inv.PeriodStart, inv.PeriodEnd, inv.IdempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("insert invoice: %w", err)
	}
	for i := range inv.Lines {
		line := &inv.Lines[i]
		meta, _ := json.Marshal(line.Metadata)
		_, err = tx.Exec(ctx,
			`INSERT INTO invoice_lines (id, invoice_id, type, description, quantity, unit_amount_cents, amount_cents, metadata)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			line.ID, inv.ID, line.Type, line.Description, line.Quantity, line.UnitAmountCents, line.AmountCents, meta)
		if err != nil {
			return nil, fmt.Errorf("insert invoice line: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, inv.ID)
}

func (s *Service) usageUsage(ctx context.Context, subscriptionID uuid.UUID, from, to time.Time) (int64, error) {
	summary, err := s.usage.Summary(ctx, subscriptionID, "api_requests", from, to)
	if err != nil {
		return 0, err
	}
	return summary.TotalQuantity, nil
}

func (s *Service) couponValid(c *models.Coupon) bool {
	if !c.IsActive {
		return false
	}
	now := time.Now()
	if c.ValidFrom.After(now) || (c.ValidUntil != nil && c.ValidUntil.Before(now)) {
		return false
	}
	if c.MaxRedemptions != nil && c.TimesRedeemed >= *c.MaxRedemptions {
		return false
	}
	return true
}

// List returns invoices for a customer.
func (s *Service) List(ctx context.Context, customerID *uuid.UUID, limit, offset int) ([]models.Invoice, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	query := `SELECT id, customer_id, subscription_id, status, currency, subtotal_cents, discount_cents, tax_cents, total_cents, amount_due_cents, amount_paid_cents, period_start, period_end, issued_at, due_at, paid_at, created_at, updated_at FROM invoices`
	countQuery := `SELECT COUNT(*) FROM invoices`
	args := []interface{}{}
	where := ""
	if customerID != nil {
		args = append(args, *customerID)
		where = fmt.Sprintf(" WHERE customer_id = $%d", len(args))
	}
	query += where + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.Invoice
	for rows.Next() {
		var inv models.Invoice
		err := rows.Scan(&inv.ID, &inv.CustomerID, &inv.SubscriptionID, &inv.Status, &inv.Currency, &inv.SubtotalCents, &inv.DiscountCents, &inv.TaxCents, &inv.TotalCents, &inv.AmountDueCents, &inv.AmountPaidCents, &inv.PeriodStart, &inv.PeriodEnd, &inv.IssuedAt, &inv.DueAt, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, inv)
	}
	countArgs := []interface{}{}
	if customerID != nil {
		countArgs = append(countArgs, *customerID)
	}
	var total int64
	_ = s.pool.QueryRow(ctx, countQuery+where, countArgs...).Scan(&total)
	return list, total, rows.Err()
}

// GetByID returns an invoice with lines.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Invoice, error) {
	var inv models.Invoice
	err := s.pool.QueryRow(ctx,
		`SELECT id, customer_id, subscription_id, status, currency, subtotal_cents, discount_cents, tax_cents, total_cents, amount_due_cents, amount_paid_cents, period_start, period_end, issued_at, due_at, paid_at, created_at, updated_at
		 FROM invoices WHERE id = $1`, id).Scan(
		&inv.ID, &inv.CustomerID, &inv.SubscriptionID, &inv.Status, &inv.Currency, &inv.SubtotalCents, &inv.DiscountCents, &inv.TaxCents, &inv.TotalCents, &inv.AmountDueCents, &inv.AmountPaidCents, &inv.PeriodStart, &inv.PeriodEnd, &inv.IssuedAt, &inv.DueAt, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, err
	}
	lines, err := s.loadLines(ctx, id)
	if err != nil {
		return nil, err
	}
	inv.Lines = lines
	return &inv, nil
}

func (s *Service) loadLines(ctx context.Context, invoiceID uuid.UUID) ([]models.InvoiceLine, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, invoice_id, type, description, quantity, unit_amount_cents, amount_cents, metadata, created_at
		 FROM invoice_lines WHERE invoice_id=$1 ORDER BY created_at`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lines []models.InvoiceLine
	for rows.Next() {
		var line models.InvoiceLine
		var meta []byte
		err := rows.Scan(&line.ID, &line.InvoiceID, &line.Type, &line.Description, &line.Quantity, &line.UnitAmountCents, &line.AmountCents, &meta, &line.CreatedAt)
		if err != nil {
			return nil, err
		}
		if len(meta) > 0 {
			_ = json.Unmarshal(meta, &line.Metadata)
		}
		lines = append(lines, line)
	}
	return lines, rows.Err()
}

// Finalize opens an invoice.
func (s *Service) Finalize(ctx context.Context, id uuid.UUID) (*models.Invoice, error) {
	now := time.Now()
	due := now.AddDate(0, 0, 7)
	_, err := s.pool.Exec(ctx,
		`UPDATE invoices SET status='open', issued_at=$1, due_at=$2, updated_at=NOW() WHERE id=$3 AND status='draft'`,
		now, due, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// Void marks an invoice void.
func (s *Service) Void(ctx context.Context, id uuid.UUID) (*models.Invoice, error) {
	_, err := s.pool.Exec(ctx, `UPDATE invoices SET status='void', updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// Pay marks an invoice as paid.
func (s *Service) Pay(ctx context.Context, id uuid.UUID) (*models.Invoice, error) {
	now := time.Now()
	_, err := s.pool.Exec(ctx,
		`UPDATE invoices SET status='paid', amount_paid_cents=total_cents, amount_due_cents=0, paid_at=$1, updated_at=NOW() WHERE id=$2`,
		now, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *Service) findByIdempotency(ctx context.Context, key string) (*models.Invoice, error) {
	if key == "" {
		return nil, errors.New("no key")
	}
	var inv models.Invoice
	err := s.pool.QueryRow(ctx,
		`SELECT id, customer_id, subscription_id, status, currency, subtotal_cents, discount_cents, tax_cents, total_cents, amount_due_cents, amount_paid_cents, period_start, period_end, issued_at, due_at, paid_at, created_at, updated_at
		 FROM invoices WHERE idempotency_key=$1`, key).Scan(
		&inv.ID, &inv.CustomerID, &inv.SubscriptionID, &inv.Status, &inv.Currency, &inv.SubtotalCents, &inv.DiscountCents, &inv.TaxCents, &inv.TotalCents, &inv.AmountDueCents, &inv.AmountPaidCents, &inv.PeriodStart, &inv.PeriodEnd, &inv.IssuedAt, &inv.DueAt, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		return nil, err
	}
	inv.Lines, _ = s.loadLines(ctx, inv.ID)
	return &inv, nil
}

func hashRequest(body []byte) string {
	h := sha256.Sum256(body)
	return hex.EncodeToString(h[:])
}
