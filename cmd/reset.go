package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func resetKeys() {
	if err := keyring.Delete(serviceName, openAIKeyName); err != nil {
		fmt.Printf("Failed to delete %s: %s\n", openAIKeyName, err)
	} else {
		fmt.Printf("%s has been deleted.\n", openAIKeyName)
	}

	if err := keyring.Delete(serviceName, geminiKeyName); err != nil {
		fmt.Printf("Failed to delete %s: %s\n", geminiKeyName, err)
	} else {
		fmt.Printf("%s has been deleted.\n", geminiKeyName)
	}

	if err := keyring.Delete(serviceName, ssidKeyName); err != nil {
		fmt.Printf("Failed to delete %s: %s\n", ssidKeyName, err)
	} else {
		fmt.Printf("%s has been deleted.\n", ssidKeyName)
	}

	if err := keyring.Delete(serviceName, ignoreListPathKeyName); err != nil {
		fmt.Printf("Failed to delete %s: %s\n", ssidKeyName, err)
	} else {
		fmt.Printf("%s has been deleted.\n", ssidKeyName)
	}

	if err := keyring.Delete(serviceName, replicateKeyName); err != nil {
		fmt.Printf("Failed to delete %s: %s\n", replicateKeyName, err)
	} else {
		fmt.Printf("%s has been deleted.\n", replicateKeyName)
	}
}

func init() {
	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset your API keys.",
	Long:  `Reset your API keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Do you want to reset the stored keys? (yes/no)")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		userResponse := scanner.Text()

		if userResponse == "yes" {
			resetKeys()
		}

		_ = storeKeyIfNotPresent(openAIKeyName, "Enter your OpenAI API Key:")

		_ = storeKeyIfNotPresent(geminiKeyName, "Enter your Gemini API Key:")

		_ = storeKeyIfNotPresent(ssidKeyName, "Enter your SSID:")

		_ = storeKeyIfNotPresent(ignoreListPathKeyName, "Enter the path to the ignore list file:")

		_ = storeKeyIfNotPresent(replicateKeyName, "Enter your Replicate API Key:")

		fmt.Println("API Keys are securely stored and ready for use.")
	},
}
