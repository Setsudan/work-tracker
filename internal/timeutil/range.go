package timeutil

import "time"

func StartOfWeekMonday(t time.Time, loc *time.Location) time.Time {
	d := t.In(loc)
	offset := (int(d.Weekday()) + 6) % 7 // Monday=0
	start := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
	return start.AddDate(0, 0, -offset)
}

func DaysOfWeek(startMonday time.Time, loc *time.Location) []time.Time {
	res := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		res[i] = startMonday.AddDate(0, 0, i)
	}
	return res
}

func MonthStart(t time.Time, loc *time.Location) time.Time {
	d := t.In(loc)
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, loc)
}

func MonthEnd(t time.Time, loc *time.Location) time.Time {
	start := MonthStart(t, loc)
	return start.AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func PreviousMonthStart(t time.Time, loc *time.Location) time.Time {
	start := MonthStart(t, loc)
	return start.AddDate(0, -1, 0)
}

func PreviousMonthEnd(t time.Time, loc *time.Location) time.Time {
	return MonthStart(t, loc).Add(-time.Nanosecond)
}
