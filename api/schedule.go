package api

import (
	"errors"
	"fmt"
	"quickyexpensetracker/database"
	"quickyexpensetracker/models"
	"regexp"
	"time"

	"gorm.io/gorm"
)

// SetSchedule creates or updates a report schedule for a user.
func SetSchedule(userID, frequency string, dayOfWeek, dayOfMonth int, scheduledTime, timezone string) (*models.ReportScheduleLog, error) {
	// Validate frequency
	if frequency != "daily" && frequency != "weekly" && frequency != "monthly" {
		return nil, errors.New("invalid frequency: must be 'daily', 'weekly', or 'monthly'")
	}

	// Validate scheduledTime format (HH:MM)
	timeRegex := regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`)
	if !timeRegex.MatchString(scheduledTime) {
		return nil, errors.New("invalid scheduledTime format: must be HH:MM")
	}

	// Validate DayOfWeek for weekly frequency
	if frequency == "weekly" && (dayOfWeek < 0 || dayOfWeek > 6) {
		return nil, errors.New("invalid dayOfWeek: must be between 0 (Sunday) and 6 (Saturday) for weekly frequency")
	}

	// Validate DayOfMonth for monthly frequency
	if frequency == "monthly" && (dayOfMonth < 1 || dayOfMonth > 31) {
		return nil, errors.New("invalid dayOfMonth: must be between 1 and 31 for monthly frequency")
	}
	
	// Validate Timezone (basic check, can be expanded)
	if timezone == "" {
		// Default to UTC if not provided, or handle as an error
		// For now, let's require it.
		return nil, errors.New("timezone is required")
	}
	// A more robust validation would involve checking against a list of valid timezones
    // e.g., using time.LoadLocation(timezone) and checking for errors.
    // However, for this basic implementation, we'll assume it's valid if provided.


	schedule := models.ReportScheduleLog{
		UserID:        userID,
		Frequency:     frequency,
		ScheduledTime: scheduledTime,
		Timezone:      timezone,
	}

	if frequency == "weekly" {
		schedule.DayOfWeek = dayOfWeek
		schedule.DayOfMonth = 0 // Ensure DayOfMonth is not set for weekly
	} else if frequency == "monthly" {
		schedule.DayOfMonth = dayOfMonth
		schedule.DayOfWeek = 0 // Ensure DayOfWeek is not set for monthly
	} else { // daily
		schedule.DayOfWeek = 0
		schedule.DayOfMonth = 0
	}

	// Check if a schedule already exists for this user
	var existingSchedule models.ReportScheduleLog
	err := database.DB.Where("user_id = ?", userID).First(&existingSchedule).Error

	if err == nil {
		// Update existing schedule
		existingSchedule.Frequency = schedule.Frequency
		existingSchedule.DayOfWeek = schedule.DayOfWeek
		existingSchedule.DayOfMonth = schedule.DayOfMonth
		existingSchedule.ScheduledTime = schedule.ScheduledTime
		existingSchedule.Timezone = schedule.Timezone
		// existingSchedule.LastSentAt can be reset or handled as needed upon update
		result := database.DB.Save(&existingSchedule)
		if result.Error != nil {
			return nil, result.Error
		}
		return &existingSchedule, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new schedule
		result := database.DB.Create(&schedule)
		if result.Error != nil {
			return nil, result.Error
		}
		return &schedule, nil
	} else {
		// Other database error
		return nil, err
	}
}

// GetSchedule retrieves the active schedule for a user.
func GetSchedule(userID string) (*models.ReportScheduleLog, error) {
	var schedule models.ReportScheduleLog
	result := database.DB.Where("user_id = ?", userID).First(&schedule)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound // Propagate gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &schedule, nil
}

// DeleteSchedule deletes a user's schedule.
func DeleteSchedule(userID string) error {
	result := database.DB.Where("user_id = ?", userID).Delete(&models.ReportScheduleLog{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Or a custom error indicating no schedule to delete
	}
	return nil
}

// GetDueSchedules finds all schedules that are due to be sent.
// This is a placeholder for more complex logic.
func GetDueSchedules(currentTime time.Time) ([]models.ReportScheduleLog, error) {
	var dueSchedules []models.ReportScheduleLog
	// Basic idea:
	// 1. Iterate through all schedules.
	// 2. For each schedule, calculate its next expected send time in UTC based on its
	//    Frequency, DayOfWeek/DayOfMonth, ScheduledTime, and Timezone.
	// 3. If currentTime (in UTC) is after this next expected send time AND
	//    (LastSentAt is null OR LastSentAt is before this next expected send time),
	//    then the schedule is due.
	//
	// This is a simplified version that doesn't implement the full logic yet.
	// A full implementation would require careful handling of timezones and schedule frequencies.
	// For now, it returns an empty list, indicating the core logic needs to be built.
	// Example: Find schedules that should have run today but haven't.
	// This will require converting currentTime to each schedule's timezone.
	
	// Placeholder:
	// result := database.DB.Find(&dueSchedules) 
	// return dueSchedules, result.Error

	fmt.Println("GetDueSchedules: Full logic to be implemented. currentTime:", currentTime.Format(time.RFC3339))
	// This function will be significantly more complex and will be refined.
	// For now, we can fetch all schedules and the filtering logic will be outside, or
	// we can try a very naive query.
	// Let's return all schedules for now and assume the worker will filter.
	// This is NOT a production-ready implementation for GetDueSchedules.
	result := database.DB.Find(&dueSchedules)
	return dueSchedules, result.Error
}
