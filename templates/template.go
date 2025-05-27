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
		Title:    "Manage Report Schedule", // Updated Title for clarity
		Subtitle: "Choose an action for your report schedule", // Updated Subtitle
		Buttons: []Button{
			{Type: "postback", Title: "Schedule Daily Report", Payload: "SCHEDULE_REPORT_DAILY_SETUP"},
			{Type: "postback", Title: "Schedule Weekly Report", Payload: "SCHEDULE_REPORT_WEEKLY_SETUP"},
			{Type: "postback", Title: "Schedule Monthly Report", Payload: "SCHEDULE_REPORT_MONTHLY_SETUP"},
			{Type: "postback", Title: "View Current Schedule", Payload: "VIEW_SCHEDULED_REPORT"},
			{Type: "postback", Title: "Unschedule Reports", Payload: "UNSCHEDULE_REPORTS"},
		},
	},
}
