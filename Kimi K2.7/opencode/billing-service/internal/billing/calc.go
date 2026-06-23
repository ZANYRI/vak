package billing

import (
	"errors"
	"sort"
	"time"

	"billing-service/internal/models"
	"github.com/google/uuid"
)

// LineInput describes a line item to add to an invoice.
type LineInput struct {
	Type        string
	Description string
	Quantity    int64
	UnitAmount  int64
	Amount      int64
	Metadata    map[string]interface{}
}

// CalculationInput holds everything needed to generate an invoice.
type CalculationInput struct {
	Customer           *models.Customer
	Subscription       *models.Subscription
	Plan               *models.Plan
	UsageQuantity      int64
	Coupon             *models.Coupon
	Proration          *ProrationCharge
	TaxRateBasisPoints int64
}

// ProrationCharge represents a proration adjustment.
type ProrationCharge struct {
	Description string
	AmountCents int64
}

// CalculationResult holds the computed invoice totals and lines.
type CalculationResult struct {
	Lines         []models.InvoiceLine
	SubtotalCents int64
	DiscountCents int64
	TaxCents      int64
	TotalCents    int64
}

// CalculateInvoice computes charges for a subscription invoice.
func CalculateInvoice(in CalculationInput, now time.Time) (*CalculationResult, error) {
	if in.Customer == nil || in.Subscription == nil || in.Plan == nil {
		return nil, errors.New("customer, subscription and plan are required")
	}

	res := &CalculationResult{Lines: []models.InvoiceLine{}}
	res.SubtotalCents += in.Plan.BasePriceCents
	addLine(res, "base", "Base plan fee", 1, in.Plan.BasePriceCents, in.Plan.BasePriceCents, nil)

	for _, price := range in.Plan.Prices {
		switch price.Model {
		case models.PlanPerSeat:
			if price.SeatPriceCents != nil && price.IncludedSeats != nil {
				extraSeats := int64(in.Subscription.Seats - *price.IncludedSeats)
				if extraSeats < 0 {
					extraSeats = 0
				}
				amount := MulCents(*price.SeatPriceCents, extraSeats)
				res.SubtotalCents += amount
				addLine(res, "seat", "Per-seat charge", extraSeats, *price.SeatPriceCents, amount, nil)
			}
		case models.PlanUsageBased:
			if price.UnitPriceCents != nil {
				included := int64(0)
				if price.IncludedUnits != nil {
					included = *price.IncludedUnits
				}
				billable := in.UsageQuantity - included
				if billable < 0 {
					billable = 0
				}
				amount := MulCents(*price.UnitPriceCents, billable)
				res.SubtotalCents += amount
				addLine(res, "usage", "Usage charge", billable, *price.UnitPriceCents, amount, nil)
			}
		case models.PlanHybrid:
			if price.UnitPriceCents != nil {
				amount := MulCents(*price.UnitPriceCents, in.UsageQuantity)
				res.SubtotalCents += amount
				addLine(res, "usage", "Usage charge", in.UsageQuantity, *price.UnitPriceCents, amount, nil)
			}
		}
	}

	if len(in.Plan.Tiers) > 0 {
		tiered := calculateTiered(in.Plan.Tiers, in.UsageQuantity)
		if tiered > 0 {
			res.SubtotalCents += tiered
			addLine(res, "tiered", "Tiered usage charge", in.UsageQuantity, 0, tiered, nil)
		}
	}

	if in.Proration != nil {
		res.SubtotalCents += in.Proration.AmountCents
		addLine(res, "proration", in.Proration.Description, 1, in.Proration.AmountCents, in.Proration.AmountCents, nil)
	}

	// Coupons
	if in.Coupon != nil {
		discount := calculateCouponDiscount(in.Coupon, res.SubtotalCents, in.Plan.Currency)
		res.DiscountCents += discount
		if discount > 0 {
			addLine(res, "discount", "Coupon discount", 1, -discount, -discount, map[string]interface{}{"coupon_id": in.Coupon.ID})
		}
	}
	 taxableAfterDiscount := ApplyFixedDiscount(res.SubtotalCents, res.DiscountCents)

	// Taxes
	res.TaxCents = TaxCents(taxableAfterDiscount, in.TaxRateBasisPoints)

	res.TotalCents = taxableAfterDiscount + res.TaxCents
	return res, nil
}

func calculateTiered(tiers []models.PricingTier, qty int64) int64 {
	sort.Slice(tiers, func(i, j int) bool { return tiers[i].From < tiers[j].From })
	total := int64(0)
	consumed := int64(0)
	for _, tier := range tiers {
		if consumed >= qty {
			break
		}
		if tier.From >= qty {
			break
		}
		end := qty
		if tier.To != nil {
			end = *tier.To
		}
		if end > qty {
			end = qty
		}
		tierQty := end - tier.From
		if tierQty <= 0 {
			continue
		}
		remaining := qty - consumed
		if tierQty > remaining {
			tierQty = remaining
		}
		total += MulCents(tier.UnitPriceCents, tierQty)
		consumed += tierQty
	}
	return total
}

func calculateCouponDiscount(coupon *models.Coupon, subtotal int64, currency string) int64 {
	if coupon == nil || !coupon.IsActive {
		return 0
	}
	now := time.Now()
	if coupon.ValidFrom.After(now) || (coupon.ValidUntil != nil && coupon.ValidUntil.Before(now)) {
		return 0
	}
	if coupon.MaxRedemptions != nil && coupon.TimesRedeemed >= *coupon.MaxRedemptions {
		return 0
	}
	switch coupon.Type {
	case models.CouponPercentage:
		if coupon.PercentOff != nil {
			return (subtotal * int64(*coupon.PercentOff)) / 100
		}
	case models.CouponFixedAmount:
		if coupon.AmountOffCents != nil && (coupon.Currency == nil || *coupon.Currency == currency) {
			if *coupon.AmountOffCents > subtotal {
				return subtotal
			}
			return *coupon.AmountOffCents
		}
	}
	return 0
}

func addLine(res *CalculationResult, typ, desc string, qty, unit, amount int64, meta map[string]interface{}) {
	res.Lines = append(res.Lines, models.InvoiceLine{
		ID:              uuid.New(),
		Type:            typ,
		Description:     desc,
		Quantity:        qty,
		UnitAmountCents: unit,
		AmountCents:     amount,
		Metadata:        meta,
	})
}

// ProrationAmount calculates the net amount to charge for a plan change.
func ProrationAmount(oldPlan, newPlan *models.Plan, daysUsed, totalDays int) int64 {
	if totalDays <= 0 {
		return 0
	}
	unusedPct := int64((totalDays - daysUsed) * 10000 / totalDays)
	unusedCredit := ProrationFactor(oldPlan.BasePriceCents, unusedPct)
	newCharge := ProrationFactor(newPlan.BasePriceCents, unusedPct)
	return newCharge - unusedCredit
}
