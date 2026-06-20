package coupons

import (
	"testing"
	"time"
)

func pct(v int) *int       { return &v }
func cents(v int64) *int64 { return &v }
func str(s string) *string { return &s }

func TestValidateActivePercentage(t *testing.T) {
	s := NewService(nil, nil)
	c := &Coupon{Type: "percentage", PercentOff: pct(20), IsActive: true}
	if err := s.Validate(c, "USD", time.Now()); err != nil {
		t.Fatalf("valid percentage coupon rejected: %v", err)
	}
}

func TestValidateInactive(t *testing.T) {
	s := NewService(nil, nil)
	c := &Coupon{Type: "percentage", PercentOff: pct(20), IsActive: false}
	if err := s.Validate(c, "USD", time.Now()); err == nil {
		t.Fatal("inactive coupon should be rejected")
	}
}

func TestValidateExpired(t *testing.T) {
	s := NewService(nil, nil)
	past := time.Now().Add(-time.Hour)
	c := &Coupon{Type: "percentage", PercentOff: pct(20), IsActive: true, ValidUntil: &past}
	if err := s.Validate(c, "USD", time.Now()); err == nil {
		t.Fatal("expired coupon should be rejected")
	}
}

func TestValidatePercentageOutOfRange(t *testing.T) {
	s := NewService(nil, nil)
	c := &Coupon{Type: "percentage", PercentOff: pct(150), IsActive: true}
	if err := s.Validate(c, "USD", time.Now()); err == nil {
		t.Fatal("percent_off > 100 should be rejected")
	}
}

func TestValidateFixedCurrencyMismatch(t *testing.T) {
	s := NewService(nil, nil)
	c := &Coupon{Type: "fixed_amount", AmountOffCents: cents(500), Currency: str("EUR"), IsActive: true}
	if err := s.Validate(c, "USD", time.Now()); err == nil {
		t.Fatal("fixed coupon with mismatched currency should be rejected")
	}
}

func TestValidateRedemptionLimit(t *testing.T) {
	s := NewService(nil, nil)
	max := 5
	c := &Coupon{Type: "percentage", PercentOff: pct(10), IsActive: true, MaxRedemptions: &max, TimesRedeemed: 5}
	if err := s.Validate(c, "USD", time.Now()); err == nil {
		t.Fatal("coupon at redemption limit should be rejected")
	}
}
