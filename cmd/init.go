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
const ignoreListPathKeyName = "ignore_list_path"
const replicateKeyName = "replicate_api_key"

func getAPIKeyFromUser(promptMessage string) string {
	fmt.Println(promptMessage)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // Wait for input
	return scanner.Text()
}

func storeKeyIfNotPresent(accountName string, promptMessage string) string {
	apiKey, err := keyring.Get(serviceName, accountName)

	if err != nil {
		apiKey = getAPIKeyFromUser(promptMessage)

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

		_, errOpenAI := keyring.Get(serviceName, openAIKeyName)
		_, errGemini := keyring.Get(serviceName, geminiKeyName)
		_, errSSID := keyring.Get(serviceName, ssidKeyName)
		_, errIgnoreListPath := keyring.Get(serviceName, ignoreListPathKeyName)
		if errOpenAI == nil || errGemini == nil || errSSID == nil || errIgnoreListPath == nil {
			fmt.Println("Some or all keys are already present. Please reset using `genie reset` if you want to update them.")
			return
		}

		openAIKey := storeKeyIfNotPresent(openAIKeyName, "Enter your OpenAI API Key:")

		geminiKey := storeKeyIfNotPresent(geminiKeyName, "Enter your Gemini API Key:")

		ssidKey := storeKeyIfNotPresent(ssidKeyName, "Enter your SSID:")

		ignoreListPath := storeKeyIfNotPresent(ignoreListPathKeyName, "Enter the path to your ignore list file:")

		replicateKey := storeKeyIfNotPresent(replicateKeyName, "Enter your Replicate API Key:")

		fmt.Println("API Keys are securely stored and ready for use.")

		fmt.Println("OpenAI API Key:", openAIKey)
		fmt.Println("Gemini API Key:", geminiKey)
		fmt.Println("SSID:", ssidKey)
		fmt.Println("Ignore List Path:", ignoreListPath)
		fmt.Println("Replicate API Key:", replicateKey)
	},
}
