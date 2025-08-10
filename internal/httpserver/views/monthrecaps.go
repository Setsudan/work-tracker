package views

import (
	"net/http"
	"strconv"
	"time"

	"work-tracker/internal/config"
	"work-tracker/internal/httpserver/middleware"
	"work-tracker/internal/model"
	"work-tracker/internal/service/calc"
	"work-tracker/internal/store"
	"work-tracker/internal/timeutil"

	"github.com/go-chi/chi/v5"
)

func MonthRecaps(cfg config.Config, stores store.Stores) chi.Router {
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
		year, _ := strconv.Atoi(q.Get("year"))
		month, _ := strconv.Atoi(q.Get("month"))
		var monthStart time.Time
		if year == 0 || month == 0 {
			monthStart = timeutil.MonthStart(time.Now(), loc)
		} else {
			monthStart = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
		}
		monthEnd := timeutil.MonthEnd(monthStart, loc)
		logs, err := stores.TimeLogs.Range(r.Context(), uid, monthStart, monthEnd)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		// group by day
		byDay := make(map[string][]model.TimeLog)
		for _, l := range logs {
			date := l.Timestamp.In(loc).Format("2006-01-02")
			byDay[date] = append(byDay[date], l)
		}
		// days off
		daysOffList, _ := stores.DaysOff.List(r.Context(), uid)
		daysOff := make(map[string]bool)
		for _, d := range daysOffList {
			daysOff[d.DateISO] = true
		}
		// build day summaries for all days in month
		var daySummaries []model.DaySummary
		for d := monthStart; !d.After(monthEnd); d = d.AddDate(0, 0, 1) {
			dateISO := d.In(loc).Format("2006-01-02")
			expected := int((u.WeeklyHours / 5.0) * 60.0)
			if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
				expected = 0
			}
			if daysOff[dateISO] {
				daySummaries = append(daySummaries, model.DaySummary{DateISO: dateISO, WorkedMinutes: expected, ExpectedMinutes: expected, OvertimeMinutes: 0})
				continue
			}
			day := calc.ComputeDaySummary(byDay[dateISO], expected, u.DefaultLunchBreakMinutes, loc)
			day.DateISO = dateISO
			daySummaries = append(daySummaries, day)
		}
		recap := calc.ComputeMonthRecap(u, daySummaries, monthStart, loc)
		_ = stores.MonthRecaps.Save(r.Context(), &recap)
		writeJSON(w, http.StatusOK, recap)
	})

	return r
}
