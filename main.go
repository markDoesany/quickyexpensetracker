package main

import (
	"fmt"
	"log"
	"time" // Added for ticker

	"quickyexpensetracker/database"
	"quickyexpensetracker/handlers"
	"quickyexpensetracker/services" // Added for reminder processor

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Err loading .env file: %v", err)
	}
	database.InitDB()

	// Start the reminder processor
	go func() {
		// Run once immediately at startup, then tick.
		fmt.Println("Starting initial check for due reminders...")
		services.CheckDueReminders()

		// Then, check periodically.
		// For example, check every 1 hour. Adjust the duration as needed.
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop() // Ensure the ticker is stopped if the goroutine exits

		for range ticker.C {
			fmt.Println("Periodic check for due reminders triggered by ticker...")
			services.CheckDueReminders()
		}
	}()

	router := gin.Default()
	router.GET("/", handlers.HandleVerification)
	router.POST("/", handlers.HandleWebhook)

	router.Run(":8080")
}
