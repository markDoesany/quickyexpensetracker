package services

import (
	"fmt"
	"os" // Required for Getenv
	"quickyexpensetracker/api"
	"quickyexpensetracker/models"
	"quickyexpensetracker/utils" // Assuming utils.CalculateNextDueDate will be here
	"time"
)

// CheckDueReminders fetches pending reminders, checks if they are due,
// sends notifications, and updates them according to their type and frequency.
func CheckDueReminders() {
	fmt.Println("Reminder Processor: Checking for due reminders...")

	token := os.Getenv("FB_PAGE_ACCESS_TOKEN")
	if token == "" {
		fmt.Println("Reminder Processor: Error - FB_PAGE_ACCESS_TOKEN not set. Cannot send messages.")
		return
	}

	reminders, err := api.GetPendingUnnotifiedReminders()
	if err != nil {
		fmt.Printf("Reminder Processor: Error fetching reminders: %v\n", err)
		return
	}

	if len(reminders) == 0 {
		fmt.Println("Reminder Processor: No pending unnotified reminders found.")
		return
	}

	fmt.Printf("Reminder Processor: Found %d pending unnotified reminders.\n", len(reminders))

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for _, reminder := range reminders {
		dueDate := time.Date(reminder.DueDate.Year(), reminder.DueDate.Month(), reminder.DueDate.Day(), 0, 0, 0, 0, reminder.DueDate.Location())

		if dueDate.Before(today) || dueDate.Equal(today) {
			fmt.Printf("Reminder Processor: Processing Reminder ID %d (Type: %s, Frequency: %s) for User %s.\n", reminder.ID, reminder.ReminderType, reminder.Frequency, reminder.UserID)

			var message string
			var processingError error
			var notificationSent bool

			switch reminder.ReminderType {
			case "payment":
				message = fmt.Sprintf("Hi there! This is a friendly reminder that your payment of ₱%.2f to %s is due today (%s).",
					reminder.Amount, reminder.Recipient, reminder.DueDate.Format("Jan 2, 2006"))

				err = utils.SendTextMessage(message, reminder.UserID, token)
				if err != nil {
					processingError = fmt.Errorf("error sending payment notification: %w", err)
				} else {
					notificationSent = true
					fmt.Printf("Reminder Processor: Payment notification sent for reminder ID %d.\n", reminder.ID)
				}

			case "expense_summary":
				// Determine period for expense summary
				var periodStartDate, periodEndDate time.Time
				periodEndDate = reminder.DueDate // Summary up to the due date
				// This is a simplified assumption for period start.
				// A more robust solution might be needed in utils.CalculatePeriod.
				switch reminder.Frequency {
				case "daily":
					periodStartDate = periodEndDate.AddDate(0, 0, -1)
				case "weekly":
					periodStartDate = periodEndDate.AddDate(0, 0, -7)
				case "monthly":
					periodStartDate = periodEndDate.AddDate(0, -1, 0)
				default: // "once" or unknown, summarize for the day of the reminder
					periodStartDate = time.Date(reminder.DueDate.Year(), reminder.DueDate.Month(), reminder.DueDate.Day(), 0, 0, 0, 0, reminder.DueDate.Location())
					periodEndDate = periodStartDate.AddDate(0,0,1).Add(-time.Nanosecond) // Full day of reminder.DueDate
				}

				expenses, totalAmount, err := api.GetExpensesForPeriod(reminder.UserID, periodStartDate, periodEndDate)
				if err != nil {
					processingError = fmt.Errorf("error fetching expenses for summary: %w", err)
					break // Break from switch, will go to error logging
				}

				summaryMessage := utils.GenerateExpenseSummaryMessage(reminder.Frequency, expenses, totalAmount)

				err = utils.SendTextMessage(summaryMessage, reminder.UserID, token)
				if err != nil {
					processingError = fmt.Errorf("error sending expense summary: %w", err)
				} else {
					notificationSent = true
					fmt.Printf("Reminder Processor: Expense summary sent for reminder ID %d.\n", reminder.ID)
				}

			default:
				fmt.Printf("Reminder Processor: Unknown reminder type '%s' for reminder ID %d. Skipping.\n", reminder.ReminderType, reminder.ID)
				// Optionally mark as notified to prevent reprocessing if it's an unrecoverable situation.
				// For now, we'll let it be, it might be a temporary issue or a new type not yet supported.
				// Consider adding `api.MarkReminderAsNotified(fmt.Sprint(reminder.ID))` if that's desired.
				continue // Skip to next reminder
			}

			if processingError != nil {
				fmt.Printf("Reminder Processor: Error processing reminder ID %d: %v\n", reminder.ID, processingError)
				// Decide on retry logic or if to mark as notified despite error.
				// For now, if notification failed, we don't update status to allow retry on next cycle.
				// If other processing failed (e.g. fetching expenses), notification wouldn't have been sent.
				continue // Skip to next reminder
			}

			// If notification was attempted (either successfully or failed but we want to proceed)
			if notificationSent || reminder.ReminderType == "payment" { // Ensure payment reminders are always updated or marked
				if reminder.Frequency == "once" || reminder.Frequency == "" { // Treat empty frequency as "once"
					err = api.MarkReminderAsNotified(fmt.Sprint(reminder.ID))
					if err != nil {
						fmt.Printf("Reminder Processor: Error marking reminder ID %d as notified: %v\n", reminder.ID, err)
					} else {
						fmt.Printf("Reminder Processor: Reminder ID %d (once) marked as notified.\n", reminder.ID)
					}
				} else if reminder.Frequency == "daily" || reminder.Frequency == "weekly" || reminder.Frequency == "monthly" {
					nextDueDate, err := utils.CalculateNextDueDate(reminder.DueDate, reminder.Frequency)
					if err != nil {
						fmt.Printf("Reminder Processor: Error calculating next due date for reminder ID %d: %v. Skipping reschedule.\n", reminder.ID, err)
						// Mark as notified to prevent loop if calculation consistently fails for this reminder
						if markErr := api.MarkReminderAsNotified(fmt.Sprint(reminder.ID)); markErr != nil {
							fmt.Printf("Reminder Processor: Error marking reminder ID %d as notified after calc error: %v\n", reminder.ID, markErr)
						}
						continue
					}
					err = api.UpdateReminderDueDateAndNotifiedStatus(fmt.Sprint(reminder.ID), nextDueDate, false)
					if err != nil {
						fmt.Printf("Reminder Processor: Error updating reminder ID %d for next occurrence: %v\n", reminder.ID, err)
					} else {
						fmt.Printf("Reminder Processor: Reminder ID %d rescheduled to %s. Notified status reset.\n", reminder.ID, nextDueDate.Format("Jan 2, 2006"))
					}
				} else {
					fmt.Printf("Reminder Processor: Unknown frequency '%s' for reminder ID %d. Marking as notified to prevent loop.\n", reminder.Frequency, reminder.ID)
					err = api.MarkReminderAsNotified(fmt.Sprint(reminder.ID))
					if err != nil {
						fmt.Printf("Reminder Processor: Error marking reminder ID %d with unknown frequency as notified: %v\n", reminder.ID, err)
					}
				}
			}
		}
	}
	fmt.Println("Reminder Processor: Finished checking due reminders.")
}
