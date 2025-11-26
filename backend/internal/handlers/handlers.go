package handlers

import (
	"backend-app/internal/handlers/health"
  "backend-app/internal/handlers/sse"

	"github.com/go-chi/chi/v5"
)

func Register(r *chi.Mux) {
	r.Get("/health", health.Handler)

	r.Get("/sse", sse.Handler)
}
