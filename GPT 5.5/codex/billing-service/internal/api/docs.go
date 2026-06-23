package api

import (
	"net/http"
	"os"
)

func (a *App) docs(w http.ResponseWriter, r *http.Request) {
	b, e := os.ReadFile("docs/openapi.yaml")
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "OpenAPI document is unavailable", nil)
		return
	}
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	_, _ = w.Write(b)
}
