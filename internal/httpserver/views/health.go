package views

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	writeJSON(w, http.StatusOK, response)
}

func RegisterHealthRoutes(r chi.Router) {
	r.Get("/health", HealthHandler)
}
