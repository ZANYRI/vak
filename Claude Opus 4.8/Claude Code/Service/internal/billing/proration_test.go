package billing

import (
	"testing"
	"time"
)

// TestProrationHalfPeriodUpgrade reproduces the spec example: a 3000 monthly
// plan upgraded halfway to a 6000 plan yields a 1500 proration difference.
func TestProrationHalfPeriodUpgrade(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 30)
	change := start.AddDate(0, 0, 15) // halfway

	p := Prorate(3000, 6000, start, end, change)

	if p.UnusedOldCreditCents != 1500 {
		t.Fatalf("credit got %d want 1500", p.UnusedOldCreditCents)
	}
	if p.RemainingNewChargeCents != 3000 {
		t.Fatalf("charge got %d want 3000", p.RemainingNewChargeCents)
	}
	if p.DifferenceCents != 1500 {
		t.Fatalf("difference got %d want 1500", p.DifferenceCents)
	}
	if len(p.Lines) != 2 {
		t.Fatalf("expected 2 proration lines, got %d", len(p.Lines))
	}
}

func TestProrationNoRemainingTime(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 30)
	p := Prorate(3000, 6000, start, end, end) // change at the very end
	if p.DifferenceCents != 0 {
		t.Fatalf("no remaining time should yield 0 difference, got %d", p.DifferenceCents)
	}
}

func TestProrationDowngrade(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 30)
	change := start.AddDate(0, 0, 15)
	p := Prorate(6000, 3000, start, end, change)
	// credit 3000, charge 1500 => negative difference (a credit to the customer)
	if p.DifferenceCents != -1500 {
		t.Fatalf("downgrade difference got %d want -1500", p.DifferenceCents)
	}
}
