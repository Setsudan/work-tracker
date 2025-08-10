package calc

import (
	"sort"
	"time"

	"work-tracker/internal/model"
)

func ComputeDaySummary(logs []model.TimeLog, expectedMinutes int, defaultLunchMinutes int, loc *time.Location) model.DaySummary {
	if len(logs) == 0 {
		return model.DaySummary{
			WorkedMinutes:   0,
			ExpectedMinutes: expectedMinutes,
			OvertimeMinutes: 0 - expectedMinutes,
			LogsCount:       0,
		}
	}
	// ensure chronological order
	sort.Slice(logs, func(i, j int) bool { return logs[i].Timestamp.Before(logs[j].Timestamp) })

	var startTimes []time.Time
	var stopTimes []time.Time
	for _, l := range logs {
		if l.Type == model.TimeLogStart {
			startTimes = append(startTimes, l.Timestamp.In(loc))
		}
		if l.Type == model.TimeLogStop {
			stopTimes = append(stopTimes, l.Timestamp.In(loc))
		}
	}

	// pair starts and stops sequentially
	pairs := min(len(startTimes), len(stopTimes))
	worked := 0
	for i := 0; i < pairs; i++ {
		d := stopTimes[i].Sub(startTimes[i])
		if d > 0 {
			worked += int(d.Minutes())
		}
	}

	// lunch rule: only if exactly one start and one stop
	if len(startTimes) == 1 && len(stopTimes) == 1 {
		worked -= defaultLunchMinutes
		if worked < 0 {
			worked = 0
		}
	}

	var firstStart, lastStop, avgStart, avgStop string
	if len(startTimes) > 0 {
		firstStart = startTimes[0].Format(time.RFC3339)
	}
	if len(stopTimes) > 0 {
		lastStop = stopTimes[len(stopTimes)-1].Format(time.RFC3339)
	}
	if len(startTimes) > 0 {
		avgStart = averageTimes(startTimes).Format(time.RFC3339)
	}
	if len(stopTimes) > 0 {
		avgStop = averageTimes(stopTimes).Format(time.RFC3339)
	}

	overtime := worked - expectedMinutes
	return model.DaySummary{
		WorkedMinutes:   worked,
		ExpectedMinutes: expectedMinutes,
		OvertimeMinutes: overtime,
		FirstStartISO:   firstStart,
		LastStopISO:     lastStop,
		AverageStartISO: avgStart,
		AverageStopISO:  avgStop,
		LogsCount:       len(logs),
	}
}

func ComputeWeek(user *model.User, daysLogs map[string][]model.TimeLog, daysOff map[string]bool, weekDays []time.Time, loc *time.Location) model.TimeSheet {
	expectedPerWorkday := int((user.WeeklyHours / 5.0) * 60.0)
	var days []model.DaySummary
	totalWorked := 0
	totalExpected := 0
	for _, day := range weekDays {
		dateISO := day.In(loc).Format("2006-01-02")
		logs := daysLogs[dateISO]
		expected := expectedPerWorkday
		if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
			expected = 0
		}
		if daysOff[dateISO] {
			// count as fulfilled expected minutes
			days = append(days, model.DaySummary{
				DateISO:         dateISO,
				WorkedMinutes:   expected,
				ExpectedMinutes: expected,
				OvertimeMinutes: 0,
				LogsCount:       len(logs),
			})
			totalWorked += expected
			totalExpected += expected
			continue
		}
		d := ComputeDaySummary(logs, expected, user.DefaultLunchBreakMinutes, loc)
		d.DateISO = dateISO
		days = append(days, d)
		totalWorked += d.WorkedMinutes
		totalExpected += expected
	}
	weekStart := weekDays[0].In(loc)
	year, week := weekStart.ISOWeek()
	return model.TimeSheet{
		UserID:           user.ID,
		WeekStartISO:     weekStart.Format("2006-01-02"),
		WeekNumber:       week,
		Year:             year,
		Days:             days,
		TotalWorkedMin:   totalWorked,
		TotalExpectedMin: totalExpected,
		TotalOvertimeMin: totalWorked - totalExpected,
		GeneratedAt:      time.Now(),
	}
}

func ComputeMonthRecap(user *model.User, daySummaries []model.DaySummary, monthStart time.Time, loc *time.Location) model.MonthRecap {
	totalWorked := 0
	totalExpected := 0
	var avgStartTimes []time.Time
	var avgStopTimes []time.Time
	for _, d := range daySummaries {
		totalWorked += d.WorkedMinutes
		totalExpected += d.ExpectedMinutes
		if d.AverageStartISO != "" {
			if t, err := time.Parse(time.RFC3339, d.AverageStartISO); err == nil {
				avgStartTimes = append(avgStartTimes, t)
			}
		}
		if d.AverageStopISO != "" {
			if t, err := time.Parse(time.RFC3339, d.AverageStopISO); err == nil {
				avgStopTimes = append(avgStopTimes, t)
			}
		}
	}
	avgStart := ""
	avgStop := ""
	if len(avgStartTimes) > 0 {
		avgStart = averageTimes(avgStartTimes).Format(time.RFC3339)
	}
	if len(avgStopTimes) > 0 {
		avgStop = averageTimes(avgStopTimes).Format(time.RFC3339)
	}
	return model.MonthRecap{
		UserID:           user.ID,
		Year:             monthStart.In(loc).Year(),
		Month:            int(monthStart.In(loc).Month()),
		WeeksIncluded:    0, // computed elsewhere if needed
		TotalWorkedMin:   totalWorked,
		TotalExpectedMin: totalExpected,
		TotalOvertimeMin: totalWorked - totalExpected,
		AverageStartISO:  avgStart,
		AverageStopISO:   avgStop,
		GeneratedAt:      time.Now(),
	}
}

func averageTimes(ts []time.Time) time.Time {
	if len(ts) == 0 {
		return time.Time{}
	}
	sum := int64(0)
	for _, t := range ts {
		sum += t.Unix()
	}
	avg := sum / int64(len(ts))
	return time.Unix(avg, 0).UTC()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
