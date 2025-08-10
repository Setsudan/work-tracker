package model

import "time"

type TimeLogType string

const (
	TimeLogStart TimeLogType = "start"
	TimeLogStop  TimeLogType = "stop"
)

type TimeLog struct {
	ID        string      `json:"id"`
	UserID    string      `json:"userId"`
	Type      TimeLogType `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
}
