package coupons

import (
	"time"

	"github.com/google/uuid"
)

// Coupon is a discount that can be applied to subscriptions.
type Coupon struct {
	ID             uuid.UUID  `json:"id"`
	Code           string     `json:"code"`
	Type           string     `json:"type"`
	PercentOff     *int       `json:"percent_off"`
	AmountOffCents *int64     `json:"amount_off_cents"`
	Currency       *string    `json:"currency"`
	MaxRedemptions *int       `json:"max_redemptions"`
	TimesRedeemed  int        `json:"times_redeemed"`
	ValidFrom      *time.Time `json:"valid_from"`
	ValidUntil     *time.Time `json:"valid_until"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CreateRequest is the create-coupon payload.
type CreateRequest struct {
	Code           string     `json:"code" validate:"required"`
	Type           string     `json:"type" validate:"required,oneof=percentage fixed_amount"`
	PercentOff     *int       `json:"percent_off"`
	AmountOffCents *int64     `json:"amount_off_cents"`
	Currency       *string    `json:"currency"`
	MaxRedemptions *int       `json:"max_redemptions"`
	ValidFrom      *time.Time `json:"valid_from"`
	ValidUntil     *time.Time `json:"valid_until"`
}

// UpdateRequest patches mutable coupon fields.
type UpdateRequest struct {
	MaxRedemptions *int       `json:"max_redemptions"`
	ValidUntil     *time.Time `json:"valid_until"`
	IsActive       *bool      `json:"is_active"`
}

// ApplyRequest is the body for applying a coupon to a subscription.
type ApplyRequest struct {
	Code     string `json:"code" validate:"required"`
	Currency string `json:"currency" validate:"required,len=3"`
}
