package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/harshalranjhani/genie/internal/structs"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func GetStatusFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".genie")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.Mkdir(configDir, 0700)
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(configDir, "user_status.json"), nil
}

func LoadStatus() (*structs.UserStatus, error) {
	statusFile, err := GetStatusFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(statusFile); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(statusFile)
	if err != nil {
		return nil, err
	}

	var status structs.UserStatus
	err = json.Unmarshal(data, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func SendVerificationEmail(email string) error {
	client := &http.Client{Timeout: time.Second * 10}
	reqBody, _ := json.Marshal(map[string]string{
		"email": email,
	})
	resp, err := client.Post("https://api.harshalranjhani.in/genie/send-verification", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send verification email")
	}
	return nil
}

func SaveStatus(status *structs.UserStatus) error {
	statusFile, err := GetStatusFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(status)
	if err != nil {
		return err
	}

	tempFile := statusFile + ".temp"
	err = os.WriteFile(tempFile, data, 0444) // write with read-only permissions
	if err != nil {
		return err
	}

	// Rename the temporary file to the actual status file
	err = os.Rename(tempFile, statusFile)
	if err != nil {
		return err
	}

	return nil
}

func WaitForVerification(email string) (string, error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://api.harshalranjhani.in/", nil)
	if err != nil {
		return "", fmt.Errorf("failed to connect to WebSocket server: %v", err)
	}
	defer c.Close()

	message := map[string]string{"email": email}
	msg, _ := json.Marshal(message)
	c.WriteMessage(websocket.TextMessage, msg)

	_, messageBytes, err := c.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("error reading WebSocket message: %v", err)
	}

	var response map[string]string
	if err := json.Unmarshal(messageBytes, &response); err != nil {
		return "", fmt.Errorf("error unmarshaling WebSocket message: %v", err)
	}

	token, ok := response["token"]
	if !ok {
		return "", fmt.Errorf("Payment not found.")
	}

	return token, nil
}

func TokenValid() (bool, error) {
	status, err := LoadStatus()
	if err != nil {
		return false, err
	}
	if status == nil {
		return false, fmt.Errorf("no status found")
	}

	if time.Now().Unix() > status.Expiry {
		return false, fmt.Errorf("token expired")
	}

	return true, nil
}

func VerifySubscriptionMiddleware(cmd *cobra.Command, args []string) error {
	valid, err := TokenValid()
	if !valid {
		message := fmt.Sprintf(color.RedString("Subscription verification required: %v\n", err) +
			color.CyanString("Please run the following command to re-verify your email:\n") +
			color.YellowString("\tgenie verify [email]\n"))
		fmt.Println(message)
		return fmt.Errorf("subscription verification failed")
	}
	// fmt.Println(color.GreenString("Subscription verified successfully."))
	return nil
}

func CheckKeyringSetup() error {
	// Check if the service exists by attempting to read a key
	_, err := keyring.Get("genie", "ignore_list_path")
	if err != nil {
		return fmt.Errorf("genie service not initialized. Please run 'genie init' to set up required configurations")
	}

	// Check for required fields
	ignorePath, err := keyring.Get("genie", "ignore_list_path")
	if err != nil || ignorePath == "" {
		return fmt.Errorf("ignore list path not configured. Please run 'genie init' to set up required configurations")
	}

	return nil
}
