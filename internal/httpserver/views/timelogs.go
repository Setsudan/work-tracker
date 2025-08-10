package views

import (
	"net/http"
	"time"

	"work-tracker/internal/config"
	"work-tracker/internal/httpserver/middleware"
	"work-tracker/internal/model"
	"work-tracker/internal/store"

	"github.com/go-chi/chi/v5"
)

type toggleResp struct {
	Created model.TimeLog `json:"created"`
}

func TimeLogs(cfg config.Config, stores store.Stores) chi.Router {
	r := chi.NewRouter()

	r.Post("/toggle", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := middleware.UserIDFromContext(r.Context())
		last, _ := stores.TimeLogs.GetLast(r.Context(), uid)
		now := time.Now()
		nextType := model.TimeLogStart
		if last != nil && last.Type == model.TimeLogStart {
			nextType = model.TimeLogStop
		}
		created, err := stores.TimeLogs.Add(r.Context(), uid, nextType, now)
		if err != nil {
			http.Error(w, "failed to create log", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, toggleResp{Created: *created})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := middleware.UserIDFromContext(r.Context())
		q := r.URL.Query()
		fromStr := q.Get("from")
		toStr := q.Get("to")
		if fromStr == "" || toStr == "" {
			http.Error(w, "from and to required (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		from, err1 := time.Parse("2006-01-02", fromStr)
		to, err2 := time.Parse("2006-01-02", toStr)
		if err1 != nil || err2 != nil {
			http.Error(w, "invalid date", http.StatusBadRequest)
			return
		}
		to = to.Add(24*time.Hour - time.Nanosecond)
		logs, err := stores.TimeLogs.Range(r.Context(), uid, from, to)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, logs)
	})

	return r
}
