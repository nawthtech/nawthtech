package handlers

import (
	"github.com/go-chi/chi/v5"
)

type Services struct {
	// سيتم إضافة الخدمات لاحقاً
}

func Register(r chi.Router, services *Services) {
	// تنفيذ أساسي للمسارات
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}