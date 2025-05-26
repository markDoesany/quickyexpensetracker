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

func GetRemindersReport(reminders []models.RemindersLog) string {
	report := "Reminders Report\n"
	report += "-----------------------------\n"

	if len(reminders) == 0 {
		report += "No reminders found.\n"
		return report
	}

	for i, reminder := range reminders {
		report += fmt.Sprintf(
			"%d. ₱%.2f to %s (Account Number: %s) - Due: %s\n",
			i+1,
			reminder.Amount,
			reminder.Rececipient,
			reminder.GcashNumber,
			reminder.DueDate.Format("2006-01-02"),
		)
	}

	return report
}
