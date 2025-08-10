package views

import (
	"encoding/json"
	"net/http"

	"work-tracker/internal/config"
	"work-tracker/internal/httpserver/middleware"
	"work-tracker/internal/model"
	"work-tracker/internal/store"

	"github.com/go-chi/chi/v5"
)

type dayOffReq struct {
	DateISO string `json:"dateISO"`
	Reason  string `json:"reason"`
}

func DaysOff(cfg config.Config, stores store.Stores) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := middleware.UserIDFromContext(r.Context())
		d, err := stores.DaysOff.List(r.Context(), uid)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, d)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := middleware.UserIDFromContext(r.Context())
		var req dayOffReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if req.DateISO == "" {
			http.Error(w, "dateISO required", http.StatusBadRequest)
			return
		}
		if err := stores.DaysOff.Add(r.Context(), uid, model.DayOff{DateISO: req.DateISO, Reason: req.Reason}); err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	r.Delete("/{date}", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := middleware.UserIDFromContext(r.Context())
		date := chi.URLParam(r, "date")
		if date == "" {
			http.Error(w, "date required", http.StatusBadRequest)
			return
		}
		if err := stores.DaysOff.Remove(r.Context(), uid, date); err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return r
}
