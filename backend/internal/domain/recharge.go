package domain

import (
	"time"
)

type Recharge struct {
	ID              uint      `json:"id"`
	MemberID        int       `json:"memberid"`
	UserID          int       `json:"userid"`
	DurationMinutes int       `json:"durationminutes"`
	Amount          float64   `json:"amount"`
	BonusAmount     float64   `json:"bonusamount"`
	PaymentMethod   int       `json:"paymentmethod"`
	BeforeBalance   float64   `json:"beforebalance"`
	AfterBalance    float64   `json:"afterbalance"`
	Remark          string    `json:"remark"`
	CreatedAt       time.Time `json:"createdat"`
}
