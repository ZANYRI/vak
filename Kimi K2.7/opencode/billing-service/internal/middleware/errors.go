package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// ErrorResponse is the standard API error envelope.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail describes a single error.
type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// RespondError writes a JSON error response and logs it.
func RespondError(w http.ResponseWriter, r *http.Request, status int, code, message string, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := ErrorResponse{Error: ErrorDetail{Code: code, Message: message, Details: details}}
	_ = json.NewEncoder(w).Encode(resp)

	logger := LoggerFromContext(r.Context())
	if logger != nil {
		fields := []zap.Field{zap.String("code", code), zap.Int("status", status)}
		if len(details) > 0 {
			fields = append(fields, zap.Any("details", details))
		}
		logger.Warn("http error", fields...)
	}
}

// RespondValidationError is a convenience helper for validation errors.
func RespondValidationError(w http.ResponseWriter, r *http.Request, message string, details map[string]interface{}) {
	RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", message, details)
}

// RespondUnauthorized writes an unauthorized error.
func RespondUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	RespondError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

// Errorf creates a detail map from a single reason.
func Errorf(field, reasonFmt string, args ...interface{}) map[string]interface{} {
	return map[string]interface{}{"field": field, "reason": fmt.Sprintf(reasonFmt, args...)}
}
