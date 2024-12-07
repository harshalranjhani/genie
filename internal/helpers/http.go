package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendChatHistoryEmail(email string, chatHistory interface{}) error {
	payload := map[string]interface{}{
		"email":       email,
		"chatHistory": chatHistory,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	resp, err := http.Post(
		"https://api.harshalranjhani.in/genie/send-chat-history",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
