package utils

import "regexp"

func IsExpenseLogFormatCorrect(text string) bool {
	pattern := `(?i)^\s*(\d+(\.\d{1,2})?)\s+for\s+(.+?)\s*$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(text)
}

func IsReminderLogFormatCorrect(text string) bool {
	pattern := `(?i)^\s*(\d+(\.\d{1,2})?)\s+to\s+(.+?):(09\d{9})\s+on\s+(\d{2}/\d{2}/\d{4})\s*$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(text)
}
