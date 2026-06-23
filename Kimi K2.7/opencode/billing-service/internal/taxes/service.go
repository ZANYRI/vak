package taxes

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

// Service provides tax rule management.
type Service struct {
	pool *pgxpool.Pool
}

// NewService creates a tax service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// CreateRequest represents a tax rule creation payload.
type CreateRequest struct {
	Country            string `json:"country" validate:"required"`
	Region             string `json:"region,omitempty"`
	TaxName            string `json:"tax_name" validate:"required"`
	TaxRateBasisPoints int    `json:"tax_rate_basis_points" validate:"min=0,max=100000"`
	IsActive           bool   `json:"is_active"`
}

// Create inserts a tax rule.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*models.TaxRule, error) {
	t := &models.TaxRule{
		ID:                 uuid.New(),
		Country:            strings.ToUpper(req.Country),
		Region:             req.Region,
		TaxName:            req.TaxName,
		TaxRateBasisPoints: req.TaxRateBasisPoints,
		IsActive:           req.IsActive,
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO tax_rules (id, country, region, tax_name, tax_rate_basis_points, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		t.ID, t.Country, t.Region, t.TaxName, t.TaxRateBasisPoints, t.IsActive)
	if err != nil {
		return nil, fmt.Errorf("insert tax rule: %w", err)
	}
	return t, nil
}

// List returns tax rules.
func (s *Service) List(ctx context.Context, limit, offset int) ([]models.TaxRule, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := s.pool.Query(ctx,
		`SELECT id, country, region, tax_name, tax_rate_basis_points, is_active, created_at, updated_at
		 FROM tax_rules ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.TaxRule
	for rows.Next() {
		var t models.TaxRule
		err := rows.Scan(&t.ID, &t.Country, &t.Region, &t.TaxName, &t.TaxRateBasisPoints, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, t)
	}
	var total int64
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tax_rules`).Scan(&total)
	return list, total, rows.Err()
}

// GetForRegion returns the active tax rule matching country/region.
func (s *Service) GetForRegion(ctx context.Context, country, region string) (*models.TaxRule, error) {
	var t models.TaxRule
	err := s.pool.QueryRow(ctx,
		`SELECT id, country, region, tax_name, tax_rate_basis_points, is_active, created_at, updated_at
		 FROM tax_rules WHERE country=$1 AND (region='' OR region=$2) AND is_active=true
		 ORDER BY (region=$2) DESC LIMIT 1`, strings.ToUpper(country), region).Scan(
		&t.ID, &t.Country, &t.Region, &t.TaxName, &t.TaxRateBasisPoints, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// UpdateRequest updates mutable fields.
type UpdateRequest struct {
	TaxRateBasisPoints *int  `json:"tax_rate_basis_points,omitempty" validate:"omitempty,min=0,max=100000"`
	IsActive           *bool `json:"is_active,omitempty"`
}

// Update modifies a tax rule.
func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*models.TaxRule, error) {
	if req.TaxRateBasisPoints == nil && req.IsActive == nil {
		return s.GetByID(ctx, id)
	}
	t, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.TaxRateBasisPoints != nil {
		t.TaxRateBasisPoints = *req.TaxRateBasisPoints
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}
	_, err = s.pool.Exec(ctx,
		`UPDATE tax_rules SET tax_rate_basis_points=$1, is_active=$2, updated_at=NOW() WHERE id=$3`,
		t.TaxRateBasisPoints, t.IsActive, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// GetByID returns a tax rule by id.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.TaxRule, error) {
	var t models.TaxRule
	err := s.pool.QueryRow(ctx,
		`SELECT id, country, region, tax_name, tax_rate_basis_points, is_active, created_at, updated_at
		 FROM tax_rules WHERE id=$1`, id).Scan(
		&t.ID, &t.Country, &t.Region, &t.TaxName, &t.TaxRateBasisPoints, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tax rule not found")
		}
		return nil, err
	}
	return &t, nil
}

// Delete removes a tax rule.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM tax_rules WHERE id=$1`, id)
	return err
}
