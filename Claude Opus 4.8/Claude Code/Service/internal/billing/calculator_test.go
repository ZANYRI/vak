package billing

import "testing"

func TestFlatPlan(t *testing.T) {
	res := Calculate(Input{
		Plan:     Plan{PricingModel: ModelFlat, BasePriceCents: 1900},
		Currency: "USD",
	})
	if res.SubtotalCents != 1900 || res.TotalCents != 1900 {
		t.Fatalf("flat: got subtotal=%d total=%d want 1900/1900", res.SubtotalCents, res.TotalCents)
	}
}

func TestPerSeatPlan(t *testing.T) {
	res := Calculate(Input{
		Plan: Plan{PricingModel: ModelPerSeat, BasePriceCents: 1200,
			SeatPriceCents: 700, IncludedSeats: 3},
		Seats:    5,
		Currency: "USD",
	})
	// base 1200 + 2 extra seats * 700 = 2600
	if res.SubtotalCents != 2600 {
		t.Fatalf("per_seat: got %d want 2600", res.SubtotalCents)
	}
}

func TestUsagePlan(t *testing.T) {
	res := Calculate(Input{
		Plan: Plan{PricingModel: ModelUsageBased, BasePriceCents: 2900,
			IncludedUnits: 10000, UnitPriceCents: 2},
		UsageUnits: 15000,
		Currency:   "USD",
	})
	// base 2900 + 5000 units * 2 = 12900
	if res.SubtotalCents != 12900 {
		t.Fatalf("usage: got %d want 12900", res.SubtotalCents)
	}
}

func TestTieredPlan(t *testing.T) {
	to1 := int64(1000)
	to2 := int64(10000)
	res := Calculate(Input{
		Plan: Plan{PricingModel: ModelTiered, Tiers: []Tier{
			{From: 0, To: &to1, UnitPriceCents: 5},
			{From: 1001, To: &to2, UnitPriceCents: 3},
			{From: 10001, To: nil, UnitPriceCents: 1},
		}},
		UsageUnits: 1500,
		Currency:   "USD",
	})
	// tier1: units 0..1000 => 1001 * 5 = 5005
	// tier2: units 1001..1500 => 500 * 3 = 1500
	want := int64(5005 + 1500)
	if res.SubtotalCents != want {
		t.Fatalf("tiered: got %d want %d", res.SubtotalCents, want)
	}
}

// TestInvoiceExampleFromSpec reproduces the documented invoice example:
// subtotal 7900, fixed coupon 1000, 20% tax => tax 1380, total 8280.
func TestInvoiceExampleFromSpec(t *testing.T) {
	res := Calculate(Input{
		Plan:               Plan{PricingModel: ModelFlat, BasePriceCents: 7900},
		Currency:           "USD",
		Coupons:            []Coupon{{Type: CouponFixed, AmountOffCents: 1000, Currency: "USD"}},
		TaxRateBasisPoints: 2000,
	})
	if res.SubtotalCents != 7900 {
		t.Fatalf("subtotal got %d want 7900", res.SubtotalCents)
	}
	if res.DiscountCents != 1000 {
		t.Fatalf("discount got %d want 1000", res.DiscountCents)
	}
	if res.TaxCents != 1380 {
		t.Fatalf("tax got %d want 1380", res.TaxCents)
	}
	if res.TotalCents != 8280 {
		t.Fatalf("total got %d want 8280", res.TotalCents)
	}
}

func TestPercentageCouponCapsAtSubtotal(t *testing.T) {
	res := Calculate(Input{
		Plan:     Plan{PricingModel: ModelFlat, BasePriceCents: 1000},
		Currency: "USD",
		Coupons:  []Coupon{{Type: CouponPercentage, PercentOff: 100}},
	})
	if res.DiscountCents != 1000 || res.TotalCents != 0 {
		t.Fatalf("100%% coupon: discount=%d total=%d want 1000/0", res.DiscountCents, res.TotalCents)
	}
}

func TestFixedCouponWrongCurrencyIgnored(t *testing.T) {
	res := Calculate(Input{
		Plan:     Plan{PricingModel: ModelFlat, BasePriceCents: 1000},
		Currency: "USD",
		Coupons:  []Coupon{{Type: CouponFixed, AmountOffCents: 500, Currency: "EUR"}},
	})
	if res.DiscountCents != 0 {
		t.Fatalf("currency-mismatched coupon should be ignored, got discount %d", res.DiscountCents)
	}
}

func TestInvalidPercentageIgnored(t *testing.T) {
	res := Calculate(Input{
		Plan:     Plan{PricingModel: ModelFlat, BasePriceCents: 1000},
		Currency: "USD",
		Coupons:  []Coupon{{Type: CouponPercentage, PercentOff: 150}},
	})
	if res.DiscountCents != 0 {
		t.Fatalf("out-of-range percentage should be ignored, got %d", res.DiscountCents)
	}
}

func TestTaxFor(t *testing.T) {
	if got := TaxFor(10000, 2100); got != 2100 {
		t.Fatalf("tax 21%% of 10000 = %d want 2100", got)
	}
	if got := TaxFor(0, 2100); got != 0 {
		t.Fatalf("tax of zero should be zero, got %d", got)
	}
}
