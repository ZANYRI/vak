package billing

import (
	"errors"
	"fmt"
)

type Tier struct {
	From           int64  `json:"from"`
	To             *int64 `json:"to"`
	UnitPriceCents int64  `json:"unit_price_cents"`
}
type Pricing struct {
	SeatPriceCents int64  `json:"seat_price_cents"`
	IncludedSeats  int64  `json:"included_seats"`
	IncludedUnits  int64  `json:"included_units"`
	UnitPriceCents int64  `json:"unit_price_cents"`
	Tiers          []Tier `json:"tiers"`
}
type Totals struct{ SubtotalCents, DiscountCents, TaxCents, TotalCents int64 }

func CalculateSubtotal(model string, base, quantity, usage int64, p Pricing) (int64, error) {
	if base < 0 || quantity < 0 || usage < 0 {
		return 0, errors.New("amounts and quantities must be non-negative")
	}
	total := base
	switch model {
	case "flat":
	case "per_seat":
		total += max(quantity-p.IncludedSeats, 0) * p.SeatPriceCents
	case "usage_based":
		total += max(usage-p.IncludedUnits, 0) * p.UnitPriceCents
	case "tiered":
		total += tierCharge(usage, p.Tiers)
	case "hybrid":
		total += max(quantity-p.IncludedSeats, 0)*p.SeatPriceCents + max(usage-p.IncludedUnits, 0)*p.UnitPriceCents + tierCharge(usage, p.Tiers)
	default:
		return 0, fmt.Errorf("unsupported pricing model %q", model)
	}
	if total < 0 {
		return 0, errors.New("amount overflow")
	}
	return total, nil
}
func tierCharge(usage int64, tiers []Tier) int64 {
	var total int64
	for _, tier := range tiers {
		if usage <= tier.From {
			continue
		}
		end := usage
		if tier.To != nil && *tier.To < end {
			end = *tier.To
		}
		if end > tier.From {
			total += (end - tier.From) * tier.UnitPriceCents
		}
	}
	return total
}
func Calculate(subtotal int64, percentOff *int, fixedOff int64, taxBPS int64) (Totals, error) {
	if subtotal < 0 || fixedOff < 0 || taxBPS < 0 || taxBPS > 100000 {
		return Totals{}, errors.New("invalid monetary input")
	}
	discount := fixedOff
	if percentOff != nil {
		if *percentOff < 1 || *percentOff > 100 {
			return Totals{}, errors.New("percent_off must be between 1 and 100")
		}
		discount += subtotal * int64(*percentOff) / 100
	}
	if discount > subtotal {
		discount = subtotal
	}
	taxable := subtotal - discount
	tax := taxable * taxBPS / 10000
	return Totals{subtotal, discount, tax, taxable + tax}, nil
}

// Proration returns the credit for unused old service and charge for new service, rounded toward zero.
func Proration(oldPrice, newPrice int64, periodStart, periodEnd, changedAtUnix int64) (credit, charge, difference int64, err error) {
	if oldPrice < 0 || newPrice < 0 || periodEnd <= periodStart {
		return 0, 0, 0, errors.New("invalid period or price")
	}
	remaining := max(periodEnd-changedAtUnix, 0)
	if remaining > periodEnd-periodStart {
		remaining = periodEnd - periodStart
	}
	credit = oldPrice * remaining / (periodEnd - periodStart)
	charge = newPrice * remaining / (periodEnd - periodStart)
	return credit, charge, charge - credit, nil
}
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
