package main

import (
	"log"
	"quickyexpensetracker/database"
	"quickyexpensetracker/handlers"
	"quickyexpensetracker/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Err loading .env file: %v", err)
	}
	database.InitDB()

	c := cron.New()
	_, err = c.AddFunc("@daily", services.SendDueReminderNotifications)
	if err != nil {
		log.Fatalf("Could not add 'SendDueReminderNotifications' cron job: %v", err)
	}
	c.Start()
	log.Println("Cron scheduler started.")

	router := gin.Default()
	router.GET("/", handlers.HandleVerification)
	router.POST("/", handlers.HandleWebhook)

	router.Run(":8080")
}
