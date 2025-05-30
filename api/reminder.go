package api

import (
	"fmt"
	"quickyexpensetracker/database"
	"quickyexpensetracker/models"
	"strconv"
	"time"
)

func SaveReminder(userID string, amount float64, accountName string, gcashNumber string, dueDate time.Time, paymentMethod string, status string, reminderType string, frequency string) error {
	reminder := models.RemindersLog{
		Amount:        amount,
		GcashNumber:   gcashNumber,
		Recipient:     accountName,
		DueDate:       dueDate,
		PaymentMethod: paymentMethod,
		Status:        status,
		UserID:        userID,
		Notified:      false, // Explicitly set Notified to false
		ReminderType:  reminderType,
		Frequency:     frequency,
	}

	result := database.DB.Create(&reminder)

	return result.Error
}

func GetReminderByID(reminderID string) (*models.RemindersLog, error) {
	var reminder models.RemindersLog
	reminderIDUint, err := strconv.ParseUint(reminderID, 10, 32) // Assuming ID is uint
	if err != nil {
		return nil, fmt.Errorf("error converting reminderID to uint: %w", err)
	}
	result := database.DB.First(&reminder, uint(reminderIDUint)) // Use uint for GORM
	if result.Error != nil {
		return nil, result.Error
	}
	return &reminder, nil
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

// GetPendingUnnotifiedReminders retrieves all reminders that are pending and for which notifications have not yet been sent.
func GetPendingUnnotifiedReminders() ([]models.RemindersLog, error) {
	var reminders []models.RemindersLog
	result := database.DB.Where("status = ? AND notified = ?", "pending", false).Find(&reminders)
	return reminders, result.Error
}

// MarkReminderAsNotified updates the 'notified' status of a specific reminder to true.
func MarkReminderAsNotified(reminderID string) error {
	reminderIDUint, err := strconv.ParseUint(reminderID, 10, 64) // gorm.Model ID is uint
	if err != nil {
		return fmt.Errorf("error converting reminderID to uint: %w", err)
	}

	result := database.DB.Model(&models.RemindersLog{}).Where("id = ?", reminderIDUint).Update("notified", true)
	return result.Error
}

// UpdateReminderDueDateAndNotifiedStatus updates the 'DueDate' and 'notified' status of a specific reminder.
func UpdateReminderDueDateAndNotifiedStatus(reminderID string, newDueDate time.Time, notified bool) error {
	reminderIDUint, err := strconv.ParseUint(reminderID, 10, 64) // gorm.Model ID is uint
	if err != nil {
		return fmt.Errorf("error converting reminderID to uint: %w", err)
	}

	result := database.DB.Model(&models.RemindersLog{}).Where("id = ?", reminderIDUint).Updates(models.RemindersLog{DueDate: newDueDate, Notified: notified})
	return result.Error
}
