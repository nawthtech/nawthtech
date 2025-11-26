package handlers

import (
	"github.com/nawthtech/nawthtech/backend/internal/handlers/health"
  "github.com/nawthtech/nawthtech/backend/internal/handlers/sse"

	"github.com/go-chi/chi/v5"
)

func Register(r *chi.Mux) {
	r.Get("/health", health.Handler)

	r.Get("/sse", sse.Handler)
}
