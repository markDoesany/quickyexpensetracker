package templates

var MenuTemplate = map[int]Template{
	1: {
		Title:    "Welcome to Quick-E Tracker",
		Subtitle: "What quicky transaction do you want to do?",
		Buttons: []Button{
			{Type: "postback", Title: "Log Expenses", Payload: "LOG_EXPENSES_MENU"},
			{Type: "postback", Title: "Remind Payments", Payload: "REMIND_PAYMENTS_MENU"},
			{Type: "postback", Title: "Subscribe", Payload: "SUBSCRIPTION_STATUS"},
		},
	},
	2: {
		Title:    "Log Your Expenses",
		Subtitle: "What quicky expenses log do you want to do?",
		Buttons: []Button{
			{Type: "postback", Title: "Record Expenses", Payload: "LOG_EXPENSES"},
			{Type: "postback", Title: "Generate Report", Payload: "GENERATE_REPORT_SUBMENU"},
			{Type: "postback", Title: "More", Payload: "EXPAND_MENU"},
		},
	},
	3: {
		Title:    "Payments Reminder",
		Subtitle: "What quicky payment reminders do you want to do?",
		Buttons: []Button{
			{Type: "postback", Title: "Set Reminder", Payload: "SET_REMINDER"},
			{Type: "postback", Title: "View Pending Payments", Payload: "VIEW_PENDING_PAYMENTS"},
			{Type: "postback", Title: "View Accomplished Payments", Payload: "VIEW_ACCOMPLISHED_PAYMENTS"},
		},
	},
	4: {
		Title:    "Log Your Expenses",
		Subtitle: "What quicky expenses log do you want to do?",
		Buttons: []Button{
			{Type: "postback", Title: "Set Report Sched", Payload: "SET_REPORT_SCHED_SUBMENU"},
			{Type: "postback", Title: "Reset Sched", Payload: "RESET_SCHED"},
		},
	},
}

var SubMenuTemplate = map[int]Template{
	1: {
		Title:    "Generate Financial Report",
		Subtitle: "Generate reports for your expenses for the past..",
		Buttons: []Button{
			{Type: "postback", Title: "Day", Payload: "GENERATE_REPORT_DAY"},
			{Type: "postback", Title: "Week", Payload: "GENERATE_REPORT_WEEK"},
			{Type: "postback", Title: "Month", Payload: "GENERATE_REPORT_MONTH"},
		},
	},
	2: {
		Title:    "Set Report Schedule",
		Subtitle: "Choose when you want to receive your report",
		Buttons: []Button{
			{Type: "postback", Title: "Daily", Payload: "SET_REPORT_DAILY"},
			{Type: "postback", Title: "Weekly", Payload: "SET_REPORT_WEEKLY"},
			{Type: "postback", Title: "Monthly", Payload: "SET_REPORT_MONTHLY"},
		},
	},
}
