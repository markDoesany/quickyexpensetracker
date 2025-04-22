package api

import (
	"errors"
	"quickyexpensetracker/database"
	"quickyexpensetracker/models"
	"time"
)

func SaveExpense(amount float64, category string, psid string) error {
	expense := models.ExpensesLog{
		Amount:   amount,
		Category: category,
		UserID:   psid,
	}

	result := database.DB.Create(&expense)

	return result.Error
}

func GetExpensesByUserAndRange(userID string, rangeType string) ([]models.ExpensesLog, error) {
	var expenses []models.ExpensesLog
	var startTime time.Time

	now := time.Now()

	switch rangeType {
	case "day":
		startTime = now.AddDate(0, 0, -1) // last 24 hours
	case "week":
		startTime = now.AddDate(0, 0, -7) // last 7 days
	case "month":
		startTime = now.AddDate(0, -1, 0) // last 1 month
	default:
		return nil, errors.New("invalid range type: choose 'day', 'week', or 'month'")
	}

	result := database.DB.
		Where("user_id = ? AND created_at >= ?", userID, startTime).
		Order("created_at desc").
		Find(&expenses)

	return expenses, result.Error
}
