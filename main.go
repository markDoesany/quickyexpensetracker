package main

import (
	"log"
	"quickyexpensetracker/database"
	"quickyexpensetracker/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Err loading .env file: %v", err)
	}
	database.InitDB()

	router := gin.Default()
	router.GET("/", handlers.HandleVerification)
	router.POST("/webhook", handlers.HandleWebhook)

	router.Run(":8080")
}
