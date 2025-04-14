package models

import (
	"gorm.io/gorm"
)

type ExpensesLog struct {
	gorm.Model
	Amount   float64
	Category string
	UserID   int
}

type RemindersLog struct {
	gorm.Model
	Amount      float64
	Rececipient string
	GcashNumber string
	DueDate     string
	UserID      int
}
