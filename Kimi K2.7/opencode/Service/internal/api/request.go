package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"billing-service/internal/middleware"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// DecodeJSON decodes the request body into dst and validates it.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request payload", map[string]interface{}{"reason": err.Error()})
		return false
	}
	if err := validate.Struct(dst); err != nil {
		fields := make(map[string]interface{})
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			for _, fe := range verr {
				fields[fe.Field()] = fmt.Sprintf("failed %s", fe.Tag())
			}
		} else {
			fields["error"] = err.Error()
		}
		middleware.RespondValidationError(w, r, "validation failed", fields)
		return false
	}
	return true
}

// RespondJSON encodes a value as JSON and writes it with the given status.
func RespondJSON(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// IdempotencyKey returns the Idempotency-Key header value trimmed.
func IdempotencyKey(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get("Idempotency-Key"))
}
