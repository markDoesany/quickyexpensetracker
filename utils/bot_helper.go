package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"quickyexpensetracker/templates"
)

func ComputeHMAC(message []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

func SendGenerateRequest(elements []templates.Template, PSID string, pageAccessToken string) error {
	fmt.Println("Sending Generate Request")

	if len(elements) == 0 {
		return fmt.Errorf("elements slice cannot be empty when calling SendGenerateRequest")
	}

	var interfaceSlice []interface{} = make([]interface{}, len(elements))
	for i, d := range elements {
		interfaceSlice[i] = d
	}

	payload := templates.RequestPayload{
		Recipient: templates.Recipient{ID: PSID},
		Message: templates.Message{
			Attachment: templates.Attachment{
				Type: "template",
				Payload: templates.AttachmentPayload{
					TemplateType: "generic",
					Elements:     interfaceSlice,
				},
			},
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	url := fmt.Sprintf("https://graph.facebook.com/v21.0/me/messages?access_token=%s", pageAccessToken)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return nil
}

type Message struct {
	Text string `json:"text"`
}

type TextPayload struct {
	Recipient     templates.Recipient `json:"recipient"`
	Message       Message             `json:"message"`
	MessagingType string              `json:"messaging_type"`
}

func SendTextMessage(message string, PSID string, pageAccessToken string) error {
	client := &http.Client{}

	payload := TextPayload{
		Recipient:     templates.Recipient{ID: PSID},
		Message:       Message{Text: message},
		MessagingType: "RESPONSE",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %v", err)
	}

	url := fmt.Sprintf("https://graph.facebook.com/v21.0/me/messages?access_token=%s", pageAccessToken)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
