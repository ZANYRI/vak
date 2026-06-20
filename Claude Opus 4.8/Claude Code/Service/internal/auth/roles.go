package auth

// Role constants for RBAC.
const (
	RoleAdmin          = "admin"
	RoleBillingManager = "billing_manager"
	RoleSupport        = "support"
	RoleCustomer       = "customer"
)

// ValidRoles is the set of assignable roles.
var ValidRoles = map[string]bool{
	RoleAdmin:          true,
	RoleBillingManager: true,
	RoleSupport:        true,
	RoleCustomer:       true,
}

// CanManageBilling reports whether the role may create/modify billing rules
// (plans, coupons, tax rules, subscriptions, invoices).
func CanManageBilling(role string) bool {
	return role == RoleAdmin || role == RoleBillingManager
}

// CanViewBilling reports whether the role may read billing data broadly.
// (support can view; customer is restricted to its own resources at the handler.)
func CanViewBilling(role string) bool {
	return role == RoleAdmin || role == RoleBillingManager || role == RoleSupport
}
