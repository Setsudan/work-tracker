package model

type DayOff struct {
	DateISO string `json:"dateISO"`
	Reason  string `json:"reason"`
}
