package api

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/example/billing-service/internal/auth"
)

func TestCouponValidation(t *testing.T) {
	pct := 20
	if !validCoupon(couponInput{Code: "SAVE20", Type: "percentage", PercentOff: &pct}) {
		t.Fatal("valid percentage coupon rejected")
	}
	invalid := 101
	if validCoupon(couponInput{Code: "NOPE", Type: "percentage", PercentOff: &invalid}) {
		t.Fatal("invalid percentage coupon accepted")
	}
}

func TestAuthorizationAndPeriod(t *testing.T) {
	if !auth.Allowed(auth.RoleAdmin, auth.RoleBillingManager) || auth.Allowed(auth.RoleSupport, auth.RoleBillingManager) {
		t.Fatal("unexpected role decision")
	}
	if got := addPeriod(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), "yearly"); got.Year() != 2027 {
		t.Fatalf("got %s", got)
	}
}

func TestOnlyDigits(t *testing.T) {
	if got := onlyDigits("4242-4242 4242 0000"); got != "4242424242420000" {
		t.Fatalf("got %q", got)
	}
}

func TestDecodeJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"email":"admin@example.com","password":"admin-password-change-me"}`))
	var in loginInput
	if err := decode(req, &in); err != nil || in.Email != "admin@example.com" {
		t.Fatalf("decode: %#v, %v", in, err)
	}
}
