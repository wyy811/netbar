package domain

import (
	"encoding/json"
	"time"
)

type Auditlog struct {
	ID         uint            `json:"id"`
	UserID     uint            `json:"userid"`
	Action     string          `json:"action"`
	TargetType string          `json:"targettype"`
	TargetID   uint            `json:"targetid"`
	Detail     json.RawMessage `json:"detail"`
	IpAdress   string          `json:"ipadress"`
	UserAgent  string          `json:"useragent"`
	CreatedAt  time.Time       `json:"createdat"`
}
