package taxes

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("tax rule not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, t *TaxRule) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO tax_rules (country, region, tax_name, tax_rate_basis_points, is_active)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, created_at, updated_at`,
		t.Country, t.Region, t.TaxName, t.TaxRateBasisPoints, t.IsActive,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*TaxRule, error) {
	t := &TaxRule{}
	err := r.db.QueryRow(ctx, `
		SELECT id, country, region, tax_name, tax_rate_basis_points, is_active, created_at, updated_at
		FROM tax_rules WHERE id = $1`, id,
	).Scan(&t.ID, &t.Country, &t.Region, &t.TaxName, &t.TaxRateBasisPoints, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]TaxRule, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, country, region, tax_name, tax_rate_basis_points, is_active, created_at, updated_at
		FROM tax_rules ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []TaxRule
	for rows.Next() {
		var t TaxRule
		if err := rows.Scan(&t.ID, &t.Country, &t.Region, &t.TaxName, &t.TaxRateBasisPoints, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*TaxRule, error) {
	_, err := r.db.Exec(ctx, `
		UPDATE tax_rules SET
			tax_name = COALESCE($2, tax_name),
			tax_rate_basis_points = COALESCE($3, tax_rate_basis_points),
			is_active = COALESCE($4, is_active),
			updated_at = now()
		WHERE id = $1`,
		id, req.TaxName, req.TaxRateBasisPoints, req.IsActive)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM tax_rules WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// FindActiveFor returns the active rule matching (country, region). If a region
// is given but no region-specific rule exists, it falls back to the country-level
// rule (region = ''). Returns ErrNotFound if no active rule matches.
func (r *Repository) FindActiveFor(ctx context.Context, country, region string) (*TaxRule, error) {
	if region != "" {
		t, err := r.findActive(ctx, country, region)
		if err == nil {
			return t, nil
		}
		if !errors.Is(err, ErrNotFound) {
			return nil, err
		}
		// Fall back to country-level rule.
	}
	return r.findActive(ctx, country, "")
}

func (r *Repository) findActive(ctx context.Context, country, region string) (*TaxRule, error) {
	t := &TaxRule{}
	err := r.db.QueryRow(ctx, `
		SELECT id, country, region, tax_name, tax_rate_basis_points, is_active, created_at, updated_at
		FROM tax_rules
		WHERE country = $1 AND region = $2 AND is_active = true
		ORDER BY created_at DESC LIMIT 1`, country, region,
	).Scan(&t.ID, &t.Country, &t.Region, &t.TaxName, &t.TaxRateBasisPoints, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}
