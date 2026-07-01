package domain

import (
	"time"
)

type Session struct {
	ID              uint      `json:"id"`
	ComputerID      int       `json:"computerid"`
	MemberID        int       `json:"memberid"`
	UserID          int       `json:"userid"`
	StartTime       time.Time `json:"starttime"`
	EndTime         time.Time `json:"endtime"`
	DurationMinutes int       `json:"durationminutes"`
	TotalAmount     float64   `json:"totalamount"`
	DiscountAmount  float64   `json:"discountamount"`
	PaidAmount      float64   `json:"paidamount"`
	PaymentMethod   int       `json:"paymentmethod"`
	Status          int       `json:"status"`
	Note            string    `json:"note"`
	CreatedAt       time.Time `json:"createdat"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
