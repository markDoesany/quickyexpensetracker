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

func DeleteExpensesByUser(userID string) error {
	result := database.DB.Where("user_id = ?", userID).Delete(&models.ExpensesLog{})
	return result.Error
}

// GetExpensesForPeriod retrieves expenses for a user within a specific date range and calculates the total amount.
// Note: GORM typically uses `created_at` for timestamping. Adjust field name if different in your model.
func GetExpensesForPeriod(userID string, periodStartDate time.Time, periodEndDate time.Time) ([]models.ExpensesLog, float64, error) {
	var expenses []models.ExpensesLog
	var totalAmount float64

	// Ensure periodEndDate is at the end of its day to be inclusive for the whole day.
	// For example, if periodEndDate is 2023-10-26, we want to include all expenses on that day.
	// So, the query should be < 2023-10-27 00:00:00
	// Or, if using direct date comparison, ensure the time part of periodEndDate is 23:59:59.999
	// For simplicity, we'll use exclusive end date (CreatedAt < periodEndDate.AddDate(0,0,1)) if periodEndDate is just a date.
	// However, the provided periodEndDate is already a time.Time, so we can use it directly if it's set to EOD, or use a strict < for the next day.
	// Let's assume periodEndDate is exclusive for simplicity here. If it needs to be inclusive, the query would adjust.
	// A common pattern is StartDate <= CreatedAt < EndDate (exclusive of EndDate)
	// Or StartDate <= CreatedAt <= EndDate (inclusive of EndDate, ensure EndDate has time set to 23:59:59)

	result := database.DB.
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, periodStartDate, periodEndDate).
		Order("created_at asc"). // Order by date for easier reading of summaries if needed
		Find(&expenses)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	for _, expense := range expenses {
		totalAmount += expense.Amount
	}

	return expenses, totalAmount, nil
}
