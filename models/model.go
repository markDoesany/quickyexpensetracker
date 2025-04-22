package models

import (
	"time"

	"gorm.io/gorm"
)

type ExpensesLog struct {
	gorm.Model
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	UserID   string  `json:"user_id"`
}

type RemindersLog struct {
	gorm.Model
	Amount      float64   `json:"amount"`
	Rececipient string    `json:"recipient"`
	GcashNumber string    `json:"gcash_number"`
	DueDate     time.Time `json:"due_date"`
	UserID      string    `json:"user_id"`
}
