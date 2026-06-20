package httpx

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// envelope wraps API errors per the documented format: {"error": {...}}.
type errEnvelope struct {
	Error *APIError `json:"error"`
}

// JSON writes a success payload with the given status code.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// Error writes an error response. Non-APIError values are masked as INTERNAL_ERROR.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		apiErr = ErrInternal("")
		// Log the real cause but never leak it to the client.
		slog.ErrorContext(r.Context(), "unhandled error", "error", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.Status)
	_ = json.NewEncoder(w).Encode(errEnvelope{Error: apiErr})
}
