package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/example/billing-service/internal/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type customerInput struct {
	UserID         *string         `json:"user_id"`
	Email          string          `json:"email"`
	Name           string          `json:"name"`
	CompanyName    string          `json:"company_name"`
	BillingAddress json.RawMessage `json:"billing_address"`
	Country        string          `json:"country"`
	TaxID          string          `json:"tax_id"`
	Currency       string          `json:"currency"`
}

func (a *App) createCustomer(w http.ResponseWriter, r *http.Request) {
	var in customerInput
	if e := decode(r, &in); e != nil || !validEmail(in.Email) || strings.TrimSpace(in.Name) == "" || !currency(in.Currency) {
		fail(w, 400, "VALIDATION_ERROR", "Invalid customer payload", nil)
		return
	}
	if len(in.BillingAddress) == 0 {
		in.BillingAddress = []byte(`{}`)
	}
	var uid any
	if in.UserID != nil {
		uid = *in.UserID
	}
	cid := uuid.New()
	_, e := a.DB.Exec(r.Context(), `INSERT INTO customers(id,user_id,email,name,company_name,billing_address,country,tax_id,currency) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`, cid, uid, strings.TrimSpace(in.Email), strings.TrimSpace(in.Name), in.CompanyName, in.BillingAddress, in.Country, in.TaxID, strings.ToUpper(in.Currency))
	if e != nil {
		fail(w, 409, "CONFLICT", "Customer email or user assignment already exists", nil)
		return
	}
	a.audit(r, "customer.created", "customer", cid.String(), nil)
	a.getCustomerByID(w, r, cid.String(), 201)
}
func (a *App) listCustomers(w http.ResponseWriter, r *http.Request) {
	p := user(r)
	q := `SELECT id::text,user_id::text,email::text,name,company_name,billing_address,country,tax_id,currency,created_at,updated_at FROM customers`
	var rows pgx.Rows
	var e error
	if p.Role == auth.RoleCustomer {
		rows, e = a.DB.Query(r.Context(), q+` WHERE user_id=$1 ORDER BY created_at DESC`, p.ID)
	} else {
		rows, e = a.DB.Query(r.Context(), q+` ORDER BY created_at DESC`)
	}
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to list customers", nil)
		return
	}
	defer rows.Close()
	out := []any{}
	for rows.Next() {
		v, e := scanCustomer(rows)
		if e != nil {
			fail(w, 500, "INTERNAL_ERROR", "Unable to read customers", nil)
			return
		}
		out = append(out, v)
	}
	writeJSON(w, 200, map[string]any{"data": out})
}
func (a *App) getCustomer(w http.ResponseWriter, r *http.Request) {
	a.getCustomerByID(w, r, id(r), 200)
}
func (a *App) getCustomerByID(w http.ResponseWriter, r *http.Request, cid string, status int) {
	if !a.canCustomer(r, cid) {
		fail(w, 403, "FORBIDDEN", "Customer is not accessible", nil)
		return
	}
	row := a.DB.QueryRow(r.Context(), `SELECT id::text,user_id::text,email::text,name,company_name,billing_address,country,tax_id,currency,created_at,updated_at FROM customers WHERE id=$1`, cid)
	v, e := scanCustomer(row)
	if e != nil {
		if e == pgx.ErrNoRows {
			fail(w, 404, "NOT_FOUND", "Customer not found", nil)
		} else {
			fail(w, 500, "INTERNAL_ERROR", "Unable to load customer", nil)
		}
		return
	}
	writeJSON(w, status, v)
}
func scanCustomer(row interface{ Scan(...any) error }) (map[string]any, error) {
	var i, email, name, company, address, country, tax, curr string
	var uid *string
	var created, updated time.Time
	e := row.Scan(&i, &uid, &email, &name, &company, &address, &country, &tax, &curr, &created, &updated)
	return map[string]any{"id": i, "user_id": uid, "email": email, "name": name, "company_name": company, "billing_address": json.RawMessage(address), "country": country, "tax_id": tax, "currency": curr, "created_at": created, "updated_at": updated}, e
}
func (a *App) updateCustomer(w http.ResponseWriter, r *http.Request) {
	cid := id(r)
	if !a.canCustomer(r, cid) {
		fail(w, 403, "FORBIDDEN", "Customer is not accessible", nil)
		return
	}
	var in customerInput
	if e := decode(r, &in); e != nil {
		fail(w, 400, "VALIDATION_ERROR", "Invalid customer payload", nil)
		return
	}
	if in.Email != "" && !validEmail(in.Email) {
		fail(w, 400, "VALIDATION_ERROR", "Invalid email", nil)
		return
	}
	if in.Currency != "" && !currency(in.Currency) {
		fail(w, 400, "VALIDATION_ERROR", "Invalid currency", nil)
		return
	}
	addr := in.BillingAddress
	if len(addr) == 0 {
		addr = []byte(`{}`)
	}
	_, e := a.DB.Exec(r.Context(), `UPDATE customers SET email=COALESCE(NULLIF($2,''),email),name=COALESCE(NULLIF($3,''),name),company_name=COALESCE(NULLIF($4,''),company_name),billing_address=CASE WHEN $5::jsonb='{}'::jsonb THEN billing_address ELSE $5::jsonb END,country=COALESCE(NULLIF($6,''),country),tax_id=COALESCE(NULLIF($7,''),tax_id),currency=COALESCE(NULLIF($8,''),currency),updated_at=now() WHERE id=$1`, cid, in.Email, in.Name, in.CompanyName, addr, in.Country, in.TaxID, strings.ToUpper(in.Currency))
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to update customer", nil)
		return
	}
	a.audit(r, "customer.updated", "customer", cid, nil)
	a.getCustomerByID(w, r, cid, 200)
}
func (a *App) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	tag, e := a.DB.Exec(r.Context(), `DELETE FROM customers WHERE id=$1`, id(r))
	if e != nil {
		fail(w, 409, "CONFLICT", "Customer has related records", nil)
		return
	}
	if tag.RowsAffected() == 0 {
		fail(w, 404, "NOT_FOUND", "Customer not found", nil)
		return
	}
	a.audit(r, "customer.deleted", "customer", id(r), nil)
	w.WriteHeader(204)
}
func (a *App) canCustomer(r *http.Request, cid string) bool {
	p := user(r)
	if auth.Allowed(p.Role, auth.RoleAdmin, auth.RoleBillingManager, auth.RoleSupport) {
		return true
	}
	var yes bool
	e := a.DB.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM customers WHERE id=$1 AND user_id=$2)`, cid, p.ID).Scan(&yes)
	return e == nil && yes
}

type planInput struct {
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Currency        string          `json:"currency"`
	BillingInterval string          `json:"billing_interval"`
	PricingModel    string          `json:"pricing_model"`
	BasePriceCents  int64           `json:"base_price_cents"`
	TrialDays       int             `json:"trial_days"`
	Pricing         json.RawMessage `json:"pricing"`
	IsActive        *bool           `json:"is_active"`
}

func (a *App) createPlan(w http.ResponseWriter, r *http.Request) {
	var in planInput
	if e := decode(r, &in); e != nil || strings.TrimSpace(in.Name) == "" || !currency(in.Currency) || !oneOf(in.BillingInterval, "monthly", "yearly") || !oneOf(in.PricingModel, "flat", "per_seat", "usage_based", "tiered", "hybrid") || in.BasePriceCents < 0 || in.TrialDays < 0 {
		fail(w, 400, "VALIDATION_ERROR", "Invalid plan payload", nil)
		return
	}
	if len(in.Pricing) == 0 {
		in.Pricing = []byte(`{}`)
	}
	active := true
	if in.IsActive != nil {
		active = *in.IsActive
	}
	pid := uuid.New()
	_, e := a.DB.Exec(r.Context(), `INSERT INTO plans(id,name,description,currency,billing_interval,pricing_model,base_price_cents,trial_days,pricing,is_active) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, pid, in.Name, in.Description, strings.ToUpper(in.Currency), in.BillingInterval, in.PricingModel, in.BasePriceCents, in.TrialDays, in.Pricing, active)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to create plan", nil)
		return
	}
	a.audit(r, "plan.created", "plan", pid.String(), nil)
	a.getPlanByID(w, r, pid.String(), 201)
}
func (a *App) listPlans(w http.ResponseWriter, r *http.Request) {
	rows, e := a.DB.Query(r.Context(), `SELECT id::text,name,description,currency,billing_interval,pricing_model,base_price_cents,trial_days,pricing,is_active,created_at,updated_at FROM plans ORDER BY created_at DESC`)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to list plans", nil)
		return
	}
	defer rows.Close()
	out := []any{}
	for rows.Next() {
		v, e := scanPlan(rows)
		if e != nil {
			fail(w, 500, "INTERNAL_ERROR", "Unable to read plans", nil)
			return
		}
		out = append(out, v)
	}
	writeJSON(w, 200, map[string]any{"data": out})
}
func (a *App) getPlan(w http.ResponseWriter, r *http.Request) { a.getPlanByID(w, r, id(r), 200) }
func (a *App) getPlanByID(w http.ResponseWriter, r *http.Request, pid string, status int) {
	v, e := scanPlan(a.DB.QueryRow(r.Context(), `SELECT id::text,name,description,currency,billing_interval,pricing_model,base_price_cents,trial_days,pricing,is_active,created_at,updated_at FROM plans WHERE id=$1`, pid))
	if e != nil {
		if e == pgx.ErrNoRows {
			fail(w, 404, "NOT_FOUND", "Plan not found", nil)
		} else {
			fail(w, 500, "INTERNAL_ERROR", "Unable to load plan", nil)
		}
		return
	}
	writeJSON(w, status, v)
}
func scanPlan(row interface{ Scan(...any) error }) (map[string]any, error) {
	var i, n, d, c, bi, pm, p string
	var base int64
	var trial int
	var active bool
	var created, updated time.Time
	e := row.Scan(&i, &n, &d, &c, &bi, &pm, &base, &trial, &p, &active, &created, &updated)
	return map[string]any{"id": i, "name": n, "description": d, "currency": c, "billing_interval": bi, "pricing_model": pm, "base_price_cents": base, "trial_days": trial, "pricing": json.RawMessage(p), "is_active": active, "created_at": created, "updated_at": updated}, e
}
func (a *App) updatePlan(w http.ResponseWriter, r *http.Request) {
	var in planInput
	if e := decode(r, &in); e != nil || in.BasePriceCents < 0 || in.TrialDays < 0 {
		fail(w, 400, "VALIDATION_ERROR", "Invalid plan payload", nil)
		return
	}
	if in.Currency != "" && !currency(in.Currency) || in.BillingInterval != "" && !oneOf(in.BillingInterval, "monthly", "yearly") || in.PricingModel != "" && !oneOf(in.PricingModel, "flat", "per_seat", "usage_based", "tiered", "hybrid") {
		fail(w, 400, "VALIDATION_ERROR", "Invalid plan field", nil)
		return
	}
	pricing := in.Pricing
	if len(pricing) == 0 {
		pricing = []byte(`{}`)
	}
	_, e := a.DB.Exec(r.Context(), `UPDATE plans SET name=COALESCE(NULLIF($2,''),name),description=COALESCE(NULLIF($3,''),description),currency=COALESCE(NULLIF($4,''),currency),billing_interval=COALESCE(NULLIF($5,''),billing_interval),pricing_model=COALESCE(NULLIF($6,''),pricing_model),base_price_cents=CASE WHEN $7=0 THEN base_price_cents ELSE $7 END,trial_days=CASE WHEN $8=0 THEN trial_days ELSE $8 END,pricing=CASE WHEN $9::jsonb='{}'::jsonb THEN pricing ELSE $9::jsonb END,is_active=COALESCE($10,is_active),updated_at=now() WHERE id=$1`, id(r), in.Name, in.Description, strings.ToUpper(in.Currency), in.BillingInterval, in.PricingModel, in.BasePriceCents, in.TrialDays, pricing, in.IsActive)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to update plan", nil)
		return
	}
	a.audit(r, "plan.updated", "plan", id(r), nil)
	a.getPlanByID(w, r, id(r), 200)
}
func (a *App) deletePlan(w http.ResponseWriter, r *http.Request) {
	tag, e := a.DB.Exec(r.Context(), `DELETE FROM plans WHERE id=$1`, id(r))
	if e != nil {
		fail(w, 409, "CONFLICT", "Plan is in use and cannot be deleted", nil)
		return
	}
	if tag.RowsAffected() == 0 {
		fail(w, 404, "NOT_FOUND", "Plan not found", nil)
		return
	}
	a.audit(r, "plan.deleted", "plan", id(r), nil)
	w.WriteHeader(204)
}
func currency(v string) bool { v = strings.ToUpper(v); return len(v) == 3 && v >= "AAA" && v <= "ZZZ" }
func oneOf(v string, choices ...string) bool {
	for _, c := range choices {
		if v == c {
			return true
		}
	}
	return false
}
