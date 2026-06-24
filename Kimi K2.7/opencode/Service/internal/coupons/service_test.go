package coupons

import (
	"testing"
)

func TestValidateCoupon(t *testing.T) {
	percent := 20
	if err := ValidateCoupon(CreateRequest{Type: "percentage", PercentOff: &percent}); err != nil {
		t.Errorf("valid percentage coupon got error: %v", err)
	}
	if err := ValidateCoupon(CreateRequest{Type: "percentage"}); err == nil {
		t.Error("expected error for missing percent_off")
	}
	amount := int64(500)
	if err := ValidateCoupon(CreateRequest{Type: "fixed_amount", AmountOffCents: &amount}); err != nil {
		t.Errorf("valid fixed coupon got error: %v", err)
	}
}
