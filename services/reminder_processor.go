package services

import (
	"fmt"
	"os"
	"quickyexpensetracker/api"
	"quickyexpensetracker/templates"
	"quickyexpensetracker/utils"
	"time"
)

// CheckDueReminders fetches pending reminders, checks if they are due,
// sends notifications, and updates them according to their type and frequency.
func CheckDueReminders() {
	fmt.Println("Reminder Processor: Checking for due reminders...")

	token := os.Getenv("PAGE_TOKEN")
	if token == "" {
		fmt.Println("Reminder Processor: Error - PAGE_TOKEN not set. Cannot send messages.")
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
			fmt.Printf("Reminder Processor: Processing Reminder ID %d (Type: %s, Frequency: %s) for User %s.\n",
				reminder.ID, reminder.ReminderType, reminder.Frequency, reminder.UserID)

			var processingError error
			var notificationSent bool

			switch reminder.ReminderType {
			case "payment":
				// Send plain text message first
				message := fmt.Sprintf("Hi there! This is a friendly reminder that your payment of ₱%.2f to %s is due today (%s).",
					reminder.Amount, reminder.Recipient, reminder.DueDate.Format("Jan 2, 2006"))

				err := utils.SendTextMessage(message, reminder.UserID, token)
				if err != nil {
					processingError = fmt.Errorf("error sending payment notification: %w", err)
					break
				}

				// Then send the payment template with buttons
				title := fmt.Sprintf("Payment to %s", reminder.Recipient)
				subtitle := fmt.Sprintf("Amount: ₱%.2f\nGCash: %s\nDue: %s",
					reminder.Amount,
					reminder.GcashNumber,
					dueDate.Format("2006-01-02"))

				var buttons []templates.Button
				if reminder.PaymentMethod == "Gcash" && reminder.Status == "pending" {
					buttons = append(buttons, templates.Button{
						Type:    "postback",
						Title:   "Pay with GCash",
						Payload: "PAY_GCASH_" + fmt.Sprint(reminder.ID),
					})
					buttons = append(buttons, templates.Button{
						Type:    "postback",
						Title:   "Mark as Paid",
						Payload: "MARK_AS_PAID_" + fmt.Sprint(reminder.ID),
					})
				}

				element := templates.Template{
					Title:    title,
					Subtitle: subtitle,
					Buttons:  buttons,
				}

				// Send the template
				err = utils.SendTemplateMessage([]templates.Template{element}, reminder.UserID, token)
				if err != nil {
					processingError = fmt.Errorf("error sending payment template: %w", err)
				} else {
					notificationSent = true
					fmt.Printf("Reminder Processor: Payment template sent for reminder ID %d.\n", reminder.ID)
				}

			case "expense_summary":
				// Determine period for expense summary
				var periodStartDate, periodEndDate time.Time
				periodEndDate = reminder.DueDate
				switch reminder.Frequency {
				case "daily":
					periodStartDate = periodEndDate.AddDate(0, 0, -1)
				case "weekly":
					periodStartDate = periodEndDate.AddDate(0, 0, -7)
				case "monthly":
					periodStartDate = periodEndDate.AddDate(0, -1, 0)
				default:
					periodStartDate = time.Date(reminder.DueDate.Year(), reminder.DueDate.Month(),
						reminder.DueDate.Day(), 0, 0, 0, 0, reminder.DueDate.Location())
					periodEndDate = periodStartDate.AddDate(0, 0, 1).Add(-time.Nanosecond)
				}

				expenses, totalAmount, err := api.GetExpensesForPeriod(reminder.UserID, periodStartDate, periodEndDate)
				if err != nil {
					processingError = fmt.Errorf("error fetching expenses for summary: %w", err)
					break
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
				fmt.Printf("Reminder Processor: Unknown reminder type '%s' for reminder ID %d. Skipping.\n",
					reminder.ReminderType, reminder.ID)
				continue
			}

			if processingError != nil {
				fmt.Printf("Reminder Processor: Error processing reminder ID %d: %v\n",
					reminder.ID, processingError)
				continue
			}

			// Handle reminder status update based on frequency
			if notificationSent || reminder.ReminderType == "payment" {
				if reminder.Frequency == "once" || reminder.Frequency == "" {
					err = api.MarkReminderAsNotified(fmt.Sprint(reminder.ID))
					if err != nil {
						fmt.Printf("Reminder Processor: Error marking reminder ID %d as notified: %v\n",
							reminder.ID, err)
					} else {
						fmt.Printf("Reminder Processor: Reminder ID %d (once) marked as notified.\n", reminder.ID)
					}
				} else if reminder.Frequency == "daily" || reminder.Frequency == "weekly" || reminder.Frequency == "monthly" {
					nextDueDate, err := utils.CalculateNextDueDate(reminder.DueDate, reminder.Frequency)
					if err != nil {
						fmt.Printf("Reminder Processor: Error calculating next due date for reminder ID %d: %v. Marking as notified.\n",
							reminder.ID, err)
						if markErr := api.MarkReminderAsNotified(fmt.Sprint(reminder.ID)); markErr != nil {
							fmt.Printf("Reminder Processor: Error marking reminder ID %d as notified after calc error: %v\n",
								reminder.ID, markErr)
						}
						continue
					}

					err = api.UpdateReminderDueDateAndNotifiedStatus(fmt.Sprint(reminder.ID), nextDueDate, false)
					if err != nil {
						fmt.Printf("Reminder Processor: Error updating reminder ID %d for next occurrence: %v\n",
							reminder.ID, err)
					} else {
						fmt.Printf("Reminder Processor: Reminder ID %d rescheduled to %s. Notified status reset.\n",
							reminder.ID, nextDueDate.Format("Jan 2, 2006"))
					}
				} else {
					fmt.Printf("Reminder Processor: Unknown frequency '%s' for reminder ID %d. Marking as notified.\n",
						reminder.Frequency, reminder.ID)
					if err := api.MarkReminderAsNotified(fmt.Sprint(reminder.ID)); err != nil {
						fmt.Printf("Reminder Processor: Error marking reminder ID %d with unknown frequency as notified: %v\n",
							reminder.ID, err)
					}
				}
			}
		}
	}
	fmt.Println("Reminder Processor: Finished checking due reminders.")
}
