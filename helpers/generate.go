package helpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/zalando/go-keyring"
)

func GenerateImage(prompt string) (string, error) {
	// get ssid from keyring
	ssid, keyRingError := keyring.Get("genie", "ssid")
	fmt.Println("SSID:", ssid)
	if keyRingError != nil {
		fmt.Println("SSID not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}
	cmd := exec.Command("python", "scripts/generate.py", prompt, ssid)

	// A buffer to capture the output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute python script: %w", err)
	}

	filename := out.String()

	return filename, nil
}
