package model

import "time"

type MonthRecap struct {
	ID               string    `json:"id"`
	UserID           string    `json:"userId"`
	Year             int       `json:"year"`
	Month            int       `json:"month"`
	WeeksIncluded    int       `json:"weeksIncluded"`
	TotalWorkedMin   int       `json:"totalWorkedMin"`
	TotalExpectedMin int       `json:"totalExpectedMin"`
	TotalOvertimeMin int       `json:"totalOvertimeMin"`
	AverageStartISO  string    `json:"averageStartISO"`
	AverageStopISO   string    `json:"averageStopISO"`
	GeneratedAt      time.Time `json:"generatedAt"`
}
