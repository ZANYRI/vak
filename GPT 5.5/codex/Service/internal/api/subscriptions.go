package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/example/billing-service/internal/auth"
	"github.com/example/billing-service/internal/billing"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type subscriptionInput struct {
	CustomerID        string `json:"customer_id"`
	PlanID            string `json:"plan_id"`
	Quantity          int64  `json:"quantity"`
	CancelAtPeriodEnd *bool  `json:"cancel_at_period_end"`
}

func (a *App) createSubscription(w http.ResponseWriter, r *http.Request) {
	var in subscriptionInput
	if e := decode(r, &in); e != nil || in.CustomerID == "" || in.PlanID == "" || in.Quantity < 1 {
		fail(w, 400, "VALIDATION_ERROR", "customer_id, plan_id and positive quantity are required", nil)
		return
	}
	if !a.canCustomer(r, in.CustomerID) {
		fail(w, 403, "FORBIDDEN", "Customer is not accessible", nil)
		return
	}
	var interval string
	var trial int
	var active bool
	e := a.DB.QueryRow(r.Context(), `SELECT billing_interval,trial_days,is_active FROM plans WHERE id=$1`, in.PlanID).Scan(&interval, &trial, &active)
	if e != nil || !active {
		fail(w, 404, "NOT_FOUND", "Active plan not found", nil)
		return
	}
	now := time.Now().UTC()
	end := addPeriod(now, interval)
	status := "active"
	var ts, te any
	if trial > 0 {
		status = "trialing"
		ts = now
		te = now.AddDate(0, 0, trial)
	}
	sid := uuid.New()
	_, e = a.DB.Exec(r.Context(), `INSERT INTO subscriptions(id,customer_id,plan_id,status,quantity,current_period_start,current_period_end,trial_start,trial_end) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`, sid, in.CustomerID, in.PlanID, status, in.Quantity, now, end, ts, te)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to create subscription", nil)
		return
	}
	a.audit(r, "subscription.created", "subscription", sid.String(), nil)
	_ = a.enqueue(r.Context(), "subscription.created", map[string]string{"subscription_id": sid.String()})
	a.getSubscriptionByID(w, r, sid.String(), 201)
}
func (a *App) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	p := user(r)
	q := `SELECT s.id::text,s.customer_id::text,s.plan_id::text,s.status,s.quantity,s.current_period_start,s.current_period_end,s.trial_start,s.trial_end,s.cancel_at_period_end,s.created_at,s.updated_at FROM subscriptions s`
	var rows pgx.Rows
	var e error
	if p.Role == auth.RoleCustomer {
		rows, e = a.DB.Query(r.Context(), q+` JOIN customers c ON c.id=s.customer_id WHERE c.user_id=$1 ORDER BY s.created_at DESC`, p.ID)
	} else if p.Role == auth.RoleSupport {
		fail(w, 403, "INSUFFICIENT_PERMISSIONS", "Insufficient permissions", nil)
		return
	} else {
		rows, e = a.DB.Query(r.Context(), q+` ORDER BY s.created_at DESC`)
	}
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to list subscriptions", nil)
		return
	}
	defer rows.Close()
	out := []any{}
	for rows.Next() {
		v, e := scanSubscription(rows)
		if e != nil {
			fail(w, 500, "INTERNAL_ERROR", "Unable to read subscriptions", nil)
			return
		}
		out = append(out, v)
	}
	writeJSON(w, 200, map[string]any{"data": out})
}
func (a *App) getSubscription(w http.ResponseWriter, r *http.Request) {
	a.getSubscriptionByID(w, r, id(r), 200)
}
func (a *App) getSubscriptionByID(w http.ResponseWriter, r *http.Request, sid string, status int) {
	if !a.canSubscription(r, sid) {
		fail(w, 403, "FORBIDDEN", "Subscription is not accessible", nil)
		return
	}
	v, e := scanSubscription(a.DB.QueryRow(r.Context(), `SELECT id::text,customer_id::text,plan_id::text,status,quantity,current_period_start,current_period_end,trial_start,trial_end,cancel_at_period_end,created_at,updated_at FROM subscriptions WHERE id=$1`, sid))
	if e != nil {
		if e == pgx.ErrNoRows {
			fail(w, 404, "NOT_FOUND", "Subscription not found", nil)
		} else {
			fail(w, 500, "INTERNAL_ERROR", "Unable to load subscription", nil)
		}
		return
	}
	writeJSON(w, status, v)
}
func scanSubscription(row interface{ Scan(...any) error }) (map[string]any, error) {
	var i, c, p, status string
	var q int64
	var start, end, created, updated time.Time
	var ts, te *time.Time
	var cancel bool
	e := row.Scan(&i, &c, &p, &status, &q, &start, &end, &ts, &te, &cancel, &created, &updated)
	return map[string]any{"id": i, "customer_id": c, "plan_id": p, "status": status, "quantity": q, "current_period_start": start, "current_period_end": end, "trial_start": ts, "trial_end": te, "cancel_at_period_end": cancel, "created_at": created, "updated_at": updated}, e
}
func (a *App) canSubscription(r *http.Request, sid string) bool {
	p := user(r)
	if auth.Allowed(p.Role, auth.RoleAdmin, auth.RoleBillingManager) {
		return true
	}
	if p.Role != auth.RoleCustomer {
		return false
	}
	var ok bool
	e := a.DB.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM subscriptions s JOIN customers c ON c.id=s.customer_id WHERE s.id=$1 AND c.user_id=$2)`, sid, p.ID).Scan(&ok)
	return e == nil && ok
}
func (a *App) updateSubscription(w http.ResponseWriter, r *http.Request) {
	sid := id(r)
	if !a.canSubscription(r, sid) {
		fail(w, 403, "FORBIDDEN", "Subscription is not accessible", nil)
		return
	}
	var in subscriptionInput
	if e := decode(r, &in); e != nil || in.Quantity < 0 {
		fail(w, 400, "VALIDATION_ERROR", "Invalid subscription payload", nil)
		return
	}
	_, e := a.DB.Exec(r.Context(), `UPDATE subscriptions SET quantity=CASE WHEN $2=0 THEN quantity ELSE $2 END,cancel_at_period_end=COALESCE($3,cancel_at_period_end),updated_at=now() WHERE id=$1`, sid, in.Quantity, in.CancelAtPeriodEnd)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to update subscription", nil)
		return
	}
	a.audit(r, "subscription.updated", "subscription", sid, nil)
	a.getSubscriptionByID(w, r, sid, 200)
}
func (a *App) setSubscriptionStatus(w http.ResponseWriter, r *http.Request, status string) {
	sid := id(r)
	if !a.canSubscription(r, sid) {
		fail(w, 403, "FORBIDDEN", "Subscription is not accessible", nil)
		return
	}
	tag, e := a.DB.Exec(r.Context(), `UPDATE subscriptions SET status=$2,updated_at=now() WHERE id=$1 AND status NOT IN ('cancelled','expired')`, sid, status)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to update subscription", nil)
		return
	}
	if tag.RowsAffected() == 0 {
		fail(w, 409, "CONFLICT", "Subscription cannot transition to requested state", nil)
		return
	}
	a.audit(r, "subscription."+status, "subscription", sid, nil)
	a.getSubscriptionByID(w, r, sid, 200)
}
func (a *App) cancelSubscription(w http.ResponseWriter, r *http.Request) {
	a.setSubscriptionStatus(w, r, "cancelled")
}
func (a *App) pauseSubscription(w http.ResponseWriter, r *http.Request) {
	a.setSubscriptionStatus(w, r, "paused")
}
func (a *App) resumeSubscription(w http.ResponseWriter, r *http.Request) {
	a.setSubscriptionStatus(w, r, "active")
}

type changePlanInput struct {
	PlanID string `json:"plan_id"`
}

func (a *App) changePlan(w http.ResponseWriter, r *http.Request) {
	sid := id(r)
	if !a.canSubscription(r, sid) {
		fail(w, 403, "FORBIDDEN", "Subscription is not accessible", nil)
		return
	}
	var in changePlanInput
	if e := decode(r, &in); e != nil || in.PlanID == "" {
		fail(w, 400, "VALIDATION_ERROR", "plan_id is required", nil)
		return
	}
	var customerID, oldPlan, currency string
	var q, oldBase int64
	var start, end time.Time
	e := a.DB.QueryRow(r.Context(), `SELECT s.customer_id::text,s.plan_id::text,s.quantity,s.current_period_start,s.current_period_end,p.base_price_cents,p.currency FROM subscriptions s JOIN plans p ON p.id=s.plan_id WHERE s.id=$1`, sid).Scan(&customerID, &oldPlan, &q, &start, &end, &oldBase, &currency)
	if e != nil {
		fail(w, 404, "NOT_FOUND", "Subscription not found", nil)
		return
	}
	var newBase int64
	var newCurrency string
	e = a.DB.QueryRow(r.Context(), `SELECT base_price_cents,currency FROM plans WHERE id=$1 AND is_active`, in.PlanID).Scan(&newBase, &newCurrency)
	if e != nil {
		fail(w, 404, "NOT_FOUND", "Active target plan not found", nil)
		return
	}
	if newCurrency != currency {
		fail(w, 400, "VALIDATION_ERROR", "Plan currency cannot change during a subscription", nil)
		return
	}
	credit, charge, diff, e := billing.Proration(oldBase, newBase, start.Unix(), end.Unix(), time.Now().UTC().Unix())
	if e != nil {
		fail(w, 400, "VALIDATION_ERROR", e.Error(), nil)
		return
	}
	tx, e := a.DB.Begin(r.Context())
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to change plan", nil)
		return
	}
	defer tx.Rollback(r.Context())
	_, e = tx.Exec(r.Context(), `UPDATE subscriptions SET plan_id=$2,updated_at=now() WHERE id=$1`, sid, in.PlanID)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to change plan", nil)
		return
	}
	var invoiceID string
	if diff != 0 {
		inv := uuid.New()
		due := diff
		if due < 0 {
			due = 0
		}
		_, e = tx.Exec(r.Context(), `INSERT INTO invoices(id,customer_id,subscription_id,status,currency,subtotal_cents,discount_cents,tax_cents,total_cents,amount_due_cents,period_start,period_end,issued_at,due_at) VALUES($1,$2,$3,'open',$4,$5,0,0,$5,$6,$7,$8,now(),now()+interval '14 days')`, inv, customerID, sid, currency, diff, due, start, end)
		if e == nil {
			_, e = tx.Exec(r.Context(), `INSERT INTO invoice_lines(id,invoice_id,type,description,quantity,unit_amount_cents,amount_cents,metadata) VALUES($1,$2,'proration','Unused old plan credit',1,$3,$4,'{}'),($5,$2,'proration','New plan charge',1,$6,$7,'{}')`, uuid.New(), inv, -credit, -credit, uuid.New(), charge, charge)
		}
		invoiceID = inv.String()
	}
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to record proration", nil)
		return
	}
	if e = tx.Commit(r.Context()); e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to change plan", nil)
		return
	}
	a.audit(r, "subscription.plan_changed", "subscription", sid, map[string]any{"old_plan_id": oldPlan, "new_plan_id": in.PlanID, "proration_cents": diff})
	if invoiceID != "" {
		_ = a.enqueue(r.Context(), "invoice.finalize", map[string]string{"invoice_id": invoiceID})
	}
	writeJSON(w, 200, map[string]any{"subscription_id": sid, "old_plan_id": oldPlan, "new_plan_id": in.PlanID, "proration_credit_cents": credit, "proration_charge_cents": charge, "proration_difference_cents": diff, "invoice_id": invoiceID})
}

type couponApplyInput struct {
	Code string `json:"code"`
}

func (a *App) applyCoupon(w http.ResponseWriter, r *http.Request) {
	sid := id(r)
	if !a.canSubscription(r, sid) {
		fail(w, 403, "FORBIDDEN", "Subscription is not accessible", nil)
		return
	}
	var in couponApplyInput
	if e := decode(r, &in); e != nil || strings.TrimSpace(in.Code) == "" {
		fail(w, 400, "VALIDATION_ERROR", "Coupon code is required", nil)
		return
	}
	tx, e := a.DB.Begin(r.Context())
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to apply coupon", nil)
		return
	}
	defer tx.Rollback(r.Context())
	var cid string
	var max *int
	var redeemed int
	var active bool
	var validFrom, validUntil *time.Time
	e = tx.QueryRow(r.Context(), `SELECT id::text,max_redemptions,times_redeemed,is_active,valid_from,valid_until FROM coupons WHERE code=$1 FOR UPDATE`, strings.TrimSpace(in.Code)).Scan(&cid, &max, &redeemed, &active, &validFrom, &validUntil)
	now := time.Now()
	if e != nil || !active || (validFrom != nil && now.Before(*validFrom)) || (validUntil != nil && !now.Before(*validUntil)) || (max != nil && redeemed >= *max) {
		fail(w, 400, "VALIDATION_ERROR", "Coupon is not valid", nil)
		return
	}
	_, e = tx.Exec(r.Context(), `INSERT INTO subscription_coupons(subscription_id,coupon_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, sid, cid)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to apply coupon", nil)
		return
	}
	_, e = tx.Exec(r.Context(), `UPDATE coupons SET times_redeemed=times_redeemed+1,updated_at=now() WHERE id=$1`, cid)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to apply coupon", nil)
		return
	}
	if e = tx.Commit(r.Context()); e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to apply coupon", nil)
		return
	}
	a.audit(r, "coupon.applied", "subscription", sid, map[string]string{"coupon_id": cid})
	writeJSON(w, 200, map[string]string{"subscription_id": sid, "coupon_id": cid, "status": "applied"})
}

type usageInput struct {
	CustomerID     string     `json:"customer_id"`
	SubscriptionID string     `json:"subscription_id"`
	Metric         string     `json:"metric"`
	Quantity       int64      `json:"quantity"`
	RecordedAt     *time.Time `json:"recorded_at"`
}

func (a *App) recordUsage(w http.ResponseWriter, r *http.Request) {
	var in usageInput
	if e := decode(r, &in); e != nil || in.CustomerID == "" || in.SubscriptionID == "" || strings.TrimSpace(in.Metric) == "" || in.Quantity < 0 {
		fail(w, 400, "VALIDATION_ERROR", "Invalid usage payload", nil)
		return
	}
	if !a.canSubscription(r, in.SubscriptionID) {
		fail(w, 403, "FORBIDDEN", "Subscription is not accessible", nil)
		return
	}
	key := r.Header.Get("Idempotency-Key")
	at := time.Now().UTC()
	if in.RecordedAt != nil {
		at = *in.RecordedAt
	}
	eid := uuid.New()
	_, e := a.DB.Exec(r.Context(), `INSERT INTO usage_events(id,customer_id,subscription_id,metric,quantity,idempotency_key,recorded_at) VALUES($1,$2,$3,$4,$5,$6,$7)`, eid, in.CustomerID, in.SubscriptionID, in.Metric, in.Quantity, key, at)
	if e != nil {
		fail(w, 409, "CONFLICT", "Usage event was already recorded", nil)
		return
	}
	_ = a.enqueue(r.Context(), "usage.aggregate", map[string]string{"subscription_id": in.SubscriptionID})
	writeJSON(w, 201, map[string]any{"id": eid, "customer_id": in.CustomerID, "subscription_id": in.SubscriptionID, "metric": in.Metric, "quantity": in.Quantity, "recorded_at": at})
}
func (a *App) listUsage(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("subscription_id")
	if sid == "" || !a.canSubscription(r, sid) {
		fail(w, 403, "FORBIDDEN", "Subscription is not accessible", nil)
		return
	}
	rows, e := a.DB.Query(r.Context(), `SELECT id::text,customer_id::text,subscription_id::text,metric,quantity,recorded_at,created_at FROM usage_events WHERE subscription_id=$1 ORDER BY recorded_at DESC`, sid)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to list usage", nil)
		return
	}
	defer rows.Close()
	out := []any{}
	for rows.Next() {
		var i, c, s, m string
		var q int64
		var at, created time.Time
		if e = rows.Scan(&i, &c, &s, &m, &q, &at, &created); e != nil {
			fail(w, 500, "INTERNAL_ERROR", "Unable to read usage", nil)
			return
		}
		out = append(out, map[string]any{"id": i, "customer_id": c, "subscription_id": s, "metric": m, "quantity": q, "recorded_at": at, "created_at": created})
	}
	writeJSON(w, 200, map[string]any{"data": out})
}
func (a *App) usageSummary(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query().Get("subscription_id")
	if sid == "" || !a.canSubscription(r, sid) {
		fail(w, 403, "FORBIDDEN", "Subscription is not accessible", nil)
		return
	}
	rows, e := a.DB.Query(r.Context(), `SELECT metric,COALESCE(sum(quantity),0) FROM usage_events WHERE subscription_id=$1 GROUP BY metric ORDER BY metric`, sid)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to summarize usage", nil)
		return
	}
	defer rows.Close()
	out := []any{}
	for rows.Next() {
		var m string
		var q int64
		_ = rows.Scan(&m, &q)
		out = append(out, map[string]any{"metric": m, "quantity": q})
	}
	writeJSON(w, 200, map[string]any{"subscription_id": sid, "data": out})
}
func addPeriod(t time.Time, interval string) time.Time {
	if interval == "yearly" {
		return t.AddDate(1, 0, 0)
	}
	return t.AddDate(0, 1, 0)
}

var _ = json.RawMessage{}
