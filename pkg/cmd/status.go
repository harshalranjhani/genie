package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/middleware"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display genie status information",
	Long:  `Show detailed information about genie's current configuration, engine status, and system information.`,
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Prefix = color.HiCyanString("Fetching status: ")
		s.Start()

		// Get version from the Version constant
		version := Version

		// Get current engine
		engineName, _ := keyring.Get(serviceName, "engineName")
		if engineName == "" {
			engineName = "Gemini (default)"
		}

		// Check API keys status
		openAIKey, _ := keyring.Get(serviceName, openAIKeyName)
		geminiKey, _ := keyring.Get(serviceName, geminiKeyName)
		replicateKey, _ := keyring.Get(serviceName, replicateKeyName)
		ssidKey, _ := keyring.Get(serviceName, ssidKeyName)
		ignoreListPath, _ := keyring.Get(serviceName, ignoreListPathKeyName)

		// Add verification status check
		status, _ := middleware.LoadStatus()

		time.Sleep(500 * time.Millisecond)
		s.Stop()

		// Clear screen
		fmt.Print("\033[H\033[2J")

		// Header
		cyan := color.New(color.FgCyan, color.Bold)
		cyan.Println("\n🧞 Genie Status Dashboard")
		fmt.Println(strings.Repeat("─", 50))

		// Version and System Info
		fmt.Printf("📌 %s: %s\n", color.HiBlackString("Version"), color.HiGreenString(version))
		fmt.Printf("💻 %s: %s\n", color.HiBlackString("System"), color.HiGreenString(runtime.GOOS))
		fmt.Printf("⚙️  %s: %s\n", color.HiBlackString("Engine"), color.HiGreenString(engineName))

		// Add Verification Status section before Configuration Status
		fmt.Println("\n🔐 Verification Status")
		fmt.Println(strings.Repeat("─", 25))
		if status != nil && status.Email != "" {
			fmt.Printf("📧 %s: ", color.HiBlackString("Status"))
			color.Green("✓ Verified")
			fmt.Printf("   %s: %s\n", color.HiBlackString("Email"), color.HiBlackString(status.Email))
		} else {
			fmt.Printf("📧 %s: ", color.HiBlackString("Status"))
			color.Yellow("! Not verified")
		}

		// Configuration Status
		fmt.Println("\n🔧 Configuration Status")
		fmt.Println(strings.Repeat("─", 25))

		printKeyStatus("OpenAI API", openAIKey)
		printKeyStatus("Gemini API", geminiKey)
		printKeyStatus("Replicate API", replicateKey)
		printKeyStatus("SSID", ssidKey)

		// Ignore List Status
		fmt.Printf("📝 %s: ", color.HiBlackString("Ignore List"))
		if ignoreListPath != "" {
			if _, err := os.Stat(ignoreListPath); err == nil {
				color.Green("✓ Configured")
				fmt.Printf("   %s: %s\n", color.HiBlackString("Path"), color.HiBlackString(ignoreListPath))
			} else {
				color.Red("✗ File not found")
			}
		} else {
			color.Yellow("! Not configured")
		}

		// Footer
		fmt.Println("\n💡 Tips")
		fmt.Println(strings.Repeat("─", 10))
		fmt.Printf("• Use %s to update configuration\n", color.CyanString("genie init"))
		fmt.Printf("• Use %s to switch engines\n", color.CyanString("genie switch"))
		fmt.Printf("• Use %s to reset configuration\n", color.CyanString("genie reset"))
		fmt.Println()
	},
}

func printKeyStatus(name string, key string) {
	fmt.Printf("🔑 %s: ", color.HiBlackString(name))
	if key != "" {
		color.Green("✓ Configured")
		fmt.Printf("   %s: %s\n", color.HiBlackString("Key"), color.HiBlackString(maskKey(key)))
	} else {
		color.Yellow("! Not configured")
	}
}
