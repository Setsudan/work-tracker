package httpserver

import (
	"net/http"
	"strings"

	"work-tracker/internal/config"
	"work-tracker/internal/httpserver/routes"
	"work-tracker/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewServer(cfg config.Config, stores store.Stores) http.Handler {
	r := chi.NewRouter()

	allowed := strings.Split(cfg.AllowOriginsCSV, ",")
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowed,
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
	}))

	r.Mount("/v1", routes.API(cfg, stores))
	return r
}
