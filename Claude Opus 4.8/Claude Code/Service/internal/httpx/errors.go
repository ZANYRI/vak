package httpx

import "net/http"

// Error codes used across the API (see service.md "Error Handling").
const (
	CodeValidation             = "VALIDATION_ERROR"
	CodeUnauthorized           = "UNAUTHORIZED"
	CodeForbidden              = "FORBIDDEN"
	CodeNotFound               = "NOT_FOUND"
	CodeConflict               = "CONFLICT"
	CodeIdempotencyConflict    = "IDEMPOTENCY_CONFLICT"
	CodePaymentFailed          = "PAYMENT_FAILED"
	CodeInsufficientPermission = "INSUFFICIENT_PERMISSIONS"
	CodeInternal               = "INTERNAL_ERROR"
)

// APIError is a typed error that carries an HTTP status and a machine code.
type APIError struct {
	Status  int            `json:"-"`
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func (e *APIError) Error() string { return e.Message }

// WithDetails returns a copy of the error with details attached.
func (e *APIError) WithDetails(d map[string]any) *APIError {
	cp := *e
	cp.Details = d
	return &cp
}

func newErr(status int, code, msg string) *APIError {
	return &APIError{Status: status, Code: code, Message: msg}
}

// Constructors for the common error shapes.
func ErrValidation(msg string) *APIError { return newErr(http.StatusBadRequest, CodeValidation, msg) }
func ErrUnauthorized(msg string) *APIError {
	return newErr(http.StatusUnauthorized, CodeUnauthorized, msg)
}
func ErrForbidden(msg string) *APIError { return newErr(http.StatusForbidden, CodeForbidden, msg) }
func ErrNotFound(msg string) *APIError  { return newErr(http.StatusNotFound, CodeNotFound, msg) }
func ErrConflict(msg string) *APIError  { return newErr(http.StatusConflict, CodeConflict, msg) }
func ErrIdempotencyConflict(msg string) *APIError {
	return newErr(http.StatusConflict, CodeIdempotencyConflict, msg)
}
func ErrPaymentFailed(msg string) *APIError {
	return newErr(http.StatusPaymentRequired, CodePaymentFailed, msg)
}
func ErrInsufficientPermission(msg string) *APIError {
	return newErr(http.StatusForbidden, CodeInsufficientPermission, msg)
}
func ErrInternal(msg string) *APIError {
	if msg == "" {
		msg = "internal server error"
	}
	return newErr(http.StatusInternalServerError, CodeInternal, msg)
}
