package utils

import (
	"fmt"
	"time"
)

// CalculateNextDueDate calculates the next due date based on the current due date and frequency.
// It resets the time to the beginning of the day (00:00:00).
func CalculateNextDueDate(currentDueDate time.Time, frequency string) (time.Time, error) {
	// Reset currentDueDate to the beginning of its day to ensure consistent calculations
	currentDate := time.Date(currentDueDate.Year(), currentDueDate.Month(), currentDueDate.Day(), 0, 0, 0, 0, currentDueDate.Location())

	switch frequency {
	case "daily":
		return currentDate.AddDate(0, 0, 1), nil
	case "weekly":
		return currentDate.AddDate(0, 0, 7), nil
	case "monthly":
		return currentDate.AddDate(0, 1, 0), nil
	case "once":
		// For "once" frequency, it implies no next due date from recurrence.
		// However, this function is typically called for rescheduling.
		// Returning an error or a zero time might be appropriate.
		// For now, let's return an error as it's unexpected in a rescheduling context.
		return time.Time{}, fmt.Errorf("cannot calculate next due date for 'once' frequency")
	default:
		return time.Time{}, fmt.Errorf("unknown frequency: %s", frequency)
	}
}
