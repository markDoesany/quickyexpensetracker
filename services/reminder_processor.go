package services

import (
	"fmt"
	"quickyexpensetracker/api"
	"quickyexpensetracker/models"
	"quickyexpensetracker/utils"
	"time"
	"os" // Required for Getenv
)

// CheckDueReminders fetches pending reminders, checks if they are due,
// sends notifications, and marks them as notified.
func CheckDueReminders() {
	fmt.Println("Reminder Processor: Checking for due reminders...")

	// Note: The bot token retrieval might need adjustment based on how it's managed in the project.
	// Assuming it's an environment variable for now.
	token := os.Getenv("FB_PAGE_ACCESS_TOKEN") // Or your bot's access token variable
	if token == "" {
		fmt.Println("Reminder Processor: Error - FB_PAGE_ACCESS_TOKEN not set. Cannot send messages.")
		return
	}

	reminders, err := api.GetPendingUnnotifiedReminders() // This function needs to be created in api/reminder.go
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
	// Truncate current time to day for accurate date comparison (ignores time part)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())


	for _, reminder := range reminders {
		// Truncate due date to day for accurate date comparison
		dueDate := time.Date(reminder.DueDate.Year(), reminder.DueDate.Month(), reminder.DueDate.Day(), 0, 0, 0, 0, reminder.DueDate.Location())

		if dueDate.Before(today) || dueDate.Equal(today) { // Check if due date is today or in the past
			fmt.Printf("Reminder Processor: Reminder ID %d for User %s is due.\n", reminder.ID, reminder.UserID)

			message := fmt.Sprintf("Hi there! This is a friendly reminder that your payment of ₱%.2f to %s is due today (%s).",
				reminder.Amount, reminder.Recipient, reminder.DueDate.Format("Jan 2, 2006"))

			// Send the message
			// The PSID is the UserID stored in the reminder
			err := utils.SendTextMessage(message, reminder.UserID, token)
			if err != nil {
				fmt.Printf("Reminder Processor: Error sending message for reminder ID %d to User %s: %v\n", reminder.ID, reminder.UserID, err)
				// Decide if we should continue to mark as notified or retry later
				// For now, we'll continue and mark as notified to avoid spamming on next check if error is persistent
			} else {
				fmt.Printf("Reminder Processor: Notification sent for reminder ID %d to User %s.\n", reminder.ID, reminder.UserID)
			}

			// Mark as notified
			err = api.MarkReminderAsNotified(fmt.Sprint(reminder.ID)) // This function needs to be created in api/reminder.go
			if err != nil {
				fmt.Printf("Reminder Processor: Error marking reminder ID %d as notified: %v\n", reminder.ID, err)
			} else {
				fmt.Printf("Reminder Processor: Reminder ID %d marked as notified.\n", reminder.ID)
			}
		}
	}
	fmt.Println("Reminder Processor: Finished checking due reminders.")
}
