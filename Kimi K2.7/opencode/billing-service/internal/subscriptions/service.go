package subscriptions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"billing-service/internal/billing"
	"billing-service/internal/models"
	"billing-service/internal/plans"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service provides subscription lifecycle management.
type Service struct {
	pool *pgxpool.Pool
	plans *plans.Service
}

// NewService creates a subscription service.
func NewService(pool *pgxpool.Pool, p *plans.Service) *Service {
	return &Service{pool: pool, plans: p}
}

// CreateRequest creates a subscription.
type CreateRequest struct {
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	PlanID     uuid.UUID `json:"plan_id" validate:"required"`
	Quantity   int       `json:"quantity" validate:"min=1"`
	Seats      int       `json:"seats" validate:"min=1"`
}

// Create inserts a subscription and sets initial period based on plan trial.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*models.Subscription, error) {
	plan, err := s.plans.GetByID(ctx, req.PlanID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	sub := &models.Subscription{
		ID:       uuid.New(),
		CustomerID: req.CustomerID,
		PlanID:   req.PlanID,
		Quantity: req.Quantity,
		Seats:    req.Seats,
	}
	if plan.TrialDays > 0 {
		sub.Status = models.StatusTrialing
		trialEnd := now.AddDate(0, 0, plan.TrialDays)
		sub.TrialStart = &now
		sub.TrialEnd = &trialEnd
		sub.CurrentPeriodStart = trialEnd
		sub.CurrentPeriodEnd = addInterval(trialEnd, plan.BillingInterval)
	} else {
		sub.Status = models.StatusActive
		sub.CurrentPeriodStart = now
		sub.CurrentPeriodEnd = addInterval(now, plan.BillingInterval)
	}

	_, err = s.pool.Exec(ctx,
		`INSERT INTO subscriptions (id, customer_id, plan_id, status, quantity, seats, current_period_start, current_period_end, trial_start, trial_end)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		sub.ID, sub.CustomerID, sub.PlanID, sub.Status, sub.Quantity, sub.Seats, sub.CurrentPeriodStart, sub.CurrentPeriodEnd, sub.TrialStart, sub.TrialEnd)
	if err != nil {
		return nil, fmt.Errorf("insert subscription: %w", err)
	}
	return s.GetByID(ctx, sub.ID)
}

// List returns subscriptions. If customerID is provided, restricts to that customer.
func (s *Service) List(ctx context.Context, customerID *uuid.UUID, limit, offset int) ([]models.Subscription, int64, error) {
	query := `SELECT id, customer_id, plan_id, status, quantity, seats, current_period_start, current_period_end, trial_start, trial_end, cancel_at_period_end, created_at, updated_at FROM subscriptions`
	countQuery := `SELECT COUNT(*) FROM subscriptions`
	args := []interface{}{}
	where := ""
	if customerID != nil {
		args = append(args, *customerID)
		where = fmt.Sprintf(" WHERE customer_id = $%d", len(args))
	}
	query += where + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(&sub.ID, &sub.CustomerID, &sub.PlanID, &sub.Status, &sub.Quantity, &sub.Seats, &sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.TrialStart, &sub.TrialEnd, &sub.CancelAtPeriodEnd, &sub.CreatedAt, &sub.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, sub)
	}
	var countArgs []interface{}
	if customerID != nil {
		countArgs = append(countArgs, *customerID)
	}
	var total int64
	_ = s.pool.QueryRow(ctx, countQuery+where, countArgs...).Scan(&total)
	for i := range list {
		p, _ := s.plans.GetByID(ctx, list[i].PlanID)
		list[i].Plan = p
	}
	return list, total, rows.Err()
}

// GetByID returns a subscription by id with embedded plan.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	var sub models.Subscription
	err := s.pool.QueryRow(ctx,
		`SELECT id, customer_id, plan_id, status, quantity, seats, current_period_start, current_period_end, trial_start, trial_end, cancel_at_period_end, created_at, updated_at
		 FROM subscriptions WHERE id = $1`, id).Scan(
		&sub.ID, &sub.CustomerID, &sub.PlanID, &sub.Status, &sub.Quantity, &sub.Seats, &sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.TrialStart, &sub.TrialEnd, &sub.CancelAtPeriodEnd, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, err
	}
	p, err := s.plans.GetByID(ctx, sub.PlanID)
	if err == nil {
		sub.Plan = p
	}
	return &sub, nil
}

// UpdateRequest updates mutable subscription fields.
type UpdateRequest struct {
	Quantity *int `json:"quantity,omitempty" validate:"omitempty,min=1"`
	Seats    *int `json:"seats,omitempty" validate:"omitempty,min=1"`
}

// Update modifies a subscription.
func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*models.Subscription, error) {
	sub, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Quantity != nil {
		sub.Quantity = *req.Quantity
	}
	if req.Seats != nil {
		sub.Seats = *req.Seats
	}
	_, err = s.pool.Exec(ctx,
		`UPDATE subscriptions SET quantity=$1, seats=$2, updated_at=NOW() WHERE id=$3`,
		sub.Quantity, sub.Seats, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// Cancel marks a subscription as cancelled.
func (s *Service) Cancel(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	if _, err := s.GetByID(ctx, id); err != nil {
		return nil, err
	}
	_, err := s.pool.Exec(ctx,
		`UPDATE subscriptions SET cancel_at_period_end=true, status='cancelled', updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// Pause sets a subscription status to paused.
func (s *Service) Pause(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	_, err := s.pool.Exec(ctx, `UPDATE subscriptions SET status='paused', updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// Resume sets a subscription status to active.
func (s *Service) Resume(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	_, err := s.pool.Exec(ctx, `UPDATE subscriptions SET status='active', updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// ChangePlanRequest changes a subscription plan.
type ChangePlanRequest struct {
	PlanID uuid.UUID `json:"plan_id" validate:"required"`
}

// ChangePlan switches the plan and returns a proration charge if applicable.
func (s *Service) ChangePlan(ctx context.Context, id uuid.UUID, req ChangePlanRequest) (*models.Subscription, *billing.ProrationCharge, error) {
	sub, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	oldPlan, err := s.plans.GetByID(ctx, sub.PlanID)
	if err != nil {
		return nil, nil, err
	}
	newPlan, err := s.plans.GetByID(ctx, req.PlanID)
	if err != nil {
		return nil, nil, err
	}
	var proration *billing.ProrationCharge
	if sub.Status == models.StatusActive || sub.Status == models.StatusTrialing {
		totalDays := int(sub.CurrentPeriodEnd.Sub(sub.CurrentPeriodStart).Hours() / 24)
		if totalDays <= 0 {
			totalDays = 30
		}
		daysUsed := int(time.Since(sub.CurrentPeriodStart).Hours()/24)
		if daysUsed < 0 {
			daysUsed = 0
		}
		if daysUsed > totalDays {
			daysUsed = totalDays
		}
		amount := billing.ProrationAmount(oldPlan, newPlan, daysUsed, totalDays)
		if amount != 0 {
			desc := "Plan change proration"
			if amount > 0 {
				desc = "Plan upgrade proration"
			} else {
				desc = "Plan downgrade proration"
			}
			proration = &billing.ProrationCharge{Description: desc, AmountCents: amount}
		}
	}
	_, err = s.pool.Exec(ctx,
		`UPDATE subscriptions SET plan_id=$1, updated_at=NOW() WHERE id=$2`,
		req.PlanID, id)
	if err != nil {
		return nil, nil, err
	}
	sub, err = s.GetByID(ctx, id)
	return sub, proration, err
}

func addInterval(t time.Time, interval models.BillingInterval) time.Time {
	if interval == models.BillingYearly {
		return t.AddDate(1, 0, 0)
	}
	return t.AddDate(0, 1, 0)
}
