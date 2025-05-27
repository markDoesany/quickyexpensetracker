package services

import (
	"fmt"
	"log"
	"os"
	"quickyexpensetracker/api"
	"quickyexpensetracker/models"
	"quickyexpensetracker/utils"
	"time"
)

// SendDueReminderNotifications fetches reminders due today and sends notifications.
func SendDueReminderNotifications() {
	log.Println("Scheduler: Running SendDueReminderNotifications job.")
	today := time.Now()

	reminders, err := api.GetDueReminders(today)
	if err != nil {
		log.Printf("Scheduler: Error fetching due reminders: %v\n", err)
		return
	}

	if len(reminders) == 0 {
		log.Println("Scheduler: No reminders due today.")
		return
	}

	log.Printf("Scheduler: Found %d reminder(s) due today.\n", len(reminders))

	pageAccessToken := os.Getenv("PAGE_TOKEN")
	if pageAccessToken == "" {
		log.Println("Scheduler: Error: PAGE_TOKEN environment variable not set.")
		return
	}

	for _, reminder := range reminders {
		psid := reminder.UserID
		
		// Use utils.GetRemindersReport to generate the template for the single reminder.
		// GetRemindersReport expects a slice and returns a slice of templates.
		reminderSlice := []models.RemindersLog{reminder}
		notificationTemplates := utils.GetRemindersReport(reminderSlice)

		if len(notificationTemplates) > 0 {
			// Assuming GetRemindersReport for a single reminder returns one template.
			notificationTemplate := notificationTemplates[0] 

			log.Printf("Scheduler: Sending notification for reminder ID %d to user %s\n", reminder.ID, psid)
			err := utils.SendGenerateRequest(notificationTemplate, psid, pageAccessToken)
			if err != nil {
				log.Printf("Scheduler: Error sending notification to user %s for reminder ID %d: %v\n", psid, reminder.ID, err)
			} else {
				log.Printf("Scheduler: Successfully sent notification for reminder ID %d to user %s\n", reminder.ID, psid)
			}
		} else {
			log.Printf("Scheduler: Could not generate notification template for reminder ID %d for user %s\n", reminder.ID, psid)
		}
	}
	log.Println("Scheduler: Finished SendDueReminderNotifications job.")
}
