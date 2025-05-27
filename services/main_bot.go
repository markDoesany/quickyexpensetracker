package services

import (
	"fmt"
	"quickyexpensetracker/api"
	"quickyexpensetracker/templates"
	"quickyexpensetracker/utils"
	"strings"
	"time"
)

var userState = make(map[string]string)
var currentTime = time.Now()

func ProcessMainCommand(command, psid, mid, token string) {
	fmt.Printf("Processing Command: %s, PSID: %s, MID: %s\n", command, psid, mid)

	switch command {
	case "GET_STARTED":
		utils.SendGenerateRequest(templates.MenuTemplate[1], psid, token)
	case "LOG_EXPENSES_MENU":
		utils.SendGenerateRequest(templates.MenuTemplate[2], psid, token)
	case "LOG_EXPENSES":
		ProcessTextMessageSent("LOG_EXPENSE_MESSAGE", psid, mid, token)
	case "GENERATE_REPORT_SUBMENU":
		utils.SendGenerateRequest(templates.SubMenuTemplate[1], psid, token)
	case "GENERATE_REPORT_DAY":
		ProcessTextMessageSent("REPORT_LOG_DAY", psid, mid, token)
	case "GENERATE_REPORT_WEEK":
		ProcessTextMessageSent("REPORT_LOG_WEEK", psid, mid, token)
	case "GENERATE_REPORT_MONTH":
		ProcessTextMessageSent("REPORT_LOG_MONTH", psid, mid, token)
	case "REMIND_PAYMENTS_MENU":
		utils.SendGenerateRequest(templates.MenuTemplate[3], psid, token)
	case "VIEW_PENDING_PAYMENTS":
		ProcessTextMessageSent("VIEW_PENDING_PAYMENTS_MESSAGE", psid, mid, token)
	case "VIEW_ACCOMPLISHED_PAYMENTS":
		ProcessTextMessageSent("VIEW_ACCOMPLISHED_PAYMENTS_MESSAGE", psid, mid, token)
	case "SUBSCRIPTION_STATUS":
		ProcessTextMessageSent("SUBSCRIPTION_STATUS_MESSAGE", psid, mid, token)
	case "EXPAND_MENU":
		utils.SendGenerateRequest(templates.MenuTemplate[4], psid, token)
	case "SET_REPORT_SCHED_SUBMENU":
		utils.SendGenerateRequest(templates.SubMenuTemplate[2], psid, token)
	case "RESET_LOGS":
		ProcessTextMessageSent("RESET_LOGS_MESSAGE", psid, mid, token)
	case "SET_REMINDER":
		ProcessTextMessageSent("SET_REMINDER_MESSAGE", psid, mid, token)
	default:
		if strings.HasPrefix(command, "PAY_GCASH_") {
			reminderID := strings.TrimPrefix(command, "PAY_GCASH_")
			deepLink, err := api.GetGcashDeepLink(reminderID)
			if err != nil {
				fmt.Printf("Error getting Gcash deep link for reminder %s, user %s: %v\n", reminderID, psid, err)
				utils.SendTextMessage("Sorry, could not retrieve Gcash link.", psid, token)
			} else {
				message := fmt.Sprintf("Click here to open Gcash: %s", deepLink)
				utils.SendTextMessage(message, psid, token)
			}
		} else if strings.HasPrefix(command, "MARK_AS_PAID_") {
			reminderID := strings.TrimPrefix(command, "MARK_AS_PAID_")
			err := api.UpdateReminderStatus(reminderID, "paid")
			if err != nil {
				fmt.Printf("Error updating reminder status for reminder %s, user %s: %v\n", reminderID, psid, err)
				utils.SendTextMessage("Sorry, could not update payment status.", psid, token)
			} else {
				utils.SendTextMessage("Payment has been marked as paid.", psid, token)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", command)
			utils.SendGenerateRequest(templates.MenuTemplate[1], psid, token)
		}
	}
}

func ProcessTextMessageSent(command, psid, mid, token string) {
	switch command {
	case "LOG_EXPENSE_MESSAGE":
		message := "Please log in this format: \n[amount] for [item/service]\n(e.g. 200.00 for softdrinks)"
		utils.SendTextMessage(message, psid, token)
		userState[psid] = "RECORDING_EXPENSE_LOG"
	case "REPORT_LOG_DAY":
		expenses, err := api.GetExpensesByUserAndRange(psid, "day")
		if err != nil {
			fmt.Printf("Error fetching daily expenses for user %s: %v\n", psid, err)
			utils.SendTextMessage("Sorry, I couldn't fetch your expense report at the moment. Please try again later.", psid, token)
			return
		}

		report := utils.GetExpenseReport(expenses, "Daily")
		utils.SendTextMessage(report, psid, token)
	case "REPORT_LOG_WEEK":
		expenses, err := api.GetExpensesByUserAndRange(psid, "week")
		if err != nil {
			fmt.Printf("Error fetching weekly expenses for user %s: %v\n", psid, err)
			utils.SendTextMessage("Sorry, I couldn't fetch your expense report at the moment. Please try again later.", psid, token)
			return
		}

		report := utils.GetExpenseReport(expenses, "Weekly")
		utils.SendTextMessage(report, psid, token)
	case "REPORT_LOG_MONTH":
		expenses, err := api.GetExpensesByUserAndRange(psid, "month")
		if err != nil {
			fmt.Printf("Error fetching monthly expenses for user %s: %v\n", psid, err)
			utils.SendTextMessage("Sorry, I couldn't fetch your expense report at the moment. Please try again later.", psid, token)
			return
		}

		report := utils.GetExpenseReport(expenses, "Monthly")
		utils.SendTextMessage(report, psid, token)
	case "VIEW_PENDING_PAYMENTS_MESSAGE":
		reminders, err := api.GetReminders(psid, "pending")
		if err != nil {
			fmt.Printf("Error fetching pending reminders for user %s: %v\n", psid, err)
			utils.SendTextMessage("Sorry, I couldn't fetch your payment reminders at the moment. Please try again later.", psid, token)
			return
		}
		reminderTemplates := utils.GetRemindersReport(reminders)
		for _, tmpl := range reminderTemplates {
			errLoop := utils.SendGenerateRequest(tmpl, psid, token) // Use errLoop to avoid conflict
			if errLoop != nil {
				fmt.Printf("Error sending reminder template for user %s: %v\n", psid, errLoop)
				// Optionally, send a text message to the user about the specific failure
			}
		}
	case "VIEW_ACCOMPLISHED_PAYMENTS_MESSAGE":
		reminders, err := api.GetReminders(psid, "completed")
		if err != nil {
			fmt.Printf("Error fetching accomplished reminders for user %s: %v\n", psid, err)
			utils.SendTextMessage("Sorry, I couldn't fetch your payment reminders at the moment. Please try again later.", psid, token)
			return
		}
		report := utils.GetRemindersReport(reminders)
		utils.SendTextMessage(report, psid, token)
	case "SET_REPORT_SCHED_MESSAGE":
		message := "The report scheduling feature is not yet implemented. Please check back later!"
		utils.SendTextMessage(message, psid, token)
	case "RESET_LOGS_MESSAGE":
		errExpenses := api.DeleteExpensesByUser(psid)
		if errExpenses != nil {
			fmt.Printf("Error deleting expenses for user %s: %v\n", psid, errExpenses)
			// Optionally, notify the user about the error, or log it for monitoring
		}

		errReminders := api.DeleteRemindersByUser(psid)
		if errReminders != nil {
			fmt.Printf("Error deleting reminders for user %s: %v\n", psid, errReminders)
			// Optionally, notify the user about the error, or log it for monitoring
		}

		message := "All your expense and reminder logs have been reset."
		utils.SendTextMessage(message, psid, token)
	case "SET_REMINDER_MESSAGE":
		message := "Please set the reminder in this format: \n[amount] to [name]:[gcash number] on [month/day/year]\n(e.g. 200.00 to mark:09565546*** on 04/25/2025)"
		utils.SendTextMessage(message, psid, token)
		userState[psid] = "RECORDING_REMINDER"
	case "SUBSCRIPTION_STATUS_MESSAGE":
		message := "The subscription status feature is not yet implemented. Please check back later!"
		utils.SendTextMessage(message, psid, token)
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func ProcessTextMessageReceived(message, psid, mid, token string) {
	state, exists := userState[psid]
	if exists {
		switch state {
		case "RECORDING_EXPENSE_LOG":
			if utils.IsExpenseLogFormatCorrect(message) {
				amount, category, err := utils.GetExpenseDataFromMessage(message)
				if err != nil {
					fmt.Printf("Error parsing expense data for user %s: %v\n", psid, err)
					utils.SendTextMessage("There was an issue parsing your expense. Please ensure you're using the format: [amount] for [item/service]", psid, token)
					userState[psid] = "WAITING..." // Reset state as format was correct
					return
				}

				err = api.SaveExpense(amount, category, psid)
				if err != nil {
					fmt.Printf("Error saving expense for user %s: %v\n", psid, err)
					utils.SendTextMessage("Sorry, I couldn't save your expense. Please try again later.", psid, token)
					userState[psid] = "WAITING..." // Reset state as format was correct
					return
				}
				currentTime = time.Now()
				message_ := fmt.Sprintf("Got it! You spent ₱%.2f on %s on %s", amount, category, currentTime.Format("Jan 2, 2006 at 3:04 PM"))
				utils.SendTextMessage(message_, psid, token)
				fmt.Printf("Expense saved for user %s: ₱%.2f on %s\n", psid, amount, category)
				userState[psid] = "WAITING..."
			} else {
				message = "Invalid Format. Please try again."
				utils.SendTextMessage(message, psid, token)
				ProcessMainCommand("GET_STARTED", psid, mid, token)
				userState[psid] = "WAITING..."
			}
		case "RECORDING_REMINDER":
			if utils.IsReminderLogFormatCorrect(message) {
				amount, accountName, gcashNumber, dueDate, err := utils.GetReminderDataFromMessage(message)
				if err != nil {
					fmt.Printf("Error parsing reminder data for user %s: %v\n", psid, err)
					utils.SendTextMessage("There was an issue parsing your reminder. Please ensure you're using the format: [amount] to [name]:[gcash number] on [month/day/year]", psid, token)
					userState[psid] = "WAITING..." // Reset state as format was correct
					return
				}

				err = api.SaveReminder(psid, amount, accountName, gcashNumber, dueDate, "Gcash", "pending")
				if err != nil {
					fmt.Printf("Error saving reminder for user %s: %v\n", psid, err)
					utils.SendTextMessage("Sorry, I couldn't save your reminder. Please try again later.", psid, token)
					userState[psid] = "WAITING..." // Reset state as format was correct
					return
				}
				message_ := fmt.Sprintf("Reminder: Pay ₱%.2f to %s (%s) on %s", amount, accountName, gcashNumber, dueDate.Format("01/02/2006"))
				utils.SendTextMessage(message_, psid, token)
				fmt.Printf("Reminder saved for user %s: ₱%.2f to %s (%s) on %s\n", psid, amount, accountName, gcashNumber, dueDate.Format("01/02/2006"))
				userState[psid] = "WAITING..."
			} else {
				message = "Invalid Format. Follow the format or verify the Gcash Number. Please try again"
				utils.SendTextMessage(message, psid, token)
				ProcessMainCommand("GET_STARTED", psid, mid, token)
				userState[psid] = "WAITING..."
			}
		default:
			message = "Your input cannot be processed. Please select an option from the menu."
			utils.SendTextMessage(message, psid, token)
			utils.SendGenerateRequest(templates.MenuTemplate[1], psid, token)
		}
	} else {
		message = "Your input cannot be processed. Please select an option from the menu."
		utils.SendTextMessage(message, psid, token)
		utils.SendGenerateRequest(templates.MenuTemplate[1], psid, token)
	}

}
