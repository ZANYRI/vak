package customers

import (
	"time"

	"github.com/google/uuid"
)

// Customer is a billing customer.
type Customer struct {
	ID             uuid.UUID  `json:"id"`
	UserID         *uuid.UUID `json:"user_id"`
	Email          string     `json:"email"`
	Name           string     `json:"name"`
	CompanyName    *string    `json:"company_name"`
	BillingAddress *string    `json:"billing_address"`
	Country        string     `json:"country"`
	Region         string     `json:"region"`
	TaxID          *string    `json:"tax_id"`
	Currency       string     `json:"currency"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CreateRequest is the create-customer payload.
type CreateRequest struct {
	Email          string     `json:"email" validate:"required,email"`
	Name           string     `json:"name" validate:"required"`
	CompanyName    *string    `json:"company_name"`
	BillingAddress *string    `json:"billing_address"`
	Country        string     `json:"country"`
	Region         string     `json:"region"`
	TaxID          *string    `json:"tax_id"`
	Currency       string     `json:"currency" validate:"required,len=3"`
	UserID         *uuid.UUID `json:"user_id"`
}

// UpdateRequest patches mutable customer fields.
type UpdateRequest struct {
	Name           *string `json:"name"`
	CompanyName    *string `json:"company_name"`
	BillingAddress *string `json:"billing_address"`
	Country        *string `json:"country"`
	Region         *string `json:"region"`
	TaxID          *string `json:"tax_id"`
	Currency       *string `json:"currency"`
}
