package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var keys = map[int]string{
	1: openAIKeyName,
	2: geminiKeyName,
	3: deepseekKeyName,
	4: replicateKeyName,
	5: ignoreListPathKeyName,
	6: ollamaURLKeyName,
	7: "all",
	8: "purge",
}

func resetKey(keyName string) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = color.YellowString(" Removing %s...", keyName)
	s.Start()
	time.Sleep(500 * time.Millisecond)

	err := keyring.Delete(serviceName, keyName)
	s.Stop()

	if err != nil {
		color.Red("✘ Failed to delete %s: %s\n", keyName, err)
	} else {
		color.Green("✔ Successfully deleted %s\n", keyName)
	}
}

func resetKeys() {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = color.YellowString(" Removing all keys...")
	s.Start()
	time.Sleep(500 * time.Millisecond)
	s.Stop()

	keys := []string{openAIKeyName, geminiKeyName, deepseekKeyName, ignoreListPathKeyName, replicateKeyName, ollamaURLKeyName}
	for _, key := range keys {
		if err := keyring.Delete(serviceName, key); err != nil {
			color.Red("✘ Failed to delete %s: %s\n", key, err)
		} else {
			color.Green("✔ Successfully deleted %s\n", key)
		}
	}
}

func purgeGenieService() {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = color.RedString(" Purging all Genie data...")
	s.Start()
	time.Sleep(500 * time.Millisecond)
	s.Stop()

	// Delete all known keys first
	keys := []string{openAIKeyName, geminiKeyName, deepseekKeyName, ignoreListPathKeyName, replicateKeyName, ollamaURLKeyName}
	for _, key := range keys {
		_ = keyring.Delete(serviceName, key)
	}

	// Delete the service itself if supported by the system keyring
	if err := keyring.Delete(serviceName, serviceName); err != nil {
		color.Yellow("Note: Service entry could not be removed (this is normal on some systems)")
	}

	color.Red("🗑️  All Genie data has been purged from your system")
}

func init() {
	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset your API keys",
	Long: color.YellowString(`
🔄 Genie Reset Wizard
--------------------
This wizard will help you reset and update your API keys stored in the system keychain.
You can choose to reset individual keys or all of them at once.
`),
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Println(color.CyanString("\n🔑 API Key Reset Menu"))
			fmt.Println(color.HiBlackString("──────────────────────"))
			options := map[string]string{
				"1": "Reset OpenAI API Key 🤖",
				"2": "Reset Gemini API Key 🧞",
				"3": "Reset DeepSeek API Key 🔄",
				"4": "Reset Replicate API Key 🔄",
				"5": "Reset Ignore List Path 📝",
				"6": "Reset Ollama URL 🌐",
				"7": "Reset All Keys ⚠️",
				"8": "Purge All Genie Data 🗑️",
				"0": "Exit 👋",
			}

			for num, text := range options {
				if num == "0" {
					fmt.Println(color.HiBlackString("──────────────────────"))
				}
				fmt.Printf("%s: %s\n", color.HiYellowString(num), text)
			}

			fmt.Print(color.CyanString("\n➜ Your choice: "))
			scanner.Scan()
			userResponse := scanner.Text()
			choice, err := strconv.Atoi(userResponse)

			if err != nil || choice < 0 || choice > len(keys) {
				color.Red("❌ Invalid choice! Please try again.")
				continue
			}
			if choice == 0 {
				color.Cyan("👋 Goodbye!")
				break
			}

			keyName := keys[choice]
			if keyName == "purge" {
				fmt.Print(color.RedString("⚠️  WARNING: This will completely remove all Genie data from your system.\nAre you absolutely sure? (yes/N): "))
				scanner.Scan()
				confirm := scanner.Text()
				if confirm == "yes" {
					purgeGenieService()
					color.Cyan("👋 Goodbye!")
					return
				}
				continue
			}
			fmt.Println(keyName)
			if keyName == "all" {
				fmt.Print(color.YellowString("⚠️  Are you sure you want to reset all keys? (y/N): "))
				scanner.Scan()
				confirm := scanner.Text()
				if confirm == "y" || confirm == "Y" {
					resetKeys()
					color.Yellow("\n🔄 Would you like to set up new keys now? (y/N): ")
					scanner.Scan()
					confirm = scanner.Text()
					if confirm == "y" || confirm == "Y" {
						fmt.Println(color.CyanString("\n📦 Starting Genie re-initialization...\n"))
						initCmd.Run(cmd, args)
					}
					return
				}
				continue
			}

			resetKey(keyName)

			color.Yellow("\n🔄 Would you like to set a new value for this key? (y/N): ")
			scanner.Scan()
			confirm := scanner.Text()
			if confirm == "y" || confirm == "Y" {
				var prompt string
				switch keyName {
				case openAIKeyName:
					prompt = "Enter your OpenAI API Key:"
				case geminiKeyName:
					prompt = "Enter your Gemini API Key:"
				case ignoreListPathKeyName:
					prompt = "Enter the path to the ignore list file:"
				case replicateKeyName:
					prompt = "Enter your Replicate API Key:"
				case deepseekKeyName:
					prompt = "Enter your DeepSeek API Key:"
				case ollamaURLKeyName:
					prompt = "Enter your Ollama URL:"
				default:
					continue
				}

				storeKeyIfNotPresent(keyName, prompt, "")
			}
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = color.GreenString(" Finalizing changes...")
		s.Start()
		time.Sleep(1 * time.Second)
		s.Stop()

		fmt.Println(color.GreenString("\n✨ All changes have been saved successfully!\n"))
	},
}
