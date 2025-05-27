package api

import (
	"quickyexpensetracker/database"
	"quickyexpensetracker/models"
	"strconv"
	"time"
)

func SaveReminder(userID string, amount float64, accountName string, gcashNumber string, dueDate time.Time, paymentMethod string, status string) error {
	reminder := models.RemindersLog{
		Amount:        amount,
		GcashNumber:   gcashNumber,
		Recipient:     accountName,
		DueDate:       dueDate,
		PaymentMethod: paymentMethod,
		Status:        status,
		UserID:        userID,
	}

	result := database.DB.Create(&reminder)

	return result.Error
}

func GetReminders(userID string, status string) ([]models.RemindersLog, error) {
	var reminders []models.RemindersLog

	query := database.DB.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	result := query.Find(&reminders)
	return reminders, result.Error
}

func DeleteRemindersByUser(userID string) error {
	result := database.DB.Where("user_id = ?", userID).Delete(&models.RemindersLog{})
	return result.Error
}

// GetGcashDeepLink returns a hardcoded GCash deep link.
// reminderID is not used currently but is kept for future enhancements.
func GetGcashDeepLink(reminderID string) (string, error) {
	return "gcash://", nil
}

func UpdateReminderStatus(reminderID string, newStatus string) error {
	reminderIDUint, err := strconv.ParseUint(reminderID, 10, 64)
	if err != nil {
		return err
	}

	result := database.DB.Model(&models.RemindersLog{}).Where("id = ?", reminderIDUint).Update("status", newStatus)
	return result.Error
}

func GetDueReminders(targetDate time.Time) ([]models.RemindersLog, error) {
	var dueReminders []models.RemindersLog

	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 23, 59, 59, 999999999, targetDate.Location())

	result := database.DB.Where("status = ?", "pending").
		Where("due_date >= ? AND due_date <= ?", startOfDay, endOfDay).
		Find(&dueReminders)

	return dueReminders, result.Error
}
