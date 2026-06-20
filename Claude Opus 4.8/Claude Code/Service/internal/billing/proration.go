package billing

import "time"

// Proration captures the result of a mid-period plan change.
type Proration struct {
	UnusedOldCreditCents int64 // credit for the unused portion of the old plan
	RemainingNewChargeCents int64 // charge for the remaining portion of the new plan
	DifferenceCents      int64 // net amount (charge - credit)
	Lines                []Line
}

// Prorate computes the proration when a subscription changes plan partway
// through a billing period. It uses time-based proportions with integer math.
//
//	oldPrice  – price already billed for the current period
//	newPrice  – full-period price of the new plan
//	periodStart/periodEnd – the current billing period bounds
//	changeAt  – the moment of the plan change
func Prorate(oldPriceCents, newPriceCents int64, periodStart, periodEnd, changeAt time.Time) Proration {
	totalSeconds := int64(periodEnd.Sub(periodStart).Seconds())
	if totalSeconds <= 0 {
		return Proration{}
	}
	if changeAt.Before(periodStart) {
		changeAt = periodStart
	}
	if changeAt.After(periodEnd) {
		changeAt = periodEnd
	}
	remainingSeconds := int64(periodEnd.Sub(changeAt).Seconds())
	if remainingSeconds < 0 {
		remainingSeconds = 0
	}

	credit := oldPriceCents * remainingSeconds / totalSeconds
	charge := newPriceCents * remainingSeconds / totalSeconds
	diff := charge - credit

	lines := []Line{
		{
			Type:        LineProration,
			Description: "Unused time on previous plan",
			Quantity:    1,
			AmountCents: -credit,
		},
		{
			Type:        LineProration,
			Description: "Remaining time on new plan",
			Quantity:    1,
			AmountCents: charge,
		},
	}

	return Proration{
		UnusedOldCreditCents:    credit,
		RemainingNewChargeCents: charge,
		DifferenceCents:         diff,
		Lines:                   lines,
	}
}
