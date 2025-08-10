package routes

import (
	"work-tracker/internal/config"
	"work-tracker/internal/httpserver/middleware"
	"work-tracker/internal/httpserver/views"
	"work-tracker/internal/store"

	"github.com/go-chi/chi/v5"
)

func API(cfg config.Config, stores store.Stores) chi.Router {
	r := chi.NewRouter()

	// Health check endpoint (no auth required)
	views.RegisterHealthRoutes(r)

	r.Mount("/auth", views.Auth(cfg, stores))

	r.Group(func(pr chi.Router) {
		pr.Use(middleware.WithAuth(cfg, stores))
		pr.Mount("/users", views.Users(cfg, stores))
		pr.Mount("/time-logs", views.TimeLogs(cfg, stores))
		pr.Mount("/days-off", views.DaysOff(cfg, stores))
		pr.Mount("/timesheets", views.TimeSheets(cfg, stores))
		pr.Mount("/month-recaps", views.MonthRecaps(cfg, stores))
	})

	return r
}
