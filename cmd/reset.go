package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var keys = map[int]string{
	1: openAIKeyName,
	2: geminiKeyName,
	3: ssidKeyName,
	4: ignoreListPathKeyName,
	5: replicateKeyName,
	6: "all",
}

func resetKey(keyName string) {
	if err := keyring.Delete(serviceName, keyName); err != nil {
		fmt.Printf("Failed to delete %s: %s\n", keyName, err)
	} else {
		fmt.Printf("%s has been deleted.\n", keyName)
	}
}
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
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Println("\n--- API Key Reset Menu ---")
			fmt.Println("Please select an option:")
			fmt.Println("1: Reset OpenAI API Key")
			fmt.Println("2: Reset Gemini API Key")
			fmt.Println("3: Reset SSID")
			fmt.Println("4: Reset Ignore List Path")
			fmt.Println("5: Reset Replicate API Key")
			fmt.Println("6: Reset All Keys")
			fmt.Println("0: Exit")
			fmt.Print("Your choice: ")
			scanner.Scan()
			userResponse := scanner.Text()
			choice, err := strconv.Atoi(userResponse)

			if err != nil || choice < 0 || choice > len(keys) {
				fmt.Println("Invalid choice")
				continue
			}
			if choice == 0 {
				fmt.Println("Exiting...")
				break
			}

			keyName := keys[choice]
			if keyName == "all" {
				resetKeys()
				return
			}
			resetKey(keyName)

			var prompt string
			switch keyName {
			case openAIKeyName:
				prompt = "Enter your OpenAI API Key:"
			case geminiKeyName:
				prompt = "Enter your Gemini API Key:"
			case ssidKeyName:
				prompt = "Enter the SSID key:"
			case ignoreListPathKeyName:
				prompt = "Enter the path to the ignore list file:"
			case replicateKeyName:
				prompt = "Enter your Replicate API Key:"
			default:
				continue
			}

			storeKeyIfNotPresent(keyName, prompt)

		}
		fmt.Println("API Keys are securely stored and ready for use.")

	},
}
