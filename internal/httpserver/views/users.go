package views

import (
	"encoding/json"
	"net/http"

	"work-tracker/internal/config"
	"work-tracker/internal/httpserver/middleware"
	"work-tracker/internal/store"

	"github.com/go-chi/chi/v5"
)

type updateMeReq struct {
	FullName                 *string  `json:"fullName"`
	WeeklyHours              *float64 `json:"weeklyHours"`
	DefaultLunchBreakMinutes *int     `json:"defaultLunchBreakMinutes"`
	Timezone                 *string  `json:"timezone"`
}

func Users(cfg config.Config, stores store.Stores) chi.Router {
	r := chi.NewRouter()

	r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := middleware.UserIDFromContext(r.Context())
		u, err := stores.Users.GetByID(r.Context(), uid)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		u.PasswordHash = ""
		writeJSON(w, http.StatusOK, u)
	})

	r.Put("/me", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := middleware.UserIDFromContext(r.Context())
		u, err := stores.Users.GetByID(r.Context(), uid)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		var req updateMeReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if req.FullName != nil {
			u.FullName = *req.FullName
		}
		if req.WeeklyHours != nil {
			u.WeeklyHours = *req.WeeklyHours
		}
		if req.DefaultLunchBreakMinutes != nil {
			u.DefaultLunchBreakMinutes = *req.DefaultLunchBreakMinutes
		}
		if req.Timezone != nil && *req.Timezone != "" {
			u.Timezone = *req.Timezone
		}
		if err := stores.Users.Update(r.Context(), u); err != nil {
			http.Error(w, "failed to update", http.StatusInternalServerError)
			return
		}
		u.PasswordHash = ""
		writeJSON(w, http.StatusOK, u)
	})

	return r
}
