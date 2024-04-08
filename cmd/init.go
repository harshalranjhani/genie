package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

const serviceName = "genie"
const openAIKeyName = "openai_api_key"
const geminiKeyName = "gemini_api_key"
const ssidKeyName = "ssid"

func getAPIKeyFromUser(promptMessage string) string {
	fmt.Println(promptMessage)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // Wait for input
	return scanner.Text()
}

func storeKeyIfNotPresent(accountName string, promptMessage string) string {
	// Try to get the API key from keyring
	apiKey, err := keyring.Get(serviceName, accountName)

	if err != nil {
		// If API key is not found, prompt the user
		apiKey = getAPIKeyFromUser(promptMessage)

		// Store the API key securely
		err := keyring.Set(serviceName, accountName, apiKey)
		if err != nil {
			fmt.Printf("Failed to store %s: %s\n", accountName, err)
			os.Exit(1)
		}
	}

	return apiKey
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "store your API keys securely in the system keychain.",
	Long:  `store your API keys securely in the system keychain.`,
	Run: func(cmd *cobra.Command, args []string) {

		// check if the API keys are already stored
		_, err := keyring.Get(serviceName, openAIKeyName)
		_, err2 := keyring.Get(serviceName, geminiKeyName)
		if err == nil || err2 == nil {
			fmt.Println("Keys are already present.")
			return
		}

		openAIKey := storeKeyIfNotPresent(openAIKeyName, "Enter your OpenAI API Key:")

		geminiKey := storeKeyIfNotPresent(geminiKeyName, "Enter your Gemini API Key:")

		ssidKey := storeKeyIfNotPresent(ssidKeyName, "Enter your SSID:")

		// Use the API keys for your application's logic
		fmt.Println("API Keys are securely stored and ready for use.")

		fmt.Println("OpenAI API Key:", openAIKey)
		fmt.Println("Gemini API Key:", geminiKey)
		fmt.Println("SSID:", ssidKey)
	},
}
