package services

import (
	"fmt"
	"quickyexpensetracker/api"
	"quickyexpensetracker/templates"
	"quickyexpensetracker/utils"
	"time"
)

var userState = make(map[string]string)
var currentTime = time.Now()

func ProcessMainCommand(command, psid, mid, token string) {
	fmt.Printf("Processing Command: %s, PSID: %s, MID: %s\n", command, psid, mid)
	fmt.Printf("token: %s\n", token)

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
		fmt.Printf("Unknown command: %s\n", command)
		utils.SendGenerateRequest(templates.MenuTemplate[1], psid, token)
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
			fmt.Println("Error fetching items")
			return
		}

		report := utils.GetExpenseReport(expenses, "Daily")
		utils.SendTextMessage(report, psid, token)
	case "REPORT_LOG_WEEK":
		expenses, err := api.GetExpensesByUserAndRange(psid, "week")
		if err != nil {
			fmt.Println("Error fetching items")
			return
		}

		report := utils.GetExpenseReport(expenses, "Weekly")
		utils.SendTextMessage(report, psid, token)
	case "REPORT_LOG_MONTH":
		expenses, err := api.GetExpensesByUserAndRange(psid, "month")
		if err != nil {
			fmt.Println("Error fetching items")
			return
		}

		report := utils.GetExpenseReport(expenses, "Monthly")
		utils.SendTextMessage(report, psid, token)
	case "VIEW_PENDING_PAYMENTS_MESSAGE":
		reminders, err := api.GetReminders(psid, "pending")
		if err != nil {
			fmt.Println("Error fetching items")
			return
		}
		report := utils.GetRemindersReport(reminders)
		utils.SendTextMessage(report, psid, token)
	case "VIEW_ACCOMPLISHED_PAYMENTS_MESSAGE":
		reminders, err := api.GetReminders(psid, "completed")
		if err != nil {
			fmt.Println("Error fetching items")
			return
		}
		report := utils.GetRemindersReport(reminders)
		utils.SendTextMessage(report, psid, token)
	case "SET_REPORT_SCHED_MESSAGE":
		message := "You can now receive notifications Daily/Weekly/Monthly."
		utils.SendTextMessage(message, psid, token)
	case "RESET_LOGS_MESSAGE":
		message := "Logs have been permanently resetted"
		utils.SendTextMessage(message, psid, token)
	case "SET_REMINDER_MESSAGE":
		message := "Please set the reminder in this format: \n[amount] to [name]:[gcash number] on [month/day/year]\n(e.g. 200.00 to mark:09565546*** on 04/25/2025)"
		utils.SendTextMessage(message, psid, token)
		userState[psid] = "RECORDING_REMINDER"
	case "SUBSCRIPTION_STATUS_MESSAGE":
		message := "You can now receive notifications."
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
				message_ := fmt.Sprintf("Got it! You spent ₱%v on %v", message, currentTime.Format("Mon Jan 2 15:04:05 MST 2006"))
				utils.SendTextMessage(message_, psid, token)
				amount, category, err := utils.GetExpenseDataFromMessage(message)
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Printf("You spent ₱%.2f on %s\n", amount, category)
					err := api.SaveExpense(amount, category, psid)
					if err != nil {
						fmt.Println("Failed to save expense:", err)
					} else {
						fmt.Println("Expense saved successfully!")
					}
				}
			} else {
				message = "Invalid Format. Please try again."
				utils.SendTextMessage(message, psid, token)
				ProcessMainCommand("GET_STARTED", psid, mid, token)
				userState[psid] = "WAITING..."
			}
		case "RECORDING_REMINDER":
			if utils.IsReminderLogFormatCorrect(message) {
				message_ := fmt.Sprintf("Reminder: Pay %v", message)
				utils.SendTextMessage(message_, psid, token)
				amount, accountName, gcashNumber, dueDate, err := utils.GetReminderDataFromMessage(message)
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Printf("You need to send ₱%.2f to %s (%s) on %s\n", amount, accountName, gcashNumber, dueDate.Format("01/02/2006"))
					err := api.SaveReminder(psid, amount, accountName, gcashNumber, dueDate)
					if err != nil {
						fmt.Println("Failed to save reminder:", err)
					} else {
						fmt.Println("Reminder saved successfully! You will be notified.")
					}
				}
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
