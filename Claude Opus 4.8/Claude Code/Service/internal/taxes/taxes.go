package taxes

import (
	"time"

	"github.com/google/uuid"
)

// TaxRule is a tax rule applied for a country/region.
type TaxRule struct {
	ID                 uuid.UUID `json:"id"`
	Country            string    `json:"country"`
	Region             string    `json:"region"`
	TaxName            string    `json:"tax_name"`
	TaxRateBasisPoints int       `json:"tax_rate_basis_points"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// CreateRequest is the create-tax-rule payload.
type CreateRequest struct {
	Country            string `json:"country" validate:"required,len=2"`
	Region             string `json:"region"`
	TaxName            string `json:"tax_name" validate:"required"`
	TaxRateBasisPoints int    `json:"tax_rate_basis_points" validate:"gte=0,lte=10000"`
}

// UpdateRequest patches mutable tax-rule fields.
type UpdateRequest struct {
	TaxName            *string `json:"tax_name"`
	TaxRateBasisPoints *int    `json:"tax_rate_basis_points"`
	IsActive           *bool   `json:"is_active"`
}
