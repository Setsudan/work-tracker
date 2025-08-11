package model

import "time"

type User struct {
	ID                       string    `json:"id"`
	Email                    string    `json:"email"`
	PasswordHash             string    `json:"passwordHash"`
	FullName                 string    `json:"fullName"`
	WeeklyHours              float64   `json:"weeklyHours"`
	DefaultLunchBreakMinutes int       `json:"defaultLunchBreakMinutes"`
	Timezone                 string    `json:"timezone"`
	CreatedAt                time.Time `json:"createdAt"`
	UpdatedAt                time.Time `json:"updatedAt"`
}
