package billing

import "testing"

func TestTaxCents(t *testing.T) {
	cases := []struct {
		amount int64
		bp     int64
		want   int64
	}{
		{10000, 2100, 2100},
		{1999, 1000, 199},
		{100, 0, 0},
	}
	for _, tc := range cases {
		got := TaxCents(tc.amount, tc.bp)
		if got != tc.want {
			t.Errorf("TaxCents(%d, %d) = %d, want %d", tc.amount, tc.bp, got, tc.want)
		}
	}
}

func TestApplyFixedDiscount(t *testing.T) {
	if got := ApplyFixedDiscount(1000, 200); got != 800 {
		t.Errorf("ApplyFixedDiscount(1000,200)=%d", got)
	}
	if got := ApplyFixedDiscount(100, 200); got != 0 {
		t.Errorf("ApplyFixedDiscount(100,200)=%d", got)
	}
}

func TestApplyPercentageDiscount(t *testing.T) {
	if got := ApplyPercentageDiscount(1000, 20); got != 800 {
		t.Errorf("ApplyPercentageDiscount(1000,20)=%d", got)
	}
}

func TestProrationFactor(t *testing.T) {
	if got := ProrationFactor(3000, 5000); got != 1500 {
		t.Errorf("ProrationFactor(3000,5000)=%d", got)
	}
}
