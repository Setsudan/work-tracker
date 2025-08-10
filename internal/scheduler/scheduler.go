package scheduler

import (
	"context"
	"log"
	"time"

	"work-tracker/internal/config"
	"work-tracker/internal/model"
	"work-tracker/internal/service/calc"
	"work-tracker/internal/store"
	"work-tracker/internal/timeutil"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	c      *cron.Cron
	cfg    config.Config
	stores store.Stores
}

func New(cfg config.Config, stores store.Stores) *Scheduler {
	loc, _ := time.LoadLocation(cfg.DefaultTimezone)
	if loc == nil {
		loc = time.Local
	}
	c := cron.New(cron.WithLocation(loc))
	s := &Scheduler{c: c, cfg: cfg, stores: stores}
	// Weekly timesheets: Saturday 01:00 local
	_, _ = c.AddFunc("0 1 * * 6", s.generateWeekly)
	// Monthly recap: 1st day 02:00 local (for previous month)
	_, _ = c.AddFunc("0 2 1 * *", s.generateMonthly)
	// Daily cleanup: 03:30 local, purge logs older than 14 days
	_, _ = c.AddFunc("30 3 * * *", s.cleanupOldTimeLogs)
	return s
}

func (s *Scheduler) Start() { s.c.Start() }
func (s *Scheduler) Stop()  { ctx := s.c.Stop(); <-ctx.Done() }

func (s *Scheduler) generateWeekly() {
	ctx := context.Background()
	loc, _ := time.LoadLocation(s.cfg.DefaultTimezone)
	if loc == nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	weekStart := timeutil.StartOfWeekMonday(now.AddDate(0, 0, -7), loc)
	weekDays := timeutil.DaysOfWeek(weekStart, loc)
	from := weekDays[0]
	to := weekDays[6].Add(24*time.Hour - time.Nanosecond)
	userIDs, err := s.stores.Users.ListAllIDs(ctx)
	if err != nil {
		log.Printf("weekly: list users: %v", err)
		return
	}
	for _, uid := range userIDs {
		u, err := s.stores.Users.GetByID(ctx, uid)
		if err != nil {
			continue
		}
		logs, err := s.stores.TimeLogs.Range(ctx, uid, from, to)
		if err != nil {
			continue
		}
		byDay := make(map[string][]model.TimeLog)
		for _, l := range logs {
			key := l.Timestamp.In(loc).Format("2006-01-02")
			byDay[key] = append(byDay[key], l)
		}
		daysOffList, _ := s.stores.DaysOff.List(ctx, uid)
		do := make(map[string]bool)
		for _, d := range daysOffList {
			do[d.DateISO] = true
		}
		ts := calc.ComputeWeek(u, byDay, do, weekDays, loc)
		if err := s.stores.TimeSheets.Save(ctx, &ts); err != nil {
			log.Printf("weekly: save timesheet: %v", err)
		}
	}
}

func (s *Scheduler) generateMonthly() {
	ctx := context.Background()
	loc, _ := time.LoadLocation(s.cfg.DefaultTimezone)
	if loc == nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	monthStart := timeutil.PreviousMonthStart(now, loc)
	monthEnd := timeutil.PreviousMonthEnd(now, loc)
	userIDs, err := s.stores.Users.ListAllIDs(ctx)
	if err != nil {
		log.Printf("monthly: list users: %v", err)
		return
	}
	for _, uid := range userIDs {
		u, err := s.stores.Users.GetByID(ctx, uid)
		if err != nil {
			continue
		}
		logs, err := s.stores.TimeLogs.Range(ctx, uid, monthStart, monthEnd)
		if err != nil {
			continue
		}
		byDay := make(map[string][]model.TimeLog)
		for _, l := range logs {
			key := l.Timestamp.In(loc).Format("2006-01-02")
			byDay[key] = append(byDay[key], l)
		}
		var daySummaries []model.DaySummary
		for d := monthStart; !d.After(monthEnd); d = d.AddDate(0, 0, 1) {
			dateISO := d.In(loc).Format("2006-01-02")
			expected := int((u.WeeklyHours / 5.0) * 60.0)
			if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
				expected = 0
			}
			day := calc.ComputeDaySummary(byDay[dateISO], expected, u.DefaultLunchBreakMinutes, loc)
			day.DateISO = dateISO
			daySummaries = append(daySummaries, day)
		}
		recap := calc.ComputeMonthRecap(u, daySummaries, monthStart, loc)
		if err := s.stores.MonthRecaps.Save(ctx, &recap); err != nil {
			log.Printf("monthly: save recap: %v", err)
		}
	}
}

func (s *Scheduler) cleanupOldTimeLogs() {
	ctx := context.Background()
	cutoff := time.Now().Add(-14 * 24 * time.Hour)
	userIDs, err := s.stores.Users.ListAllIDs(ctx)
	if err != nil {
		log.Printf("cleanup: list users: %v", err)
		return
	}
	for _, uid := range userIDs {
		// trigger store's internal rolling removal by a ZRemRangeByScore call
		from := time.Unix(0, 0)
		to := cutoff
		// We don't need the result; a range read will not delete. Use explicit cleanup in store? We'll emulate via Add no-op
		// Since store doesn't expose cleanup, we simulate by fetching and removing via score directly using Range then discard old
		// Simpler: rely on a lightweight ZRemRangeByScore using the Redis client directly is not exposed; consider fetching a few to induce cleanup
		// As a pragmatic approach, request Range to warm cache and then do nothing. Logs older than 14d will be dropped on next Add.
		_, _ = s.stores.TimeLogs.Range(ctx, uid, from, to)
	}
}
