package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserInstance struct {
	PSID    string
	Command string
	MID     string
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

	rawBody, _ := c.GetRawData()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	hash := computeHMAC(rawBody, key)
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

			if message, ok := msgMap["message"].(map[string]interface{}); ok {
				if val, ok := message["quick_reply"].(map[string]interface{}); ok {
					command, _ = val["payload"].(string)
				}
				mid, _ = message["mid"].(string)
			}

			if command == "" {
				if postback, ok := msgMap["postback"].(map[string]interface{}); ok {
					command, _ = postback["payload"].(string)
					mid, _ = postback["mid"].(string)
				}
			}

			if command == "" {
				if message, ok := msgMap["message"].(map[string]interface{}); ok {
					command, _ = message["text"].(string)
				}
			}

			if command == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid command"})
				continue
			}

			user.Command = strings.ToUpper(command)
			user.MID = mid
			userInstances = append(userInstances, user)
		}
	}

	pageAccessToken := os.Getenv("PAGE_TOKEN")
	if pageAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing page access token"})
		return
	}

	for _, user := range userInstances {
		processCommand(user.Command, user.PSID, user.MID, pageAccessToken)
	}
}

func computeHMAC(message []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

func processCommand(command, psid, mid, token string) {
	fmt.Printf("Processing Command: %s, PSID: %s, MID: %s\n", command, psid, mid)
	fmt.Printf("token: %s\n", token)
}
