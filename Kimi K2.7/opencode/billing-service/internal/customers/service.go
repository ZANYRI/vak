package customers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"billing-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service provides customer management operations.
type Service struct {
	pool *pgxpool.Pool
}

// NewService creates a customer service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// CreateRequest represents a customer creation payload.
type CreateRequest struct {
	UserID         *string `json:"user_id,omitempty"`
	Email          string  `json:"email" validate:"required,email"`
	Name           string  `json:"name" validate:"required"`
	CompanyName    string  `json:"company_name,omitempty"`
	BillingAddress string  `json:"billing_address,omitempty"`
	Country        string  `json:"country" validate:"required"`
	Region         string  `json:"region,omitempty"`
	TaxID          string  `json:"tax_id,omitempty"`
	Currency       string  `json:"currency" validate:"required,len=3"`
}

// Create inserts a new customer.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*models.Customer, error) {
	c := &models.Customer{
		ID:             uuid.New(),
		Email:          strings.ToLower(req.Email),
		Name:           req.Name,
		CompanyName:    req.CompanyName,
		BillingAddress: req.BillingAddress,
		Country:        strings.ToUpper(req.Country),
		Region:         req.Region,
		TaxID:          req.TaxID,
		Currency:       strings.ToUpper(req.Currency),
	}
	if req.UserID != nil {
		uid, err := uuid.Parse(*req.UserID)
		if err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		c.UserID = &uid
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO customers (id, user_id, email, name, company_name, billing_address, country, region, tax_id, currency)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		c.ID, c.UserID, c.Email, c.Name, c.CompanyName, c.BillingAddress, c.Country, c.Region, c.TaxID, c.Currency)
	if err != nil {
		return nil, fmt.Errorf("insert customer: %w", err)
	}
	return c, nil
}

// List returns all customers paginated.
func (s *Service) List(ctx context.Context, limit, offset int) ([]models.Customer, int64, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, email, name, company_name, billing_address, country, region, tax_id, currency, created_at, updated_at
		 FROM customers ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Customer
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(&c.ID, &c.UserID, &c.Email, &c.Name, &c.CompanyName, &c.BillingAddress, &c.Country, &c.Region, &c.TaxID, &c.Currency, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, c)
	}

	var total int64
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM customers`).Scan(&total)
	return list, total, rows.Err()
}

// GetByID returns a customer by id.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	var c models.Customer
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, email, name, company_name, billing_address, country, region, tax_id, currency, created_at, updated_at
		 FROM customers WHERE id = $1`, id).Scan(
		&c.ID, &c.UserID, &c.Email, &c.Name, &c.CompanyName, &c.BillingAddress, &c.Country, &c.Region, &c.TaxID, &c.Currency, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, err
	}
	return &c, nil
}

// UpdateRequest is the payload to update a customer.
type UpdateRequest struct {
	Email          *string `json:"email,omitempty" validate:"omitempty,email"`
	Name           *string `json:"name,omitempty"`
	CompanyName    *string `json:"company_name,omitempty"`
	BillingAddress *string `json:"billing_address,omitempty"`
	Country        *string `json:"country,omitempty" validate:"omitempty,len=3"`
	Region         *string `json:"region,omitempty"`
	TaxID          *string `json:"tax_id,omitempty"`
	Currency       *string `json:"currency,omitempty" validate:"omitempty,len=3"`
}

// Update modifies a customer.
func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*models.Customer, error) {
	c, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Email != nil {
		c.Email = strings.ToLower(*req.Email)
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.CompanyName != nil {
		c.CompanyName = *req.CompanyName
	}
	if req.BillingAddress != nil {
		c.BillingAddress = *req.BillingAddress
	}
	if req.Country != nil {
		c.Country = strings.ToUpper(*req.Country)
	}
	if req.Region != nil {
		c.Region = *req.Region
	}
	if req.TaxID != nil {
		c.TaxID = *req.TaxID
	}
	if req.Currency != nil {
		c.Currency = strings.ToUpper(*req.Currency)
	}
	_, err = s.pool.Exec(ctx,
		`UPDATE customers SET email=$1, name=$2, company_name=$3, billing_address=$4, country=$5, region=$6, tax_id=$7, currency=$8, updated_at=NOW()
		 WHERE id=$9`,
		c.Email, c.Name, c.CompanyName, c.BillingAddress, c.Country, c.Region, c.TaxID, c.Currency, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// Delete removes a customer.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := s.pool.Exec(ctx, `DELETE FROM customers WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("customer not found")
	}
	return nil
}
