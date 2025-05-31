package utils

import (
	"testing"
	"time"
)

func TestCalculateNextDueDate(t *testing.T) {
	// Test cases
	tests := []struct {
		name             string
		currentDueDate   time.Time
		frequency        string
		expectedNextDate time.Time
		expectError      bool
		expectedHour     int
		expectedMinute   int
		expectedSecond   int
	}{
		{
			name:             "daily frequency",
			currentDueDate:   time.Date(2023, 10, 26, 10, 30, 0, 0, time.UTC),
			frequency:        "daily",
			expectedNextDate: time.Date(2023, 10, 27, 0, 0, 0, 0, time.UTC),
			expectError:      false,
			expectedHour:     0,
			expectedMinute:   0,
			expectedSecond:   0,
		},
		{
			name:             "weekly frequency",
			currentDueDate:   time.Date(2023, 10, 26, 15, 0, 0, 0, time.UTC),
			frequency:        "weekly",
			expectedNextDate: time.Date(2023, 11, 2, 0, 0, 0, 0, time.UTC),
			expectError:      false,
			expectedHour:     0,
			expectedMinute:   0,
			expectedSecond:   0,
		},
		{
			name:             "monthly frequency",
			currentDueDate:   time.Date(2023, 10, 26, 0, 0, 0, 0, time.UTC),
			frequency:        "monthly",
			expectedNextDate: time.Date(2023, 11, 26, 0, 0, 0, 0, time.UTC),
			expectError:      false,
			expectedHour:     0,
			expectedMinute:   0,
			expectedSecond:   0,
		},
		{
			name:           "monthly frequency - end of month",
			currentDueDate: time.Date(2023, 1, 31, 12, 0, 0, 0, time.UTC),
			frequency:      "monthly",
			// Adding one month to Jan 31 should result in Feb 28 (or 29 in a leap year)
			// For 2023, it's Feb 28. Time should be reset.
			expectedNextDate: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			expectError:      false,
			expectedHour:     0,
			expectedMinute:   0,
			expectedSecond:   0,
		},
		{
			name:           "once frequency - should error",
			currentDueDate: time.Date(2023, 10, 26, 0, 0, 0, 0, time.UTC),
			frequency:      "once",
			expectError:    true,
		},
		{
			name:           "unknown frequency - should error",
			currentDueDate: time.Date(2023, 10, 26, 0, 0, 0, 0, time.UTC),
			frequency:      "yearly",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextDate, err := CalculateNextDueDate(tt.currentDueDate, tt.frequency)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error for frequency '%s', but got none", tt.frequency)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for frequency '%s', but got: %v", tt.frequency, err)
				}
				if !nextDate.Equal(tt.expectedNextDate) {
					t.Errorf("Expected next date %v, but got %v", tt.expectedNextDate, nextDate)
				}
				if nextDate.Hour() != tt.expectedHour || nextDate.Minute() != tt.expectedMinute || nextDate.Second() != tt.expectedSecond {
					t.Errorf("Expected time to be reset to %02d:%02d:%02d, but got %02d:%02d:%02d",
						tt.expectedHour, tt.expectedMinute, tt.expectedSecond,
						nextDate.Hour(), nextDate.Minute(), nextDate.Second())
				}
			}
		})
	}
}
