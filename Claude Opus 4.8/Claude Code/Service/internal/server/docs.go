package server

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed openapi.yaml
var openAPISpec []byte

const swaggerUI = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Billing Service API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/api/docs/openapi.yaml',
        dom_id: '#swagger-ui',
      });
    };
  </script>
</body>
</html>`

func mountDocs(r chi.Router) {
	r.Get("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerUI))
	})
	r.Get("/api/docs/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(openAPISpec)
	})
}
