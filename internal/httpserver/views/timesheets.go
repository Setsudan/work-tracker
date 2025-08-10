package views

import (
	"net/http"
	"time"

	"work-tracker/internal/config"
	"work-tracker/internal/httpserver/middleware"
	"work-tracker/internal/model"
	"work-tracker/internal/service/calc"
	"work-tracker/internal/store"
	"work-tracker/internal/timeutil"

	"github.com/go-chi/chi/v5"
)

func TimeSheets(cfg config.Config, stores store.Stores) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := middleware.UserIDFromContext(r.Context())
		u, err := stores.Users.GetByID(r.Context(), uid)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		loc, _ := time.LoadLocation(u.Timezone)
		if loc == nil {
			loc, _ = time.LoadLocation(cfg.DefaultTimezone)
		}
		q := r.URL.Query()
		weekStartStr := q.Get("weekStart")
		var weekStart time.Time
		if weekStartStr == "" {
			weekStart = timeutil.StartOfWeekMonday(time.Now(), loc)
		} else {
			var err error
			weekStart, err = time.ParseInLocation("2006-01-02", weekStartStr, loc)
			if err != nil {
				http.Error(w, "invalid weekStart", http.StatusBadRequest)
				return
			}
		}
		weekDays := timeutil.DaysOfWeek(weekStart, loc)
		from := weekDays[0]
		to := weekDays[6].Add(24*time.Hour - time.Nanosecond)
		logs, err := stores.TimeLogs.Range(r.Context(), uid, from, to)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		// group logs by day in user's timezone
		byDay := make(map[string][]model.TimeLog)
		for _, l := range logs {
			date := l.Timestamp.In(loc).Format("2006-01-02")
			byDay[date] = append(byDay[date], l)
		}
		// days off map
		daysOffList, _ := stores.DaysOff.List(r.Context(), uid)
		daysOff := make(map[string]bool)
		for _, d := range daysOffList {
			daysOff[d.DateISO] = true
		}
		ts := calc.ComputeWeek(u, byDay, daysOff, weekDays, loc)
		_ = stores.TimeSheets.Save(r.Context(), &ts)
		writeJSON(w, http.StatusOK, ts)
	})

	return r
}
