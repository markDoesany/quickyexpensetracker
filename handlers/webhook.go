package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"quickyexpensetracker/services"
	"quickyexpensetracker/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserInstance struct {
	PSID    string
	Command string
	MID     string
	Source  string
}

func HandleVerification(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode != "subscribe" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mode"})
	}

	if token != os.Getenv("VERIFY_TOKEN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if challenge == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing challenge"})
		return
	}

	c.String(http.StatusOK, challenge)
}

func HandleWebhook(c *gin.Context) {
	log.Println("Received webhook event")
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading body"})
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if body["object"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid object"})
		return
	}
	if c.ContentType() != "application/json" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid content type"})
		return
	}

	c.String(http.StatusOK, "EVENT_RECEIVED")

	// Signature Validation
	fbSignature := c.GetHeader("X-Hub-Signature-256")
	if fbSignature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}

	key := os.Getenv("APP_SECRET")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing app secret"})
		return
	}

	hash := utils.ComputeHMAC(rawBody, key)
	mySignature := "sha256=" + hash
	if fbSignature != mySignature {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	entries, ok := body["entry"].([]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry"})
		return
	}

	var userInstances []UserInstance
	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		messages, ok := entryMap["messaging"].([]interface{})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid messaging"})
			continue
		}

		for _, msg := range messages {
			msgMap := msg.(map[string]interface{})
			user := UserInstance{}

			sender := msgMap["sender"].(map[string]interface{})
			psid, ok := sender["id"].(string)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PSID"})
				continue
			}
			user.PSID = psid
			var command string
			var mid string
			var source string

			if message, ok := msgMap["message"].(map[string]interface{}); ok {
				if val, ok := message["quick_reply"].(map[string]interface{}); ok {
					command, _ = val["payload"].(string)
					source = "PAYLOAD"
				}
				mid, _ = message["mid"].(string)

				if text, ok := message["text"].(string); ok {
					fmt.Printf("Received message from PSID %s: %s\n", psid, text)
					source = "MESSAGE"
				}
			}

			if command == "" {
				if postback, ok := msgMap["postback"].(map[string]interface{}); ok {
					command, _ = postback["payload"].(string)
					mid, _ = postback["mid"].(string)
					source = "PAYLOAD"
					fmt.Printf("Received postback from PSID %s: %s\n", psid, command)
				}
			}

			if command == "" {
				if message, ok := msgMap["message"].(map[string]interface{}); ok {
					command, _ = message["text"].(string)
					source = "MESSAGE"
				}
			}

			if command == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid command"})
				continue
			}

			fmt.Printf("Processing command from PSID %s: %s\n", psid, command)

			user.Command = command
			user.MID = mid
			user.Source = source
			userInstances = append(userInstances, user)
		}
	}

	pageAccessToken := os.Getenv("PAGE_TOKEN")
	if pageAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing page access token"})
		return
	}

	for _, user := range userInstances {
		if user.Source == "PAYLOAD" {
			user.Command = strings.ToUpper(user.Command)
			services.ProcessMainCommand(user.Command, user.PSID, user.MID, pageAccessToken)
		} else if user.Source == "MESSAGE" {
			services.ProcessTextMessageReceived(user.Command, user.PSID, user.MID, pageAccessToken)
		}
	}
}
