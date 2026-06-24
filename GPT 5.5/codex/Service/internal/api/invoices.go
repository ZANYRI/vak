package api

import (
	crand "crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/example/billing-service/internal/auth"
	"github.com/example/billing-service/internal/billing"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type invoiceGenerateInput struct {
	SubscriptionID string     `json:"subscription_id"`
	PeriodStart    *time.Time `json:"period_start"`
	PeriodEnd      *time.Time `json:"period_end"`
}

func (a *App) generateInvoice(w http.ResponseWriter, r *http.Request) {
	var in invoiceGenerateInput
	if e := decode(r, &in); e != nil || in.SubscriptionID == "" {
		fail(w, 400, "VALIDATION_ERROR", "subscription_id is required", nil)
		return
	}
	if !a.canSubscription(r, in.SubscriptionID) {
		fail(w, 403, "FORBIDDEN", "Subscription is not accessible", nil)
		return
	}
	iid, e := a.generateInvoiceForSubscription(r, in.SubscriptionID, in.PeriodStart, in.PeriodEnd)
	if e != nil {
		fail(w, 400, "VALIDATION_ERROR", e.Error(), nil)
		return
	}
	a.audit(r, "invoice.generated", "invoice", iid, nil)
	a.metric.invoices.Add(1)
	_ = a.enqueue(r.Context(), "email.invoice_created", map[string]string{"invoice_id": iid})
	a.getInvoiceByID(w, r, iid, 201)
}
func (a *App) generateInvoiceForSubscription(r *http.Request, sid string, ps, pe *time.Time) (string, error) {
	ctx := r.Context()
	var customerID, planID, country, currency, model, interval string
	var qty, base int64
	var pricingText string
	var start, end time.Time
	e := a.DB.QueryRow(ctx, `SELECT s.customer_id::text,s.plan_id::text,s.quantity,s.current_period_start,s.current_period_end,c.country,p.currency,p.pricing_model,p.billing_interval,p.base_price_cents,p.pricing::text FROM subscriptions s JOIN customers c ON c.id=s.customer_id JOIN plans p ON p.id=s.plan_id WHERE s.id=$1 AND s.status IN ('active','trialing','past_due')`, sid).Scan(&customerID, &planID, &qty, &start, &end, &country, &currency, &model, &interval, &base, &pricingText)
	if e != nil {
		return "", fmtError("subscription is not billable")
	}
	if ps != nil {
		start = *ps
	}
	if pe != nil {
		end = *pe
	}
	if !end.After(start) {
		return "", fmtError("billing period is invalid")
	}
	var usage int64
	if e = a.DB.QueryRow(ctx, `SELECT COALESCE(sum(quantity),0) FROM usage_events WHERE subscription_id=$1 AND recorded_at >= $2 AND recorded_at < $3`, sid, start, end).Scan(&usage); e != nil {
		return "", e
	}
	var pricing billing.Pricing
	if e = json.Unmarshal([]byte(pricingText), &pricing); e != nil {
		return "", fmtError("plan pricing is invalid")
	}
	subtotal, e := billing.CalculateSubtotal(model, base, qty, usage, pricing)
	if e != nil {
		return "", e
	}
	var percent *int
	var amount *int64
	var couponID, couponType string
	var couponCurrency *string
	e = a.DB.QueryRow(ctx, `SELECT c.id::text,c.type,c.percent_off,c.amount_off_cents,c.currency FROM subscription_coupons sc JOIN coupons c ON c.id=sc.coupon_id WHERE sc.subscription_id=$1 AND c.is_active AND (c.valid_from IS NULL OR c.valid_from<=now()) AND (c.valid_until IS NULL OR c.valid_until>now()) ORDER BY sc.created_at DESC LIMIT 1`, sid).Scan(&couponID, &couponType, &percent, &amount, &couponCurrency)
	if e != nil && e != pgx.ErrNoRows {
		return "", e
	}
	if couponCurrency != nil && *couponCurrency != currency {
		return "", fmtError("coupon currency does not match subscription currency")
	}
	var fixed int64
	if amount != nil {
		fixed = *amount
	}
	var bps int64
	_ = a.DB.QueryRow(ctx, `SELECT tax_rate_basis_points FROM tax_rules WHERE country=$1 AND is_active ORDER BY CASE WHEN region='' THEN 1 ELSE 0 END LIMIT 1`, country).Scan(&bps)
	totals, e := billing.Calculate(subtotal, percent, fixed, bps)
	if e != nil {
		return "", e
	}
	tx, e := a.DB.Begin(ctx)
	if e != nil {
		return "", e
	}
	defer tx.Rollback(ctx)
	iid := uuid.New()
	_, e = tx.Exec(ctx, `INSERT INTO invoices(id,customer_id,subscription_id,status,currency,subtotal_cents,discount_cents,tax_cents,total_cents,amount_due_cents,period_start,period_end,issued_at,due_at) VALUES($1,$2,$3,'open',$4,$5,$6,$7,$8,$8,$9,$10,now(),now()+interval '14 days')`, iid, customerID, sid, currency, totals.SubtotalCents, totals.DiscountCents, totals.TaxCents, totals.TotalCents, start, end)
	if e != nil {
		return "", e
	}
	lines := []struct {
		typ, desc       string
		q, unit, amount int64
	}{{"base", "Subscription base price", 1, base, base}}
	if subtotal-base != 0 {
		lines = append(lines, struct {
			typ, desc       string
			q, unit, amount int64
		}{"usage", "Usage and quantity charges", 1, subtotal - base, subtotal - base})
	}
	if totals.DiscountCents > 0 {
		lines = append(lines, struct {
			typ, desc       string
			q, unit, amount int64
		}{"discount", "Coupon discount", 1, -totals.DiscountCents, -totals.DiscountCents})
	}
	if totals.TaxCents > 0 {
		lines = append(lines, struct {
			typ, desc       string
			q, unit, amount int64
		}{"tax", "Tax", 1, totals.TaxCents, totals.TaxCents})
	}
	for _, line := range lines {
		_, e = tx.Exec(ctx, `INSERT INTO invoice_lines(id,invoice_id,type,description,quantity,unit_amount_cents,amount_cents,metadata) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, uuid.New(), iid, line.typ, line.desc, line.q, line.unit, line.amount, `{}`)
		if e != nil {
			return "", e
		}
	}
	if e = tx.Commit(ctx); e != nil {
		return "", e
	}
	_ = planID
	_ = interval
	_ = couponID
	_ = couponType
	return iid.String(), nil
}
func fmtError(s string) error { return &billingError{s} }

type billingError struct{ s string }

func (e *billingError) Error() string { return e.s }

func (a *App) listInvoices(w http.ResponseWriter, r *http.Request) {
	p := user(r)
	q := `SELECT i.id::text,i.customer_id::text,i.subscription_id::text,i.status,i.currency,i.subtotal_cents,i.discount_cents,i.tax_cents,i.total_cents,i.amount_due_cents,i.amount_paid_cents,i.period_start,i.period_end,i.issued_at,i.due_at,i.paid_at,i.created_at,i.updated_at FROM invoices i`
	var rows pgx.Rows
	var e error
	if p.Role == auth.RoleCustomer {
		rows, e = a.DB.Query(r.Context(), q+` JOIN customers c ON c.id=i.customer_id WHERE c.user_id=$1 ORDER BY i.created_at DESC`, p.ID)
	} else {
		rows, e = a.DB.Query(r.Context(), q+` ORDER BY i.created_at DESC`)
	}
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to list invoices", nil)
		return
	}
	defer rows.Close()
	out := []any{}
	for rows.Next() {
		v, e := scanInvoice(rows)
		if e != nil {
			fail(w, 500, "INTERNAL_ERROR", "Unable to read invoices", nil)
			return
		}
		out = append(out, v)
	}
	writeJSON(w, 200, map[string]any{"data": out})
}
func (a *App) getInvoice(w http.ResponseWriter, r *http.Request) { a.getInvoiceByID(w, r, id(r), 200) }
func (a *App) getInvoiceByID(w http.ResponseWriter, r *http.Request, iid string, status int) {
	if !a.canInvoice(r, iid) {
		fail(w, 403, "FORBIDDEN", "Invoice is not accessible", nil)
		return
	}
	v, e := scanInvoice(a.DB.QueryRow(r.Context(), `SELECT id::text,customer_id::text,subscription_id::text,status,currency,subtotal_cents,discount_cents,tax_cents,total_cents,amount_due_cents,amount_paid_cents,period_start,period_end,issued_at,due_at,paid_at,created_at,updated_at FROM invoices WHERE id=$1`, iid))
	if e != nil {
		if e == pgx.ErrNoRows {
			fail(w, 404, "NOT_FOUND", "Invoice not found", nil)
		} else {
			fail(w, 500, "INTERNAL_ERROR", "Unable to load invoice", nil)
		}
		return
	}
	rows, e := a.DB.Query(r.Context(), `SELECT id::text,type,description,quantity,unit_amount_cents,amount_cents,metadata::text,created_at FROM invoice_lines WHERE invoice_id=$1 ORDER BY created_at`, iid)
	if e == nil {
		defer rows.Close()
		lines := []any{}
		for rows.Next() {
			var id, typ, d, meta string
			var q, u, amount int64
			var created time.Time
			if rows.Scan(&id, &typ, &d, &q, &u, &amount, &meta, &created) == nil {
				lines = append(lines, map[string]any{"id": id, "type": typ, "description": d, "quantity": q, "unit_amount_cents": u, "amount_cents": amount, "metadata": json.RawMessage(meta), "created_at": created})
			}
		}
		v["lines"] = lines
	}
	writeJSON(w, status, v)
}
func scanInvoice(row interface{ Scan(...any) error }) (map[string]any, error) {
	var i, c string
	var s *string
	var st, curr string
	var sub, disc, tax, total, due, paid int64
	var start, end, created, updated time.Time
	var issued, dueAt, paidAt *time.Time
	e := row.Scan(&i, &c, &s, &st, &curr, &sub, &disc, &tax, &total, &due, &paid, &start, &end, &issued, &dueAt, &paidAt, &created, &updated)
	return map[string]any{"id": i, "customer_id": c, "subscription_id": s, "status": st, "currency": curr, "subtotal_cents": sub, "discount_cents": disc, "tax_cents": tax, "total_cents": total, "amount_due_cents": due, "amount_paid_cents": paid, "period_start": start, "period_end": end, "issued_at": issued, "due_at": dueAt, "paid_at": paidAt, "created_at": created, "updated_at": updated}, e
}
func (a *App) canInvoice(r *http.Request, iid string) bool {
	p := user(r)
	if auth.Allowed(p.Role, auth.RoleAdmin, auth.RoleBillingManager, auth.RoleSupport) {
		return true
	}
	if p.Role != auth.RoleCustomer {
		return false
	}
	var ok bool
	e := a.DB.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM invoices i JOIN customers c ON c.id=i.customer_id WHERE i.id=$1 AND c.user_id=$2)`, iid, p.ID).Scan(&ok)
	return e == nil && ok
}
func (a *App) finalizeInvoice(w http.ResponseWriter, r *http.Request) {
	iid := id(r)
	if !a.canInvoice(r, iid) {
		fail(w, 403, "FORBIDDEN", "Invoice is not accessible", nil)
		return
	}
	tag, e := a.DB.Exec(r.Context(), `UPDATE invoices SET status='open',issued_at=COALESCE(issued_at,now()),due_at=COALESCE(due_at,now()+interval '14 days'),updated_at=now() WHERE id=$1 AND status='draft'`, iid)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to finalize invoice", nil)
		return
	}
	if tag.RowsAffected() == 0 {
		fail(w, 409, "CONFLICT", "Only draft invoices can be finalized", nil)
		return
	}
	a.audit(r, "invoice.finalized", "invoice", iid, nil)
	a.getInvoiceByID(w, r, iid, 200)
}
func (a *App) voidInvoice(w http.ResponseWriter, r *http.Request) {
	iid := id(r)
	if !a.canInvoice(r, iid) {
		fail(w, 403, "FORBIDDEN", "Invoice is not accessible", nil)
		return
	}
	tag, e := a.DB.Exec(r.Context(), `UPDATE invoices SET status='void',amount_due_cents=0,updated_at=now() WHERE id=$1 AND status IN ('draft','open')`, iid)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to void invoice", nil)
		return
	}
	if tag.RowsAffected() == 0 {
		fail(w, 409, "CONFLICT", "Invoice cannot be voided", nil)
		return
	}
	a.audit(r, "invoice.voided", "invoice", iid, nil)
	a.getInvoiceByID(w, r, iid, 200)
}

type paymentInput struct {
	InvoiceID  string `json:"invoice_id"`
	CardNumber string `json:"card_number"`
}

func (a *App) simulatePayment(w http.ResponseWriter, r *http.Request) {
	var in paymentInput
	if e := decode(r, &in); e != nil || in.InvoiceID == "" {
		fail(w, 400, "VALIDATION_ERROR", "invoice_id and card_number are required", nil)
		return
	}
	a.payment(w, r, in.InvoiceID, in.CardNumber)
}
func (a *App) payInvoice(w http.ResponseWriter, r *http.Request) {
	var in paymentInput
	if e := decode(r, &in); e != nil {
		fail(w, 400, "VALIDATION_ERROR", "Invalid payment payload", nil)
		return
	}
	a.payment(w, r, id(r), in.CardNumber)
}
func (a *App) payment(w http.ResponseWriter, r *http.Request, iid, card string) {
	if !a.canInvoice(r, iid) {
		fail(w, 403, "FORBIDDEN", "Invoice is not accessible", nil)
		return
	}
	digits := onlyDigits(card)
	if len(digits) < 4 {
		fail(w, 400, "VALIDATION_ERROR", "A card number with at least four digits is required", nil)
		return
	}
	last4 := digits[len(digits)-4:]
	var due int64
	var subID *string
	var status string
	e := a.DB.QueryRow(r.Context(), `SELECT amount_due_cents,subscription_id::text,status FROM invoices WHERE id=$1`, iid).Scan(&due, &subID, &status)
	if e != nil {
		fail(w, 404, "NOT_FOUND", "Invoice not found", nil)
		return
	}
	if status != "open" {
		fail(w, 409, "CONFLICT", "Only open invoices can be paid", nil)
		return
	}
	success := last4 == "0000"
	if last4 != "0000" && last4 != "9999" {
		n, _ := crand.Int(crand.Reader, big.NewInt(2))
		success = n.Int64() == 1
	}
	pid := uuid.New()
	state := "failed"
	if success {
		state = "succeeded"
	}
	tx, e := a.DB.Begin(r.Context())
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to process payment", nil)
		return
	}
	defer tx.Rollback(r.Context())
	_, e = tx.Exec(r.Context(), `INSERT INTO payments(id,invoice_id,status,amount_cents,card_last4,failure_reason) VALUES($1,$2,$3,$4,$5,$6)`, pid, iid, state, due, last4, map[bool]any{true: nil, false: "simulated decline"}[success])
	if e == nil && success {
		_, e = tx.Exec(r.Context(), `UPDATE invoices SET status='paid',amount_paid_cents=amount_due_cents,amount_due_cents=0,paid_at=now(),updated_at=now() WHERE id=$1`, iid)
	}
	if e == nil && !success && subID != nil {
		_, e = tx.Exec(r.Context(), `UPDATE subscriptions SET status='past_due',updated_at=now() WHERE id=$1`, *subID)
	}
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to process payment", nil)
		return
	}
	if e = tx.Commit(r.Context()); e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to process payment", nil)
		return
	}
	if success {
		a.metric.paymentsOK.Add(1)
		a.audit(r, "invoice.paid", "invoice", iid, nil)
	} else {
		a.metric.paymentsFailed.Add(1)
		a.audit(r, "payment.failed", "invoice", iid, nil)
		_ = a.enqueue(r.Context(), "email.payment_failed", map[string]string{"invoice_id": iid})
	}
	writeJSON(w, 201, map[string]any{"id": pid, "invoice_id": iid, "status": state, "amount_cents": due, "card_last4": last4})
}
func (a *App) listPayments(w http.ResponseWriter, r *http.Request) {
	p := user(r)
	q := `SELECT p.id::text,p.invoice_id::text,p.status,p.amount_cents,p.card_last4,p.failure_reason,p.created_at,p.updated_at FROM payments p`
	var rows pgx.Rows
	var e error
	if p.Role == auth.RoleCustomer {
		rows, e = a.DB.Query(r.Context(), q+` JOIN invoices i ON i.id=p.invoice_id JOIN customers c ON c.id=i.customer_id WHERE c.user_id=$1 ORDER BY p.created_at DESC`, p.ID)
	} else {
		rows, e = a.DB.Query(r.Context(), q+` ORDER BY p.created_at DESC`)
	}
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to list payments", nil)
		return
	}
	defer rows.Close()
	out := []any{}
	for rows.Next() {
		v, e := scanPayment(rows)
		if e != nil {
			fail(w, 500, "INTERNAL_ERROR", "Unable to read payments", nil)
			return
		}
		out = append(out, v)
	}
	writeJSON(w, 200, map[string]any{"data": out})
}
func (a *App) getPayment(w http.ResponseWriter, r *http.Request) {
	pid := id(r)
	var iid string
	e := a.DB.QueryRow(r.Context(), `SELECT invoice_id::text FROM payments WHERE id=$1`, pid).Scan(&iid)
	if e != nil {
		fail(w, 404, "NOT_FOUND", "Payment not found", nil)
		return
	}
	if !a.canInvoice(r, iid) {
		fail(w, 403, "FORBIDDEN", "Payment is not accessible", nil)
		return
	}
	v, e := scanPayment(a.DB.QueryRow(r.Context(), `SELECT id::text,invoice_id::text,status,amount_cents,card_last4,failure_reason,created_at,updated_at FROM payments WHERE id=$1`, pid))
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to load payment", nil)
		return
	}
	writeJSON(w, 200, v)
}
func scanPayment(row interface{ Scan(...any) error }) (map[string]any, error) {
	var i, inv, st string
	var amount int64
	var last4 *string
	var reason *string
	var created, updated time.Time
	e := row.Scan(&i, &inv, &st, &amount, &last4, &reason, &created, &updated)
	return map[string]any{"id": i, "invoice_id": inv, "status": st, "amount_cents": amount, "card_last4": last4, "failure_reason": reason, "created_at": created, "updated_at": updated}, e
}
func onlyDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
