package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type couponInput struct {
	Code           string     `json:"code"`
	Type           string     `json:"type"`
	PercentOff     *int       `json:"percent_off"`
	AmountOffCents *int64     `json:"amount_off_cents"`
	Currency       *string    `json:"currency"`
	MaxRedemptions *int       `json:"max_redemptions"`
	ValidFrom      *time.Time `json:"valid_from"`
	ValidUntil     *time.Time `json:"valid_until"`
	IsActive       *bool      `json:"is_active"`
}

func validCoupon(in couponInput) bool {
	if strings.TrimSpace(in.Code) == "" || !oneOf(in.Type, "percentage", "fixed_amount") || (in.MaxRedemptions != nil && *in.MaxRedemptions < 1) || (in.ValidFrom != nil && in.ValidUntil != nil && !in.ValidUntil.After(*in.ValidFrom)) {
		return false
	}
	if in.Type == "percentage" {
		return in.PercentOff != nil && *in.PercentOff >= 1 && *in.PercentOff <= 100 && in.AmountOffCents == nil
	}
	return in.AmountOffCents != nil && *in.AmountOffCents > 0 && in.PercentOff == nil && (in.Currency == nil || currency(*in.Currency))
}
func (a *App) createCoupon(w http.ResponseWriter, r *http.Request) {
	var in couponInput
	if e := decode(r, &in); e != nil || !validCoupon(in) {
		fail(w, 400, "VALIDATION_ERROR", "Invalid coupon payload", nil)
		return
	}
	active := true
	if in.IsActive != nil {
		active = *in.IsActive
	}
	cid := uuid.New()
	var curr any
	if in.Currency != nil {
		curr = strings.ToUpper(*in.Currency)
	}
	_, e := a.DB.Exec(r.Context(), `INSERT INTO coupons(id,code,type,percent_off,amount_off_cents,currency,max_redemptions,valid_from,valid_until,is_active) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, cid, strings.TrimSpace(in.Code), in.Type, in.PercentOff, in.AmountOffCents, curr, in.MaxRedemptions, in.ValidFrom, in.ValidUntil, active)
	if e != nil {
		fail(w, 409, "CONFLICT", "Coupon code already exists", nil)
		return
	}
	a.audit(r, "coupon.created", "coupon", cid.String(), nil)
	a.getCouponByID(w, r, cid.String(), 201)
}
func (a *App) listCoupons(w http.ResponseWriter, r *http.Request) {
	rows, e := a.DB.Query(r.Context(), `SELECT id::text,code::text,type,percent_off,amount_off_cents,currency,max_redemptions,times_redeemed,valid_from,valid_until,is_active,created_at,updated_at FROM coupons ORDER BY created_at DESC`)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to list coupons", nil)
		return
	}
	defer rows.Close()
	out := []any{}
	for rows.Next() {
		v, e := scanCoupon(rows)
		if e != nil {
			fail(w, 500, "INTERNAL_ERROR", "Unable to read coupons", nil)
			return
		}
		out = append(out, v)
	}
	writeJSON(w, 200, map[string]any{"data": out})
}
func (a *App) getCoupon(w http.ResponseWriter, r *http.Request) { a.getCouponByID(w, r, id(r), 200) }
func (a *App) getCouponByID(w http.ResponseWriter, r *http.Request, cid string, status int) {
	v, e := scanCoupon(a.DB.QueryRow(r.Context(), `SELECT id::text,code::text,type,percent_off,amount_off_cents,currency,max_redemptions,times_redeemed,valid_from,valid_until,is_active,created_at,updated_at FROM coupons WHERE id=$1`, cid))
	if e != nil {
		if e == pgx.ErrNoRows {
			fail(w, 404, "NOT_FOUND", "Coupon not found", nil)
		} else {
			fail(w, 500, "INTERNAL_ERROR", "Unable to load coupon", nil)
		}
		return
	}
	writeJSON(w, status, v)
}
func scanCoupon(row interface{ Scan(...any) error }) (map[string]any, error) {
	var i, code, typ string
	var pct *int
	var amount *int64
	var curr *string
	var max *int
	var red int
	var vf, vu *time.Time
	var active bool
	var created, updated time.Time
	e := row.Scan(&i, &code, &typ, &pct, &amount, &curr, &max, &red, &vf, &vu, &active, &created, &updated)
	return map[string]any{"id": i, "code": code, "type": typ, "percent_off": pct, "amount_off_cents": amount, "currency": curr, "max_redemptions": max, "times_redeemed": red, "valid_from": vf, "valid_until": vu, "is_active": active, "created_at": created, "updated_at": updated}, e
}
func (a *App) updateCoupon(w http.ResponseWriter, r *http.Request) {
	var in couponInput
	if e := decode(r, &in); e != nil {
		fail(w, 400, "VALIDATION_ERROR", "Invalid coupon payload", nil)
		return
	}
	if in.PercentOff != nil && (*in.PercentOff < 1 || *in.PercentOff > 100) || in.AmountOffCents != nil && *in.AmountOffCents < 1 || in.Currency != nil && !currency(*in.Currency) || in.MaxRedemptions != nil && *in.MaxRedemptions < 1 {
		fail(w, 400, "VALIDATION_ERROR", "Invalid coupon field", nil)
		return
	}
	var curr any
	if in.Currency != nil {
		curr = strings.ToUpper(*in.Currency)
	}
	_, e := a.DB.Exec(r.Context(), `UPDATE coupons SET code=COALESCE(NULLIF($2,''),code),percent_off=COALESCE($3,percent_off),amount_off_cents=COALESCE($4,amount_off_cents),currency=COALESCE($5,currency),max_redemptions=COALESCE($6,max_redemptions),valid_from=COALESCE($7,valid_from),valid_until=COALESCE($8,valid_until),is_active=COALESCE($9,is_active),updated_at=now() WHERE id=$1`, id(r), in.Code, in.PercentOff, in.AmountOffCents, curr, in.MaxRedemptions, in.ValidFrom, in.ValidUntil, in.IsActive)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to update coupon", nil)
		return
	}
	a.audit(r, "coupon.updated", "coupon", id(r), nil)
	a.getCouponByID(w, r, id(r), 200)
}
func (a *App) deleteCoupon(w http.ResponseWriter, r *http.Request) {
	tag, e := a.DB.Exec(r.Context(), `DELETE FROM coupons WHERE id=$1`, id(r))
	if e != nil {
		fail(w, 409, "CONFLICT", "Coupon is in use and cannot be deleted", nil)
		return
	}
	if tag.RowsAffected() == 0 {
		fail(w, 404, "NOT_FOUND", "Coupon not found", nil)
		return
	}
	a.audit(r, "coupon.deleted", "coupon", id(r), nil)
	w.WriteHeader(204)
}

type taxInput struct {
	Country            string `json:"country"`
	Region             string `json:"region"`
	TaxName            string `json:"tax_name"`
	TaxRateBasisPoints int64  `json:"tax_rate_basis_points"`
	IsActive           *bool  `json:"is_active"`
}

func (a *App) createTaxRule(w http.ResponseWriter, r *http.Request) {
	var in taxInput
	if e := decode(r, &in); e != nil || strings.TrimSpace(in.Country) == "" || strings.TrimSpace(in.TaxName) == "" || in.TaxRateBasisPoints < 0 || in.TaxRateBasisPoints > 100000 {
		fail(w, 400, "VALIDATION_ERROR", "Invalid tax rule payload", nil)
		return
	}
	active := true
	if in.IsActive != nil {
		active = *in.IsActive
	}
	tid := uuid.New()
	_, e := a.DB.Exec(r.Context(), `INSERT INTO tax_rules(id,country,region,tax_name,tax_rate_basis_points,is_active) VALUES($1,$2,$3,$4,$5,$6)`, tid, strings.ToUpper(in.Country), in.Region, in.TaxName, in.TaxRateBasisPoints, active)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to create tax rule", nil)
		return
	}
	a.audit(r, "tax_rule.created", "tax_rule", tid.String(), nil)
	a.getTaxRule(w, r, tid.String(), 201)
}
func (a *App) listTaxRules(w http.ResponseWriter, r *http.Request) {
	rows, e := a.DB.Query(r.Context(), `SELECT id::text,country,region,tax_name,tax_rate_basis_points,is_active,created_at,updated_at FROM tax_rules ORDER BY country,region`)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to list tax rules", nil)
		return
	}
	defer rows.Close()
	out := []any{}
	for rows.Next() {
		v, e := scanTax(rows)
		if e != nil {
			fail(w, 500, "INTERNAL_ERROR", "Unable to read tax rules", nil)
			return
		}
		out = append(out, v)
	}
	writeJSON(w, 200, map[string]any{"data": out})
}
func (a *App) getTaxRule(w http.ResponseWriter, r *http.Request, tid string, status int) {
	v, e := scanTax(a.DB.QueryRow(r.Context(), `SELECT id::text,country,region,tax_name,tax_rate_basis_points,is_active,created_at,updated_at FROM tax_rules WHERE id=$1`, tid))
	if e != nil {
		if e == pgx.ErrNoRows {
			fail(w, 404, "NOT_FOUND", "Tax rule not found", nil)
		} else {
			fail(w, 500, "INTERNAL_ERROR", "Unable to load tax rule", nil)
		}
		return
	}
	writeJSON(w, status, v)
}
func scanTax(row interface{ Scan(...any) error }) (map[string]any, error) {
	var i, c, r, n string
	var b int64
	var active bool
	var created, updated time.Time
	e := row.Scan(&i, &c, &r, &n, &b, &active, &created, &updated)
	return map[string]any{"id": i, "country": c, "region": r, "tax_name": n, "tax_rate_basis_points": b, "is_active": active, "created_at": created, "updated_at": updated}, e
}
func (a *App) updateTaxRule(w http.ResponseWriter, r *http.Request) {
	var in taxInput
	if e := decode(r, &in); e != nil || in.TaxRateBasisPoints < 0 || in.TaxRateBasisPoints > 100000 {
		fail(w, 400, "VALIDATION_ERROR", "Invalid tax rule payload", nil)
		return
	}
	_, e := a.DB.Exec(r.Context(), `UPDATE tax_rules SET country=COALESCE(NULLIF($2,''),country),region=COALESCE(NULLIF($3,''),region),tax_name=COALESCE(NULLIF($4,''),tax_name),tax_rate_basis_points=CASE WHEN $5=0 THEN tax_rate_basis_points ELSE $5 END,is_active=COALESCE($6,is_active),updated_at=now() WHERE id=$1`, id(r), strings.ToUpper(in.Country), in.Region, in.TaxName, in.TaxRateBasisPoints, in.IsActive)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to update tax rule", nil)
		return
	}
	a.audit(r, "tax_rule.updated", "tax_rule", id(r), nil)
	a.getTaxRule(w, r, id(r), 200)
}
func (a *App) deleteTaxRule(w http.ResponseWriter, r *http.Request) {
	tag, e := a.DB.Exec(r.Context(), `DELETE FROM tax_rules WHERE id=$1`, id(r))
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to delete tax rule", nil)
		return
	}
	if tag.RowsAffected() == 0 {
		fail(w, 404, "NOT_FOUND", "Tax rule not found", nil)
		return
	}
	a.audit(r, "tax_rule.deleted", "tax_rule", id(r), nil)
	w.WriteHeader(204)
}
