package utils

import (
	"fmt"
	"quickyexpensetracker/models"
	"quickyexpensetracker/templates"
	"strings"
)

func GetExpenseReport(expenses []models.ExpensesLog, rangeDay string) string {
	var total float64 = 0
	categoryTotals := make(map[string]float64)

	// Sum total and category totals
	for _, exp := range expenses {
		total += exp.Amount
		categoryTotals[exp.Category] += exp.Amount
	}

	// Build the report
	report := fmt.Sprintf("%v Report\n", rangeDay)
	report += fmt.Sprintf("Total: %.2f\n", total)
	for category, amount := range categoryTotals {
		var percentage float64
		if total > 0 {
			percentage = (amount / total) * 100
		} else {
			percentage = 0 // Or handle as appropriate, e.g. display N/A
		}
		report += fmt.Sprintf("%s = ₱%.2f - %.2f%%\n", category, amount, percentage)
	}

	return report
}

func GetRemindersReport(reminders []models.RemindersLog) []templates.Template {
	var reportElements []templates.Template

	if len(reminders) == 0 {
		element := templates.Template{
			Title:    "Reminders Report",
			Subtitle: "No pending reminders found.",
		}
		reportElements = append(reportElements, element)
		return reportElements
	}

	for _, reminder := range reminders {
		title := fmt.Sprintf("Payment to %s", reminder.Recipient)
		subtitle := fmt.Sprintf("Amount: ₱%.2f\nGCash: %s\nDue: %s",
			reminder.Amount,
			reminder.GcashNumber,
			reminder.DueDate.Format("2006-01-02"))

		var buttons []templates.Button
		if reminder.PaymentMethod == "Gcash" && reminder.Status == "pending" {
			buttons = append(buttons, templates.Button{
				Type:    "postback",
				Title:   "Pay with Gcash",
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
		reportElements = append(reportElements, element)
	}

	return reportElements
}

// GenerateExpenseSummaryMessage creates a user-friendly expense summary message.
func GenerateExpenseSummaryMessage(frequency string, expenses []models.ExpensesLog, totalAmount float64) string {
	// Ensure frequency string is lowercase for consistent messaging if it's used directly.
	// Or, use a more display-friendly version if needed.
	displayFrequency := strings.ToLower(frequency)
	if displayFrequency == "" {
		displayFrequency = "selected period" // Fallback for empty frequency
	}

	if len(expenses) == 0 {
		return fmt.Sprintf("Hi! You had no expenses in the last %s.", displayFrequency)
	}

	return fmt.Sprintf("Hi! Here's your %s expense summary: You had %d transaction(s), totaling ₱%.2f.",
		displayFrequency, len(expenses), totalAmount)
}
