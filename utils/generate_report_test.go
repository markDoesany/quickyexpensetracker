package utils

import (
	"quickyexpensetracker/models" // Assuming models.ExpensesLog is defined here
	"testing"
	// "time" // Not needed for this specific test, but often is for model creation
)

func TestGenerateExpenseSummaryMessage(t *testing.T) {
	tests := []struct {
		name          string
		frequency     string
		expenses      []models.ExpensesLog
		totalAmount   float64
		expectedMsg   string
	}{
		{
			name:      "daily summary with expenses",
			frequency: "daily",
			expenses: []models.ExpensesLog{
				{Amount: 100.50, Category: "Food", UserID: "user1"},
				{Amount: 50.00, Category: "Transport", UserID: "user1"},
			},
			totalAmount: 150.50,
			expectedMsg: "Hi! Here's your daily expense summary: You had 2 transaction(s), totaling ₱150.50.",
		},
		{
			name:      "weekly summary with no expenses",
			frequency: "weekly",
			expenses:  []models.ExpensesLog{},
			totalAmount: 0.0,
			expectedMsg: "Hi! You had no expenses in the last weekly.",
		},
		{
			name:      "monthly summary with one expense",
			frequency: "Monthly", // Test mixed case frequency
			expenses: []models.ExpensesLog{
				{Amount: 1234.56, Category: "Bills", UserID: "user2"},
			},
			totalAmount: 1234.56,
			expectedMsg: "Hi! Here's your monthly expense summary: You had 1 transaction(s), totaling ₱1234.56.",
		},
		{
			name:      "empty frequency with expenses",
			frequency: "",
			expenses: []models.ExpensesLog{
				{Amount: 10.00, Category: "Snacks", UserID: "user3"},
			},
			totalAmount: 10.00,
			expectedMsg: "Hi! Here's your selected period expense summary: You had 1 transaction(s), totaling ₱10.00.",
		},
		{
			name:      "empty frequency with no expenses",
			frequency: "",
			expenses:  []models.ExpensesLog{},
			totalAmount: 0.0,
			expectedMsg: "Hi! You had no expenses in the last selected period.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := GenerateExpenseSummaryMessage(tt.frequency, tt.expenses, tt.totalAmount)
			if msg != tt.expectedMsg {
				t.Errorf("Expected message:\n%s\nGot message:\n%s", tt.expectedMsg, msg)
			}
		})
	}
}
