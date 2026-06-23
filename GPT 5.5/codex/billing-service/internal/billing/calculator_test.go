package billing

import "testing"

func TestCalculateSubtotalTiered(t *testing.T) {
	to := int64(10000)
	got, err := CalculateSubtotal("tiered", 0, 0, 12000, Pricing{Tiers: []Tier{{From: 0, To: &to, UnitPriceCents: 5}, {From: 10000, To: nil, UnitPriceCents: 1}}})
	if err != nil || got != 52000 {
		t.Fatalf("got %d, %v", got, err)
	}
}
func TestDiscountAndTax(t *testing.T) {
	p := 10
	got, err := Calculate(7900, &p, 0, 2100)
	if err != nil || got.DiscountCents != 790 || got.TaxCents != 1493 || got.TotalCents != 8603 {
		t.Fatalf("%+v %v", got, err)
	}
}
func TestProration(t *testing.T) {
	c, n, d, e := Proration(3000, 6000, 0, 30, 15)
	if e != nil || c != 1500 || n != 3000 || d != 1500 {
		t.Fatalf("%d %d %d %v", c, n, d, e)
	}
}
