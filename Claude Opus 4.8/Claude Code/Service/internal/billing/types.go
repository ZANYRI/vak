package billing

// Pricing models.
const (
	ModelFlat       = "flat"
	ModelPerSeat    = "per_seat"
	ModelUsageBased = "usage_based"
	ModelTiered     = "tiered"
	ModelHybrid     = "hybrid"
)

// Coupon types.
const (
	CouponPercentage = "percentage"
	CouponFixed      = "fixed_amount"
)

// Line types on an invoice.
const (
	LineBase      = "base"
	LineSeat      = "seat"
	LineUsage     = "usage"
	LineTier      = "tier"
	LineProration = "proration"
	LineDiscount  = "discount"
	LineTax       = "tax"
)

// Tier is one band of a tiered pricing model. To is nil for an unbounded top tier.
// Ranges are inclusive over cumulative units (e.g. 0–1000, 1001–10000, 10001–∞).
type Tier struct {
	From           int64
	To             *int64
	UnitPriceCents int64
}

// Plan captures the pricing inputs needed to compute charges.
type Plan struct {
	PricingModel   string
	Currency       string
	BasePriceCents int64
	SeatPriceCents int64
	IncludedSeats  int64
	IncludedUnits  int64
	UnitPriceCents int64
	Tiers          []Tier
}

// Coupon is a discount input.
type Coupon struct {
	Type           string
	PercentOff     int
	AmountOffCents int64
	Currency       string
}

// Input drives a single invoice calculation.
type Input struct {
	Plan           Plan
	Currency       string
	Seats          int64 // subscription quantity
	UsageUnits     int64 // total usage for the metric in the period
	Coupons        []Coupon
	TaxRateBasisPoints int
	ExtraLines     []Line // e.g. proration lines computed elsewhere
}

// Line is a computed invoice line.
type Line struct {
	Type            string `json:"type"`
	Description     string `json:"description"`
	Quantity        int64  `json:"quantity"`
	UnitAmountCents int64  `json:"unit_amount_cents"`
	AmountCents     int64  `json:"amount_cents"`
}

// Result is the outcome of a calculation.
type Result struct {
	Lines         []Line `json:"lines"`
	SubtotalCents int64  `json:"subtotal_cents"`
	DiscountCents int64  `json:"discount_cents"`
	TaxCents      int64  `json:"tax_cents"`
	TotalCents    int64  `json:"total_cents"`
}
