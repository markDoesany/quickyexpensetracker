package utils

import (
	"fmt"
	"quickyexpensetracker/models"
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
	var elements []templates.Template

	if len(reminders) == 0 {
		return elements // Return empty slice
	}

	for _, reminder := range reminders {
		title := fmt.Sprintf("Pay ₱%.2f to %s", reminder.Amount, reminder.Rececipient)
		subtitle := fmt.Sprintf("Due: %s | Gcash: %s", reminder.DueDate.Format("Jan 2, 2006"), reminder.GcashNumber)
		// ImageURL can be omitted or set to a default Gcash logo if available
		element := templates.Template{
			Title:    title,
			Subtitle: subtitle,
			Buttons: []templates.Button{
				{
					Type:  "web_url",
					URL:   "gcash://", // Generic GCash deep link
					Title: "Pay with Gcash",
				},
			},
		}
		elements = append(elements, element)
	}

	return elements
}
