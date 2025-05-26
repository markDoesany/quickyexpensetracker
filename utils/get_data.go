package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetExpenseDataFromMessage(message string) (amount float64, category string, err error) {
	parts := strings.Split(message, "for")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid format, expected: [amount] for [item]")
		return
	}

	amountString := strings.TrimSpace(parts[0])
	category = strings.TrimSpace(parts[1])

	amount, err = strconv.ParseFloat(amountString, 64)
	if err != nil {
		return
	}

	return
}

func GetReminderDataFromMessage(message string) (amount float64, accountName string, gcashNumber string, date time.Time, err error) {
	parts := strings.Split(message, " to ")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid format: missing 'to'")
		return
	}
	amountString := strings.TrimSpace(parts[0])
	amount, err = strconv.ParseFloat((amountString), 64)
	if err != nil {
		err = fmt.Errorf("invalid amount format")
		return
	}

	secondParts := strings.Split(parts[1], " on ")
	if len(secondParts) != 2 {
		err = fmt.Errorf("invalid format: missing 'on'")
		return
	}

	nameAndNumber := strings.Split(secondParts[0], ":")
	if len(nameAndNumber) != 2 {
		err = fmt.Errorf("invalid format: missing ':' between name and number")
		return
	}

	accountName = strings.TrimSpace(nameAndNumber[0])
	gcashNumber = strings.TrimSpace(nameAndNumber[1])

	date, err = time.Parse("01/02/2006", strings.TrimSpace(secondParts[1]))
	if err != nil {
		err = fmt.Errorf("invalid date format, expected MM/DD/YYYY")
		return
	}

	return
}
