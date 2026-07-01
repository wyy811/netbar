package domain

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID         uint            `json:"id"`
	UserID     int             `json:"userid"`
	Action     string          `json:"action"`
	TargetType string          `json:"targettype"`
	TargetID   int             `json:"targetid"`
	Detail     json.RawMessage `json:"detail"`
	IpAdress   string          `json:"ipadress"`
	UserAgent  string          `json:"useragent"`
	CreatedAt  time.Time       `json:"createdat"`
}
