package billing

import "errors"

// AddCents adds two money amounts safely.
func AddCents(a, b int64) int64 {
	return a + b
}

// SubCents subtracts b from a.
func SubCents(a, b int64) int64 {
	return a - b
}

// MulCents multiplies cents by an integer quantity.
func MulCents(cents int64, qty int64) int64 {
	return cents * qty
}

// TaxCents calculates tax using basis points.
func TaxCents(taxable, basisPoints int64) int64 {
	if basisPoints <= 0 {
		return 0
	}
	return (taxable * basisPoints) / 10000
}

// ApplyPercentageDiscount returns amount after percentage discount.
func ApplyPercentageDiscount(amount int64, percent int) int64 {
	if percent <= 0 || percent > 100 {
		return amount
	}
	discount := (amount * int64(percent)) / 100
	return amount - discount
}

// ApplyFixedDiscount returns amount after fixed discount, never negative.
func ApplyFixedDiscount(amount, discount int64) int64 {
	if discount <= 0 {
		return amount
	}
	if discount > amount {
		return 0
	}
	return amount - discount
}

// ProrationFactor returns a fraction of an amount for the given percentage (0-10000).
// pct is in ten-thousandths (50% = 5000).
func ProrationFactor(amount, pctTenThousandths int64) int64 {
	if pctTenThousandths <= 0 {
		return 0
	}
	if pctTenThousandths >= 10000 {
		return amount
	}
	return (amount * pctTenThousandths) / 10000
}

// ErrMoneyNegative is returned when a money calculation would go negative unexpectedly.
var ErrMoneyNegative = errors.New("money amount cannot be negative")
