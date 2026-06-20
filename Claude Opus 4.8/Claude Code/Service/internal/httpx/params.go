package httpx

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// PathUUID parses a UUID path parameter, returning an *APIError on failure.
func PathUUID(r *http.Request, name string) (uuid.UUID, error) {
	raw := chi.URLParam(r, name)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, ErrValidation("invalid " + name + ": must be a UUID")
	}
	return id, nil
}

// Pagination reads ?limit= and ?offset= with sane defaults and caps.
func Pagination(r *http.Request) (limit, offset int) {
	limit, offset = 50, 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
}
