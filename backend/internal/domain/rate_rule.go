package domain

import (
	"time"
)

type RateRule struct {
	ID          uint      `json:"id"`
	RuleName    string    `json:"rulename"`
	DayOfWeek   int       `json:"dayofweek"`
	StartTime   time.Time `json:"starttime"`
	EndTime     time.Time `json:"endtime"`
	HourlyRate  float64   `json:"hourlyrate"`
	IsOverNight int       `json:"isovernight"`
	Priority    int       `json:"priority"`
	Enabled     int       `json:"enabled"`
	CreatedAt   time.Time `json:"createdat"`
}
