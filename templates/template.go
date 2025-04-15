package templates

var Templates = map[int]Template{
	1: {
		Title:    "QuickyTracker",
		Subtitle: "What quicky transaction do you want to do?",
		// ImageURL: "https://goodkredit-final.s3.us-west-1.amazonaws.com/lokalwifi.jpg",
		Buttons: []Button{
			{Type: "postback", Title: "Log Expenses", Payload: "LOG_EXPENSES_MENU"},
			{Type: "postback", Title: "Remind Payments", Payload: "REMIND_PAYMENTS_MENU"},
			{Type: "postback", Title: "Subscribe", Payload: "SUBSCRIPTION_STATUS"},
		},
	},
	2: {
		Title:    "This is your Account ID:\n<AccountID>",
		ImageURL: "https://goodkredit-final.s3.us-west-1.amazonaws.com/lokalwifi.jpg",
		Buttons: []Button{
			{Type: "postback", Title: "Menus", Payload: "SHOW_MENU-OPTIONS"},
			{Type: "postback", Title: "Buy LokalWiFi Credit", Payload: "BUY_CREDIT-VOUCHERS"},
			{Type: "web_url", Title: "Visit Website", URL: "https://lokalwifi.net"},
		},
	},
}
