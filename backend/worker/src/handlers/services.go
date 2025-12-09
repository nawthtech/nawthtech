package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"nawthtech/worker/src/utils"
)

type Service struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func GetServicesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	db, _ := utils.GetD1().GetDB()
	rows, err := db.QueryContext(ctx, "SELECT id, title, description, price FROM services LIMIT 50")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	services := []Service{}
	for rows.Next() {
		s := Service{}
		rows.Scan(&s.ID, &s.Title, &s.Description, &s.Price)
		services = append(services, s)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    services,
	})
}

func GetServiceByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	db, _ := utils.GetD1().GetDB()
	id := r.URL.Path[len("/services/"):]

	row := db.QueryRowContext(ctx, "SELECT id, title, description, price FROM services WHERE id = ?", id)
	service := Service{}
	err := row.Scan(&service.ID, &service.Title, &service.Description, &service.Price)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "SERVICE_NOT_FOUND",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    service,
	})
}