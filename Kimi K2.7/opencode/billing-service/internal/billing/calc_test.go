package billing

import (
	"testing"
	"time"

	"billing-service/internal/models"
)

func TestCalculateInvoiceFlatPlan(t *testing.T) {
	plan := &models.Plan{
		ID:             [16]byte{1},
		Currency:       "USD",
		BasePriceCents: 1900,
	}
	sub := &models.Subscription{ID: [16]byte{2}, Quantity: 1, Seats: 1}
	cust := &models.Customer{ID: [16]byte{3}, Country: "US"}
	res, err := CalculateInvoice(CalculationInput{Customer: cust, Subscription: sub, Plan: plan, TaxRateBasisPoints: 0}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if res.SubtotalCents != 1900 {
		t.Errorf("subtotal = %d, want 1900", res.SubtotalCents)
	}
	if res.TotalCents != 1900 {
		t.Errorf("total = %d, want 1900", res.TotalCents)
	}
}

func TestCalculateInvoicePerSeat(t *testing.T) {
	included := 3
	seatPrice := int64(700)
	plan := &models.Plan{
		ID:             [16]byte{1},
		Currency:       "USD",
		BasePriceCents: 1200,
		Prices: []models.PlanPrice{
			{Model: models.PlanPerSeat, IncludedSeats: &included, SeatPriceCents: &seatPrice},
		},
	}
	sub := &models.Subscription{ID: [16]byte{2}, Quantity: 1, Seats: 5}
	cust := &models.Customer{ID: [16]byte{3}, Country: "US"}
	res, err := CalculateInvoice(CalculationInput{Customer: cust, Subscription: sub, Plan: plan, TaxRateBasisPoints: 0}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if res.SubtotalCents != 1200+1400 {
		t.Errorf("subtotal = %d, want 2600", res.SubtotalCents)
	}
}

func TestCalculateInvoiceTiered(t *testing.T) {
	to1 := int64(1000)
	to2 := int64(10000)
	plan := &models.Plan{
		ID:             [16]byte{1},
		Currency:       "USD",
		BasePriceCents: 0,
		Tiers: []models.PricingTier{
			{From: 0, To: &to1, UnitPriceCents: 5},
			{From: 1000, To: &to2, UnitPriceCents: 3},
			{From: 10000, To: nil, UnitPriceCents: 1},
		},
	}
	sub := &models.Subscription{ID: [16]byte{2}}
	cust := &models.Customer{ID: [16]byte{3}}
	res, err := CalculateInvoice(CalculationInput{Customer: cust, Subscription: sub, Plan: plan, UsageQuantity: 15000, TaxRateBasisPoints: 0}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	want := int64(1000*5 + 9000*3 + 5000*1) // 5000+27000+5000=37000
	if res.SubtotalCents != want {
		t.Errorf("subtotal = %d, want %d", res.SubtotalCents, want)
	}
}

func TestProrationAmount(t *testing.T) {
	old := &models.Plan{BasePriceCents: 3000}
	new := &models.Plan{BasePriceCents: 6000}
	got := ProrationAmount(old, new, 15, 30) // half period unused on new
	if got <= 0 {
		t.Errorf("expected positive proration, got %d", got)
	}
}
