package domain

import (
	"time"
)

type User struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	RealName    string    `json:"realName"`
	Role        int       `json:"role"`
	Status      int       `json:"status"`
	LastLoginAt time.Time `json:"lastLoginat"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
