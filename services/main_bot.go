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
	case "SCHEDULE_REPORT_DAILY_SETUP":
		ProcessTextMessageSent("SCHEDULE_REPORT_DAILY_SETUP_MESSAGE", psid, mid, token)
	case "SCHEDULE_REPORT_WEEKLY_SETUP":
		ProcessTextMessageSent("SCHEDULE_REPORT_WEEKLY_SETUP_MESSAGE", psid, mid, token)
	case "SCHEDULE_REPORT_MONTHLY_SETUP":
		ProcessTextMessageSent("SCHEDULE_REPORT_MONTHLY_SETUP_MESSAGE", psid, mid, token)
	case "VIEW_SCHEDULED_REPORT":
		ProcessTextMessageSent("VIEW_SCHEDULED_REPORT_MESSAGE", psid, mid, token)
	case "UNSCHEDULE_REPORTS":
		ProcessTextMessageSent("UNSCHEDULE_REPORTS_MESSAGE", psid, mid, token)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		utils.SendGenerateRequest(templates.MenuTemplate[1], psid, token)
	}
}

// formatScheduleForDisplay formats the schedule information for user-friendly display.
func formatScheduleForDisplay(schedule *models.ReportScheduleLog) string {
	if schedule == nil {
		return "You currently have no reports scheduled."
	}
	var dayInfo string
	switch schedule.Frequency {
	case "daily":
		dayInfo = "every day"
	case "weekly":
		dayInfo = fmt.Sprintf("every %s", time.Weekday(schedule.DayOfWeek).String())
	case "monthly":
		dayInfo = fmt.Sprintf("on day %d of the month", schedule.DayOfMonth)
	default:
		dayInfo = "at an unknown frequency"
	}

	lastSent := "Never"
	if !schedule.LastSentAt.IsZero() {
		lastSent = schedule.LastSentAt.Format("Jan 2, 2006 at 3:04 PM MST") // Consider using schedule.Timezone for display
	}

	return fmt.Sprintf("Your %s report is scheduled %s at %s (%s). Last sent: %s",
		schedule.Frequency, dayInfo, schedule.ScheduledTime, schedule.Timezone, lastSent)
}

func ProcessTextMessageSent(command, psid, mid, token string) {
	var messageText string // Used for simple text messages to avoid conflict with 'message' from ProcessTextMessageReceived
	switch command {
	case "LOG_EXPENSE_MESSAGE":
		messageText = "Please log in this format: \n[amount] for [item/service]\n(e.g. 200.00 for softdrinks)"
		utils.SendTextMessage(messageText, psid, token)
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

		reminderElements := utils.GetRemindersReport(reminders)
		if len(reminderElements) == 0 {
			utils.SendTextMessage("No pending reminders found.", psid, token)
		} else {
			err := utils.SendGenerateRequest(reminderElements, psid, token)
			if err != nil {
				fmt.Printf("Error sending reminder elements for user %s: %v\n", psid, err)
				utils.SendTextMessage("Sorry, there was an issue displaying your reminders. Please try again.", psid, token)
			}
		}
	case "VIEW_ACCOMPLISHED_PAYMENTS_MESSAGE":
		reminders, err := api.GetReminders(psid, "completed")
		if err != nil {
			fmt.Printf("Error fetching accomplished reminders for user %s: %v\n", psid, err)
			utils.SendTextMessage("Sorry, I couldn't fetch your payment reminders at the moment. Please try again later.", psid, token)
			return
		}
		// utils.GetRemindersReport returns []templates.Template. This case needs a string for accomplished.
		var accomplishedReportStr string
		if len(reminders) == 0 {
			accomplishedReportStr = "No accomplished payments found."
		} else {
			accomplishedReportStr = "Accomplished Payments:\n"
			for _, r := range reminders { // Assuming 'reminders' here is []models.RemindersLog
				accomplishedReportStr += fmt.Sprintf("- ₱%.2f to %s on %s\n", r.Amount, r.Rececipient, r.DueDate.Format("Jan 2, 2006"))
			}
		}
		utils.SendTextMessage(accomplishedReportStr, psid, token)
	case "SET_REPORT_SCHED_MESSAGE": // This case can be considered deprecated or a fallback
		messageText = "Please use the 'Manage Report Schedule' option in the menu to set your report schedule."
		utils.SendTextMessage(messageText, psid, token)
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
		schedule, err := api.GetSchedule(psid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendTextMessage("You are not currently subscribed to any scheduled reports. Use the 'Manage Report Schedule' menu to set one up.", psid, token)
			} else {
				fmt.Printf("Error fetching schedule for user %s: %v\n", psid, err)
				utils.SendTextMessage("Sorry, I couldn't retrieve your subscription status at the moment.", psid, token)
			}
			return
		}

		if schedule != nil && schedule.ID != 0 { // Check if a valid schedule was found
			// Use the existing formatScheduleForDisplay helper
			formattedMsg := formatScheduleForDisplay(schedule)
			utils.SendTextMessage(formattedMsg, psid, token)
		} else {
			// This case should ideally be covered by gorm.ErrRecordNotFound
			utils.SendTextMessage("You are not currently subscribed to any scheduled reports. Use the 'Manage Report Schedule' menu to set one up.", psid, token)
		}

	// New cases for report scheduling
	case "SCHEDULE_REPORT_DAILY_SETUP_MESSAGE":
		utils.SendTextMessage("What time (e.g., 09:00 in HH:MM format) and timezone (e.g., Asia/Manila) would you like to receive daily reports? Format: HH:MM Your/Timezone", psid, token)
		userState[psid] = "AWAITING_DAILY_SCHEDULE_TIME_ZONE"
	case "SCHEDULE_REPORT_WEEKLY_SETUP_MESSAGE":
		utils.SendTextMessage("Which day (e.g., Monday), time (HH:MM), and timezone (e.g., Asia/Manila) for weekly reports? Format: Day HH:MM Your/Timezone", psid, token)
		userState[psid] = "AWAITING_WEEKLY_SCHEDULE_DAY_TIME_ZONE"
	case "SCHEDULE_REPORT_MONTHLY_SETUP_MESSAGE":
		utils.SendTextMessage("Which day of the month (1-31), time (HH:MM), and timezone (e.g., Asia/Manila) for monthly reports? Format: DD HH:MM Your/Timezone", psid, token)
		userState[psid] = "AWAITING_MONTHLY_SCHEDULE_DAY_TIME_ZONE"
	case "VIEW_SCHEDULED_REPORT_MESSAGE":
		schedule, err := api.GetSchedule(psid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendTextMessage("You currently have no reports scheduled.", psid, token)
			} else {
				fmt.Printf("Error fetching schedule for user %s: %v\n", psid, err)
				utils.SendTextMessage("Sorry, I couldn't retrieve your schedule information. Please try again later.", psid, token)
			}
			return
		}
		utils.SendTextMessage(formatScheduleForDisplay(schedule), psid, token)
	case "UNSCHEDULE_REPORTS_MESSAGE":
		err := api.DeleteSchedule(psid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendTextMessage("You had no reports scheduled, so nothing was changed.", psid, token)
			} else {
				fmt.Printf("Error deleting schedule for user %s: %v\n", psid, err)
				utils.SendTextMessage("Could not remove schedule. Please try again.", psid, token)
			}
		} else {
			utils.SendTextMessage("Your report schedule has been removed.", psid, token)
		}
	default:
		fmt.Printf("Unknown command in ProcessTextMessageSent: %s\n", command)
		// utils.SendTextMessage("Sorry, I didn't understand that command.", psid, token)
	}
}

// formatScheduleForDisplay formats the schedule information for user-friendly display.
func formatScheduleForDisplay(schedule *models.ReportScheduleLog) string {
	if schedule == nil {
		return "You currently have no reports scheduled."
	}
	var dayInfo string
	switch schedule.Frequency {
	case "daily":
		dayInfo = "every day"
	case "weekly":
		dayInfo = fmt.Sprintf("every %s", time.Weekday(schedule.DayOfWeek).String())
	case "monthly":
		dayInfo = fmt.Sprintf("on day %d of the month", schedule.DayOfMonth)
	default:
		dayInfo = "at an unknown frequency"
	}

	lastSent := "Never"
	if !schedule.LastSentAt.IsZero() { // Check if LastSentAt is not the zero value
		// Ensure LastSentAt is in local time or a consistent timezone for display if needed
		// For now, assuming it's stored in UTC and displayed as such or with MST (as per previous format)
		lastSent = schedule.LastSentAt.Format("Jan 2, 2006 at 3:04 PM MST")
	}

	return fmt.Sprintf("You are scheduled for %s reports at %s (%s). Last sent: %s",
		schedule.Frequency, dayInfo, schedule.ScheduledTime, schedule.Timezone, lastSent)
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
		reminderElements := utils.GetRemindersReport(reminders)
		if len(reminderElements) == 0 {
			utils.SendTextMessage("No pending reminders found.", psid, token)
		} else {
			err := utils.SendGenerateRequest(reminderElements, psid, token)
			if err != nil {
				fmt.Printf("Error sending reminder elements for user %s: %v\n", psid, err)
				utils.SendTextMessage("Sorry, there was an issue displaying your reminders. Please try again.", psid, token)
			}
		}
	case "VIEW_ACCOMPLISHED_PAYMENTS_MESSAGE":
		reminders, err := api.GetReminders(psid, "completed")
		if err != nil {
			fmt.Printf("Error fetching accomplished reminders for user %s: %v\n", psid, err)
			utils.SendTextMessage("Sorry, I couldn't fetch your payment reminders at the moment. Please try again later.", psid, token)
			return
		}
		// utils.GetRemindersReport returns []templates.Template. This case needs a string.
		var accomplishedReportStr string
		if len(reminders) == 0 {
			accomplishedReportStr = "No accomplished payments found."
		} else {
			accomplishedReportStr = "Accomplished Payments:\n"
			for _, r := range reminders { // Assuming 'reminders' here is []models.RemindersLog
				accomplishedReportStr += fmt.Sprintf("- ₱%.2f to %s on %s\n", r.Amount, r.Rececipient, r.DueDate.Format("Jan 2, 2006"))
			}
		}
		utils.SendTextMessage(accomplishedReportStr, psid, token)
	case "SET_REPORT_SCHED_MESSAGE": // This case can be considered deprecated or a fallback
		message := "Please use the 'Manage Report Schedule' option in the menu to set your report schedule."
		utils.SendTextMessage(message, psid, token)
	case "RESET_LOGS_MESSAGE":
		errExpenses := api.DeleteExpensesByUser(psid)
		if errExpenses != nil {
			fmt.Printf("Error deleting expenses for user %s: %v\n", psid, errExpenses)
		}
		errReminders := api.DeleteRemindersByUser(psid)
		if errReminders != nil {
			fmt.Printf("Error deleting reminders for user %s: %v\n", psid, errReminders)
		}
		utils.SendTextMessage("All your expense and reminder logs have been reset.", psid, token)
	case "SET_REMINDER_MESSAGE":
		message := "Please set the reminder in this format: \n[amount] to [name]:[gcash number] on [month/day/year]\n(e.g. 200.00 to mark:09565546*** on 04/25/2025)"
		utils.SendTextMessage(message, psid, token)
		userState[psid] = "RECORDING_REMINDER"
	case "SUBSCRIPTION_STATUS_MESSAGE": // This case can be considered deprecated or a fallback
		message := "Subscription features are not yet fully implemented. Please check back later!"
		utils.SendTextMessage(message, psid, token)

	// New cases for report scheduling
	case "SCHEDULE_REPORT_DAILY_SETUP_MESSAGE":
		utils.SendTextMessage("What time (e.g., 09:00 in HH:MM format) and timezone (e.g., Asia/Manila) would you like to receive daily reports? Format: HH:MM Your/Timezone", psid, token)
		userState[psid] = "AWAITING_DAILY_SCHEDULE_TIME_ZONE"
	case "SCHEDULE_REPORT_WEEKLY_SETUP_MESSAGE":
		utils.SendTextMessage("Which day (e.g., Monday), time (HH:MM), and timezone (e.g., Asia/Manila) for weekly reports? Format: Day HH:MM Your/Timezone", psid, token)
		userState[psid] = "AWAITING_WEEKLY_SCHEDULE_DAY_TIME_ZONE"
	case "SCHEDULE_REPORT_MONTHLY_SETUP_MESSAGE":
		utils.SendTextMessage("Which day of the month (1-31), time (HH:MM), and timezone (e.g., Asia/Manila) for monthly reports? Format: DD HH:MM Your/Timezone", psid, token)
		userState[psid] = "AWAITING_MONTHLY_SCHEDULE_DAY_TIME_ZONE"
	case "VIEW_SCHEDULED_REPORT_MESSAGE":
		schedule, err := api.GetSchedule(psid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendTextMessage("You currently have no reports scheduled.", psid, token)
			} else {
				fmt.Printf("Error fetching schedule for user %s: %v\n", psid, err)
				utils.SendTextMessage("Sorry, I couldn't retrieve your schedule information. Please try again later.", psid, token)
			}
			return
		}
		utils.SendTextMessage(formatScheduleForDisplay(schedule), psid, token)
	case "UNSCHEDULE_REPORTS_MESSAGE":
		err := api.DeleteSchedule(psid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendTextMessage("You had no reports scheduled, so nothing was changed.", psid, token)
			} else {
				fmt.Printf("Error deleting schedule for user %s: %v\n", psid, err)
				utils.SendTextMessage("Could not remove schedule. Please try again.", psid, token)
			}
		} else {
			utils.SendTextMessage("Your report schedule has been removed.", psid, token)
		}
	default:
		fmt.Printf("Unknown command in ProcessTextMessageSent: %s\n", command)
		// utils.SendTextMessage("Sorry, I didn't understand that command.", psid, token)
	}
}

func formatScheduleForDisplay(schedule *models.ReportScheduleLog) string {
	if schedule == nil {
		return "You currently have no reports scheduled."
	}
	var dayInfo string
	switch schedule.Frequency {
	case "daily":
		dayInfo = "every day"
	case "weekly":
		dayInfo = fmt.Sprintf("every %s", time.Weekday(schedule.DayOfWeek).String())
	case "monthly":
		dayInfo = fmt.Sprintf("on day %d of the month", schedule.DayOfMonth)
	default:
		dayInfo = "at an unknown frequency"
	}

	lastSent := "Never"
	if !schedule.LastSentAt.IsZero() {
		lastSent = schedule.LastSentAt.Format("Jan 2, 2006 at 3:04 PM MST")
	}

	return fmt.Sprintf("You are scheduled for %s reports at %s in %s. Last sent: %s",
		schedule.Frequency, schedule.ScheduledTime, schedule.Timezone, lastSent)
}

func ProcessTextMessageSent(command, psid, mid, token string) {
	switch command {
	// ... (previous cases remain unchanged)
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

		reminderElements := utils.GetRemindersReport(reminders)
		if len(reminderElements) == 0 {
			utils.SendTextMessage("No pending reminders found.", psid, token)
		} else {
			err := utils.SendGenerateRequest(reminderElements, psid, token)
			if err != nil {
				fmt.Printf("Error sending reminder elements for user %s: %v\n", psid, err)
				utils.SendTextMessage("Sorry, there was an issue displaying your reminders. Please try again.", psid, token)
			}
		}
	case "VIEW_ACCOMPLISHED_PAYMENTS_MESSAGE":
		reminders, err := api.GetReminders(psid, "completed")
		if err != nil {
			fmt.Printf("Error fetching accomplished reminders for user %s: %v\n", psid, err)
			utils.SendTextMessage("Sorry, I couldn't fetch your payment reminders at the moment. Please try again later.", psid, token)
			return
		}
		// Assuming GetRemindersReport returns []templates.Template, but this command sends text.
		// For simplicity, let's assume a text version or adapt GetRemindersReport if needed.
		// This part of the code expects GetRemindersReport to return a string for accomplished payments.
		// However, GetRemindersReport was changed to return []templates.Template.
		// This will cause a compile error.
		// For now, I will create a simple text report for accomplished payments here.
		var accomplishedReportStr string
		if len(reminders) == 0 {
			accomplishedReportStr = "No accomplished payments found."
		} else {
			accomplishedReportStr = "Accomplished Payments:\n"
			for _, r := range reminders {
				accomplishedReportStr += fmt.Sprintf("- ₱%.2f to %s on %s\n", r.Amount, r.Rececipient, r.DueDate.Format("Jan 2, 2006"))
			}
		}
		utils.SendTextMessage(accomplishedReportStr, psid, token)

	case "SET_REPORT_SCHED_MESSAGE": // This case might be deprecated by the new submenu
		message := "The report scheduling feature is not yet implemented. Please check back later!"
		utils.SendTextMessage(message, psid, token)
	case "RESET_LOGS_MESSAGE":
		errExpenses := api.DeleteExpensesByUser(psid)
		if errExpenses != nil {
			fmt.Printf("Error deleting expenses for user %s: %v\n", psid, errExpenses)
		}
		errReminders := api.DeleteRemindersByUser(psid)
		if errReminders != nil {
			fmt.Printf("Error deleting reminders for user %s: %v\n", psid, errReminders)
		}
		message := "All your expense and reminder logs have been reset."
		utils.SendTextMessage(message, psid, token)
	case "SET_REMINDER_MESSAGE":
		message := "Please set the reminder in this format: \n[amount] to [name]:[gcash number] on [month/day/year]\n(e.g. 200.00 to mark:09565546*** on 04/25/2025)"
		utils.SendTextMessage(message, psid, token)
		userState[psid] = "RECORDING_REMINDER"
	case "SUBSCRIPTION_STATUS_MESSAGE": // This case might be deprecated
		message := "The subscription status feature is not yet implemented. Please check back later!"
		utils.SendTextMessage(message, psid, token)

	// New cases for report scheduling
	case "SCHEDULE_REPORT_DAILY_SETUP_MESSAGE":
		utils.SendTextMessage("What time (e.g., 09:00 in HH:MM format) and timezone (e.g., Asia/Manila) would you like to receive daily reports? Format: HH:MM Your/Timezone", psid, token)
		userState[psid] = "AWAITING_DAILY_SCHEDULE_TIME_ZONE"
	case "SCHEDULE_REPORT_WEEKLY_SETUP_MESSAGE":
		utils.SendTextMessage("Which day (e.g., Monday), time (HH:MM), and timezone (e.g., Asia/Manila) for weekly reports? Format: Day HH:MM Your/Timezone", psid, token)
		userState[psid] = "AWAITING_WEEKLY_SCHEDULE_DAY_TIME_ZONE"
	case "SCHEDULE_REPORT_MONTHLY_SETUP_MESSAGE":
		utils.SendTextMessage("Which day of the month (1-31), time (HH:MM), and timezone (e.g., Asia/Manila) for monthly reports? Format: DD HH:MM Your/Timezone", psid, token)
		userState[psid] = "AWAITING_MONTHLY_SCHEDULE_DAY_TIME_ZONE"
	case "VIEW_SCHEDULED_REPORT_MESSAGE":
		schedule, err := api.GetSchedule(psid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendTextMessage("You currently have no reports scheduled.", psid, token)
			} else {
				fmt.Printf("Error fetching schedule for user %s: %v\n", psid, err)
				utils.SendTextMessage("Sorry, I couldn't retrieve your schedule information. Please try again later.", psid, token)
			}
			return
		}
		utils.SendTextMessage(formatScheduleForDisplay(schedule), psid, token)
	case "UNSCHEDULE_REPORTS_MESSAGE":
		err := api.DeleteSchedule(psid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendTextMessage("You had no reports scheduled, so nothing was changed.", psid, token)
			} else {
				fmt.Printf("Error deleting schedule for user %s: %v\n", psid, err)
				utils.SendTextMessage("Could not remove schedule. Please try again.", psid, token)
			}
		} else {
			utils.SendTextMessage("Your report schedule has been removed.", psid, token)
		}
	default:
		fmt.Printf("Unknown command in ProcessTextMessageSent: %s\n", command)
		// utils.SendTextMessage("Sorry, I didn't understand that command.", psid, token)
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

				err = api.SaveReminder(psid, amount, accountName, gcashNumber, dueDate)
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

// Helper function to parse "HH:MM Your/Timezone"
func parseTimeAndZone(text string) (string, string, error) {
	parts := strings.Fields(text)
	if len(parts) != 2 {
		return "", "", errors.New("invalid format. Expected: HH:MM Your/Timezone")
	}
	// Basic validation for HH:MM can be added here if needed
	// e.g., using a regex like in api/schedule.go
	timeRegex := regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`)
	if !timeRegex.MatchString(parts[0]) {
		return "", "", errors.New("invalid time format. Expected HH:MM")
	}
	// Basic timezone validation (presence)
	if parts[1] == "" {
		return "", "", errors.New("timezone cannot be empty")
	}
	// More robust timezone validation (e.g., time.LoadLocation) can be added if needed.
	return parts[0], parts[1], nil
}

// Helper function to parse "Day HH:MM Your/Timezone"
func parseDayTimeAndZone(text string) (int, string, string, error) {
	parts := strings.Fields(text)
	if len(parts) != 3 {
		return 0, "", "", errors.New("invalid format. Expected: Day HH:MM Your/Timezone")
	}
	dayStr := strings.ToLower(parts[0])
	var dayOfWeek int
	switch dayStr {
	case "sunday":
		dayOfWeek = 0
	case "monday":
		dayOfWeek = 1
	case "tuesday":
		dayOfWeek = 2
	case "wednesday":
		dayOfWeek = 3
	case "thursday":
		dayOfWeek = 4
	case "friday":
		dayOfWeek = 5
	case "saturday":
		dayOfWeek = 6
	default:
		return 0, "", "", errors.New("invalid day of the week")
	}

	timeRegex := regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`)
	if !timeRegex.MatchString(parts[1]) {
		return 0, "", "", errors.New("invalid time format. Expected HH:MM")
	}
	if parts[2] == "" {
		return 0, "", "", errors.New("timezone cannot be empty")
	}
	return dayOfWeek, parts[1], parts[2], nil
}

// Helper function to parse "DD HH:MM Your/Timezone"
func parseDayOfMonthTimeAndZone(text string) (int, string, string, error) {
	parts := strings.Fields(text)
	if len(parts) != 3 {
		return 0, "", "", errors.New("invalid format. Expected: DD HH:MM Your/Timezone")
	}
	dayOfMonth, err := strconv.Atoi(parts[0])
	if err != nil || dayOfMonth < 1 || dayOfMonth > 31 {
		return 0, "", "", errors.New("invalid day of the month. Expected 1-31")
	}

	timeRegex := regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`)
	if !timeRegex.MatchString(parts[1]) {
		return 0, "", "", errors.New("invalid time format. Expected HH:MM")
	}
	if parts[2] == "" {
		return 0, "", "", errors.New("timezone cannot be empty")
	}
	return dayOfMonth, parts[1], parts[2], nil
}
