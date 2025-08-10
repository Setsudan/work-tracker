package model

import "time"

type DaySummary struct {
	DateISO         string `json:"dateISO"`
	WorkedMinutes   int    `json:"workedMinutes"`
	ExpectedMinutes int    `json:"expectedMinutes"`
	OvertimeMinutes int    `json:"overtimeMinutes"`
	FirstStartISO   string `json:"firstStartISO"`
	LastStopISO     string `json:"lastStopISO"`
	AverageStartISO string `json:"averageStartISO"`
	AverageStopISO  string `json:"averageStopISO"`
	LogsCount       int    `json:"logsCount"`
}

type TimeSheet struct {
	ID               string       `json:"id"`
	UserID           string       `json:"userId"`
	WeekStartISO     string       `json:"weekStartISO"`
	WeekNumber       int          `json:"weekNumber"`
	Year             int          `json:"year"`
	Days             []DaySummary `json:"days"`
	TotalWorkedMin   int          `json:"totalWorkedMin"`
	TotalExpectedMin int          `json:"totalExpectedMin"`
	TotalOvertimeMin int          `json:"totalOvertimeMin"`
	GeneratedAt      time.Time    `json:"generatedAt"`
}
