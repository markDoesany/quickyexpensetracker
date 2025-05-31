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
	Amount        float64   `json:"amount"`
	Recipient     string    `json:"recipient"`
	GcashNumber   string    `json:"gcash_number"`
	DueDate       time.Time `json:"due_date"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	UserID        string    `json:"user_id"`
	Notified      bool      `json:"notified"` // New field
}
