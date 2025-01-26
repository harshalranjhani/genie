package cmd

import (
	"bufio"
	"fmt"
	"os"

	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

const serviceName = "genie"
const openAIKeyName = "openai_api_key"
const geminiKeyName = "gemini_api_key"
const ignoreListPathKeyName = "ignore_list_path"
const replicateKeyName = "replicate_api_key"
const deepseekKeyName = "deepseek_api_key"

func getAPIKeyFromUser(promptMessage string) string {
	fmt.Print(color.HiBlackString("Enter your key: "))
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func storeKeyIfNotPresent(accountName string, promptMessage string, emoji string) string {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = color.HiBlackString(" Checking existing keys...")
	s.Start()

	apiKey, err := keyring.Get(serviceName, accountName)
	time.Sleep(1 * time.Second)
	s.Stop()

	if err != nil {
		fmt.Println(color.HiBlackString("\n────────────────────────────────────"))
		fmt.Printf("%s %s\n", emoji, color.CyanString(promptMessage))
		fmt.Println(color.HiBlackString("────────────────────────────────────"))

		apiKey = getAPIKeyFromUser(promptMessage)

		s.Suffix = color.HiBlackString(" Storing key securely...")
		s.Start()
		err := keyring.Set(serviceName, accountName, apiKey)
		time.Sleep(500 * time.Millisecond)
		s.Stop()

		if err != nil {
			color.Red("❌ Failed to store %s: %s\n", accountName, err)
			os.Exit(1)
		}
		color.Green("✅ Successfully stored %s\n", accountName)
	}

	return apiKey
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Store your API keys securely in the system keychain",
	Long: color.YellowString(`
✨ Genie Initialization Wizard
────────────────────────────
This wizard will help you set up Genie by storing your API keys securely 
in your system keychain. Each key will be encrypted and can only be 
accessed by Genie.

Let's get started! 🚀
`),
	Run: func(cmd *cobra.Command, args []string) {
		// Check for existing keys
		_, errOpenAI := keyring.Get(serviceName, openAIKeyName)
		_, errGemini := keyring.Get(serviceName, geminiKeyName)
		_, errIgnoreListPath := keyring.Get(serviceName, ignoreListPathKeyName)
		_, errDeepseek := keyring.Get(serviceName, deepseekKeyName)
		if errOpenAI == nil || errGemini == nil || errIgnoreListPath == nil || errDeepseek == nil {
			color.Yellow("\n⚠️  Some or all keys are already present!")
			fmt.Print(color.HiBlackString("Use "))
			color.Cyan("genie reset")
			fmt.Println(color.HiBlackString(" to update them."))
			return
		}

		// Welcome message
		fmt.Println(color.CyanString("\n🎮 Starting Genie initialization...\n"))
		fmt.Println(color.HiBlackString("We'll help you set up everything step by step!\n"))

		// Collect all keys
		openAIKey := storeKeyIfNotPresent(openAIKeyName, "Enter your OpenAI API Key", "🤖")
		geminiKey := storeKeyIfNotPresent(geminiKeyName, "Enter your Gemini API Key", "🧞")
		deepseekKey := storeKeyIfNotPresent(deepseekKeyName, "Enter your DeepSeek API Key", "🔄")
		replicateKey := storeKeyIfNotPresent(replicateKeyName, "Enter your Replicate API Key", "🔄")
		ignoreListPath := storeKeyIfNotPresent(ignoreListPathKeyName, "Enter the path to your ignore list file", "📝")
		// Set default engine
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = color.HiBlackString(" Setting default engine...")
		s.Start()
		err := keyring.Set(serviceName, "engineName", "Gemini")
		time.Sleep(500 * time.Millisecond)
		s.Stop()

		if err != nil {
			color.Red("❌ Failed to store engineName: %s\n", err)
			os.Exit(1)
		}

		// Final success message
		fmt.Println(color.GreenString("\n🎉 Success! Genie is now configured and ready to use!\n"))

		// Configuration summary
		fmt.Println(color.YellowString("📋 Configuration Summary"))
		fmt.Println(color.HiBlackString("────────────────────────"))
		fmt.Printf("%s OpenAI API Key: %s\n", color.HiBlackString("├─ 🤖"), maskKey(openAIKey))
		fmt.Printf("%s Gemini API Key: %s\n", color.HiBlackString("├─ 🧞"), maskKey(geminiKey))
		fmt.Printf("%s DeepSeek API Key: %s\n", color.HiBlackString("└─ 🔄"), maskKey(deepseekKey))
		fmt.Printf("%s Replicate API Key: %s\n", color.HiBlackString("├─ 🔄"), maskKey(replicateKey))
		fmt.Printf("%s Ignore List Path: %s\n", color.HiBlackString("├─ 📝"), ignoreListPath)

		// Next steps
		fmt.Println(color.CyanString("\n📚 Next Steps"))
		fmt.Println(color.HiBlackString("────────────"))
		fmt.Printf("• Run %s to see all available commands\n", color.CyanString("genie help"))
		fmt.Printf("• Run %s to update your keys\n", color.CyanString("genie reset"))
		fmt.Println()
	},
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return color.HiBlackString("********")
	}
	return color.HiBlackString(key[:4] + "..." + key[len(key)-4:])
}
