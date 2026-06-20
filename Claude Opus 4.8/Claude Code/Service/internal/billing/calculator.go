package billing

import "fmt"

// Calculate computes invoice lines and totals from the given input.
// All amounts are int64 minor units (cents); no floating point is used.
func Calculate(in Input) Result {
	var lines []Line

	// 1. Base + seat + usage charges depending on the pricing model.
	lines = append(lines, baseLines(in)...)

	// 2. Extra (proration) lines supplied by the caller.
	lines = append(lines, in.ExtraLines...)

	// 3. Subtotal = sum of all charge lines (proration may be negative).
	var subtotal int64
	for _, l := range lines {
		subtotal += l.AmountCents
	}
	if subtotal < 0 {
		subtotal = 0
	}

	// 4. Discounts from coupons, capped so the discounted subtotal is >= 0.
	discount := applyCoupons(subtotal, in.Currency, in.Coupons)
	if discount > 0 {
		lines = append(lines, Line{
			Type:        LineDiscount,
			Description: "Discount",
			Quantity:    1,
			AmountCents: -discount,
		})
	}

	taxable := subtotal - discount
	if taxable < 0 {
		taxable = 0
	}

	// 5. Tax on the discounted amount.
	tax := TaxFor(taxable, in.TaxRateBasisPoints)
	if tax > 0 {
		lines = append(lines, Line{
			Type:        LineTax,
			Description: fmt.Sprintf("Tax (%.2f%%)", float64(in.TaxRateBasisPoints)/100),
			Quantity:    1,
			AmountCents: tax,
		})
	}

	total := taxable + tax

	return Result{
		Lines:         lines,
		SubtotalCents: subtotal,
		DiscountCents: discount,
		TaxCents:      tax,
		TotalCents:    total,
	}
}

func baseLines(in Input) []Line {
	var lines []Line
	p := in.Plan

	switch p.PricingModel {
	case ModelFlat:
		lines = append(lines, baseLine(p.BasePriceCents))
	case ModelPerSeat:
		lines = append(lines, baseLine(p.BasePriceCents))
		lines = append(lines, seatLine(p, in.Seats)...)
	case ModelUsageBased:
		lines = append(lines, baseLine(p.BasePriceCents))
		lines = append(lines, usageLine(p, in.UsageUnits)...)
	case ModelTiered:
		lines = append(lines, baseLine(p.BasePriceCents))
		lines = append(lines, tieredLines(p.Tiers, in.UsageUnits)...)
	case ModelHybrid:
		lines = append(lines, baseLine(p.BasePriceCents))
		lines = append(lines, seatLine(p, in.Seats)...)
		if len(p.Tiers) > 0 {
			lines = append(lines, tieredLines(p.Tiers, in.UsageUnits)...)
		} else {
			lines = append(lines, usageLine(p, in.UsageUnits)...)
		}
	default:
		lines = append(lines, baseLine(p.BasePriceCents))
	}

	// Drop zero-amount base line only if it would be confusing; keep it for clarity.
	return lines
}

func baseLine(amount int64) Line {
	return Line{Type: LineBase, Description: "Base subscription fee", Quantity: 1,
		UnitAmountCents: amount, AmountCents: amount}
}

func seatLine(p Plan, seats int64) []Line {
	billable := seats - p.IncludedSeats
	if billable <= 0 {
		return nil
	}
	return []Line{{
		Type:            LineSeat,
		Description:     "Additional seats",
		Quantity:        billable,
		UnitAmountCents: p.SeatPriceCents,
		AmountCents:     billable * p.SeatPriceCents,
	}}
}

func usageLine(p Plan, usage int64) []Line {
	billable := usage - p.IncludedUnits
	if billable <= 0 {
		return nil
	}
	return []Line{{
		Type:            LineUsage,
		Description:     "Usage charges",
		Quantity:        billable,
		UnitAmountCents: p.UnitPriceCents,
		AmountCents:     billable * p.UnitPriceCents,
	}}
}

// tieredLines computes graduated tiered charges. Each tier bills the units that
// fall within its inclusive [From, To] cumulative range at its unit price.
func tieredLines(tiers []Tier, usage int64) []Line {
	var lines []Line
	for i, t := range tiers {
		upper := usage
		if t.To != nil && *t.To < usage {
			upper = *t.To
		}
		if upper < t.From {
			continue
		}
		units := upper - t.From + 1
		if units <= 0 {
			continue
		}
		lines = append(lines, Line{
			Type:            LineTier,
			Description:     fmt.Sprintf("Tier %d usage", i+1),
			Quantity:        units,
			UnitAmountCents: t.UnitPriceCents,
			AmountCents:     units * t.UnitPriceCents,
		})
	}
	return lines
}

// applyCoupons returns the total discount in cents, never exceeding subtotal.
func applyCoupons(subtotal int64, currency string, coupons []Coupon) int64 {
	remaining := subtotal
	var total int64
	for _, c := range coupons {
		var d int64
		switch c.Type {
		case CouponPercentage:
			if c.PercentOff < 1 || c.PercentOff > 100 {
				continue
			}
			d = remaining * int64(c.PercentOff) / 100
		case CouponFixed:
			if c.Currency != "" && c.Currency != currency {
				continue // currency mismatch: cannot apply
			}
			d = c.AmountOffCents
		}
		if d > remaining {
			d = remaining
		}
		if d < 0 {
			d = 0
		}
		total += d
		remaining -= d
		if remaining <= 0 {
			break
		}
	}
	return total
}

// TaxFor computes tax = amount * basisPoints / 10000 using integer math.
func TaxFor(amountCents int64, basisPoints int) int64 {
	if basisPoints <= 0 || amountCents <= 0 {
		return 0
	}
	return amountCents * int64(basisPoints) / 10000
}
