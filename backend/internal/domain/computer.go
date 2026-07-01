package domain

import (
	"time"
)

type Computer struct {
	ID               uint      `json:"id"`
	MachineNumber    string    `json:"machinenumber"`
	Area             string    `json:"area"`
	IpAddress        string    `json:"ipadress"`
	Status           int       `json:"status"`
	HourlyRate       float64   `json:"hourlyrate"`
	CurrentSessionID int       `json:"currentsessionid"`
	LastOnlineAt     time.Time `json:"lastonlineat"`
	CreatedAt        time.Time `json:"createdat"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
