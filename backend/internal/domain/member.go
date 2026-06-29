package domain

import (
	"time"
)

type Member struct {
	ID            uint      `json:"id"`
	CardNumber    string    `json:"carnumber"`
	Name          string    `json:"Name"`
	Phone         string    `json:"phone"`
	IdCard        string    `json:"idcard"`
	Balance       float64   `json:"balance"`
	TotalSpend    float64   `json:"totalspend"`
	DiscountLevel int       `json:"discountlevel"`
	Status        int       `json:"status"`
	RegisteredAt  time.Time `json:"registerat"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
