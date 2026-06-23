package models

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents an RBAC role.
type UserRole string

const (
	RoleAdmin          UserRole = "admin"
	RoleBillingManager UserRole = "billing_manager"
	RoleSupport        UserRole = "support"
	RoleCustomer       UserRole = "customer"
)

// User is the authentication entity.
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SubscriptionStatus represents the lifecycle state of a subscription.
type SubscriptionStatus string

const (
	StatusTrialing   SubscriptionStatus = "trialing"
	StatusActive     SubscriptionStatus = "active"
	StatusPastDue    SubscriptionStatus = "past_due"
	StatusPaused     SubscriptionStatus = "paused"
	StatusCancelled  SubscriptionStatus = "cancelled"
	StatusExpired    SubscriptionStatus = "expired"
)

// BillingInterval defines plan billing cadence.
type BillingInterval string

const (
	BillingMonthly BillingInterval = "monthly"
	BillingYearly  BillingInterval = "yearly"
)

// PlanModel defines the pricing model of a plan.
type PlanModel string

const (
	PlanFlat        PlanModel = "flat"
	PlanPerSeat     PlanModel = "per_seat"
	PlanUsageBased  PlanModel = "usage_based"
	PlanTiered      PlanModel = "tiered"
	PlanHybrid      PlanModel = "hybrid"
)

// InvoiceStatus represents the invoice state.
type InvoiceStatus string

const (
	InvoiceDraft         InvoiceStatus = "draft"
	InvoiceOpen          InvoiceStatus = "open"
	InvoicePaid          InvoiceStatus = "paid"
	InvoiceVoid          InvoiceStatus = "void"
	InvoiceUncollectible InvoiceStatus = "uncollectible"
)

// PaymentStatus represents the payment attempt result.
type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentSucceeded PaymentStatus = "succeeded"
	PaymentFailed   PaymentStatus = "failed"
	PaymentRefunded PaymentStatus = "refunded"
)

// JobStatus represents the state of a background job.
type JobStatus string

const (
	JobQueued    JobStatus = "queued"
	JobRunning   JobStatus = "running"
	JobCompleted JobStatus = "completed"
	JobFailed    JobStatus = "failed"
	JobRetrying  JobStatus = "retrying"
	JobDead      JobStatus = "dead"
)

// CouponType represents the discount type of a coupon.
type CouponType string

const (
	CouponPercentage  CouponType = "percentage"
	CouponFixedAmount CouponType = "fixed_amount"
)

// Permission describes an action allowed for a role.
type Permission string

const (
	PermManageAll      Permission = "manage_all"
	PermManagePlans    Permission = "manage_plans"
	PermManageInvoices Permission = "manage_invoices"
	PermManageSubscriptions Permission = "manage_subscriptions"
	PermManageCoupons  Permission = "manage_coupons"
	PermManageTaxRules Permission = "manage_tax_rules"
	PermViewCustomers  Permission = "view_customers"
	PermViewOwnData    Permission = "view_own_data"
)

// RolePermissions maps each role to its allowed permissions.
var RolePermissions = map[UserRole][]Permission{
	RoleAdmin:          {PermManageAll},
	RoleBillingManager: {PermManagePlans, PermManageInvoices, PermManageSubscriptions, PermManageCoupons, PermManageTaxRules, PermViewCustomers},
	RoleSupport:        {PermViewCustomers, PermViewOwnData},
	RoleCustomer:       {PermViewOwnData},
}

// HasPermission checks whether a role has a given permission.
func HasPermission(role UserRole, perm Permission) bool {
	for _, p := range RolePermissions[role] {
		if p == perm || p == PermManageAll {
			return true
		}
	}
	return false
}
