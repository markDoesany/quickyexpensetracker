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

type ReportScheduleLog struct {
	gorm.Model
	UserID        string    `gorm:"not null"`
	Frequency     string    `gorm:"not null"` // e.g., "daily", "weekly", "monthly"
	DayOfWeek     int       // For weekly reports (0=Sunday, 6=Saturday) - nullable/zero if not weekly
	DayOfMonth    int       // For monthly reports (1-31) - nullable/zero if not monthly
	ScheduledTime string    `gorm:"not null"` // e.g., "09:00" (HH:MM format, 24-hour)
	LastSentAt    time.Time `gorm:"nullable"`
	Timezone      string    // e.g., "Asia/Manila", "UTC" - Important for accurate scheduling
}
