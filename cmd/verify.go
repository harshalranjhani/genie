package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

type UserStatus struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func getStatusFilePath() (string, error) {
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

func loadStatus() (*UserStatus, error) {
	statusFile, err := getStatusFilePath()
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

	var status UserStatus
	err = json.Unmarshal(data, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func saveStatus(status *UserStatus) error {
	statusFile, err := getStatusFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(status)
	if err != nil {
		return err
	}

	return os.WriteFile(statusFile, data, 0600)
}

func sendVerificationEmail(email string) error {
	client := &http.Client{Timeout: time.Second * 10}
	reqBody, _ := json.Marshal(map[string]string{
		"email": email,
	})
	resp, err := client.Post("https://harshalranjhaniapi.onrender.com/genie/send-verification", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send verification email")
	}
	return nil
}

func waitForVerification(email string) (string, error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://harshalranjhaniapi.onrender.com", nil)
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

var verifyCmd = &cobra.Command{
	Use:   "verify [email]",
	Short: "Verify your support status and get access to extra features.",
	Long:  `If you have donated to the project, you can verify your email to get access to extra features. This command will send a verification email to the provided email address.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]

		status, err := loadStatus()
		if err != nil {
			fmt.Println(color.RedString("Error loading status:"), err)
			return
		}

		if status != nil && status.Email == email {
			fmt.Println(color.GreenString("User is already verified."))
			return
		} else if status != nil && status.Email != email {
			fmt.Println(color.YellowString("A different user is already verified. Removing existing status..."))
			// remove the status file
			statusFile, err := getStatusFilePath()
			if err != nil {
				fmt.Println(color.RedString("Error getting status file path:"), err)
				return
			}
			err = os.Remove(statusFile)
			if err != nil {
				fmt.Println(color.RedString("Error removing status file:"), err)
				return
			}
		}

		fmt.Println(color.CyanString("Sending verification email... Please do not close this tab until the process is complete."))
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Prefix = color.CyanString("Sending: ")
		s.Start()

		err = sendVerificationEmail(email)
		s.Stop()
		if err != nil {
			fmt.Println(color.RedString("Error sending verification email:"), err)
			return
		}

		fmt.Println(color.GreenString("Verification email sent. Please check your inbox."))
		fmt.Println(color.CyanString("Waiting for verification..."))

		s.Prefix = color.CyanString("Verifying: ")
		s.Start()
		token, err := waitForVerification(email)
		s.Stop()
		if err != nil {
			fmt.Println(color.RedString("Error during verification:"), err)
			return
		}

		status = &UserStatus{Email: email, Token: token}
		err = saveStatus(status)
		if err != nil {
			fmt.Println(color.RedString("Error saving status:"), err)
			return
		}

		fmt.Println(color.GreenString("Email verified successfully. You now have access to extra features!"))
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
