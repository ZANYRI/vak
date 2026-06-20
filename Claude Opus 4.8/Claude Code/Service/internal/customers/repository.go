package customers

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("customer not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, c *Customer) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO customers (user_id, email, name, company_name, billing_address,
			country, region, tax_id, currency)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at`,
		c.UserID, c.Email, c.Name, c.CompanyName, c.BillingAddress,
		c.Country, c.Region, c.TaxID, c.Currency,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Customer, error) {
	c := &Customer{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, email, name, company_name, billing_address,
			country, region, tax_id, currency, created_at, updated_at
		FROM customers WHERE id = $1`, id,
	).Scan(&c.ID, &c.UserID, &c.Email, &c.Name, &c.CompanyName, &c.BillingAddress,
		&c.Country, &c.Region, &c.TaxID, &c.Currency, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID) (*Customer, error) {
	c := &Customer{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, email, name, company_name, billing_address,
			country, region, tax_id, currency, created_at, updated_at
		FROM customers WHERE user_id = $1`, userID,
	).Scan(&c.ID, &c.UserID, &c.Email, &c.Name, &c.CompanyName, &c.BillingAddress,
		&c.Country, &c.Region, &c.TaxID, &c.Currency, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Customer, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, email, name, company_name, billing_address,
			country, region, tax_id, currency, created_at, updated_at
		FROM customers ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.UserID, &c.Email, &c.Name, &c.CompanyName, &c.BillingAddress,
			&c.Country, &c.Region, &c.TaxID, &c.Currency, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Customer, error) {
	_, err := r.db.Exec(ctx, `
		UPDATE customers SET
			name = COALESCE($2, name),
			company_name = COALESCE($3, company_name),
			billing_address = COALESCE($4, billing_address),
			country = COALESCE($5, country),
			region = COALESCE($6, region),
			tax_id = COALESCE($7, tax_id),
			currency = COALESCE($8, currency),
			updated_at = now()
		WHERE id = $1`,
		id, req.Name, req.CompanyName, req.BillingAddress,
		req.Country, req.Region, req.TaxID, req.Currency)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM customers WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
