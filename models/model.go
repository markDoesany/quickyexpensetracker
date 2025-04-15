package models

import (
	"gorm.io/gorm"
)

type ExpensesLog struct {
	gorm.Model
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	UserID   int     `json:"user_id"`
}

type RemindersLog struct {
	gorm.Model
	Amount      float64 `json:"amount"`
	Rececipient string  `json:"recipient"`
	GcashNumber string  `json:"gcash_number"`
	DueDate     string  `json:"due_date"`
	UserID      int     `json:"user_id"`
}
