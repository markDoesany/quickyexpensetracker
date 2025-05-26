package api

import (
	"quickyexpensetracker/database"
	"quickyexpensetracker/models"
	"time"
)

func SaveReminder(userID string, amount float64, accountName string, gcashNumber string, dueDate time.Time) error {
	reminder := models.RemindersLog{
		Amount:      amount,
		GcashNumber: gcashNumber,
		Rececipient: accountName,
		DueDate:     dueDate,
		UserID:      userID,
	}

	result := database.DB.Create(&reminder)

	return result.Error
}

func GetReminders(userID string, status string) ([]models.RemindersLog, error) {
	var reminders []models.RemindersLog
	now := time.Now()

	query := database.DB.Where("user_id = ?", userID)

	if status == "pending" {
		query = query.Where("due_date > ?", now)
	} else if status == "completed" {
		query = query.Where("due_date <= ?", now)
	}

	result := query.Find(&reminders)
	return reminders, result.Error
}

func DeleteRemindersByUser(userID string) error {
	result := database.DB.Where("user_id = ?", userID).Delete(&models.RemindersLog{})
	return result.Error
}
