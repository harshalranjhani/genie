package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/config"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(engineCmd)
}

var engineCmd = &cobra.Command{
	Use:   "engine",
	Short: "Get the current engine and model configuration",
	Long: `This command displays detailed information about the current engine and model configuration.
It shows which engine is active (GPT or Gemini) and which specific model is being used.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get current engine
		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			color.Red("No engine configured. Please run `genie init` to set up your configuration.")
			return
		}

		// Get current model
		modelName, err := keyring.Get(serviceName, modelAccountKey)
		if err != nil {
			modelName = "default" // Fallback if no model is explicitly set
		}

		// Print configuration details
		fmt.Println(color.HiMagentaString("🧞 Current Configuration"))
		fmt.Println(strings.Repeat("─", 50))

		// Print engine info
		color.Cyan("Engine:")
		color.Green("  • %s", engineName)

		// Print model info
		color.Cyan("Current Model:")
		color.Green("  • %s", modelName)

		// Print available models for current engine
		color.Cyan("\nAvailable Models:")
		if engineName == config.OllamaEngine {
			models, err := getRunningOllamaModels()
			if err != nil {
				color.Red("Error: %v", err)
				return
			}
			if len(models) == 0 {
				color.Yellow("  No running Ollama models found. Please start some models first.")
				return
			}
			for _, model := range models {
				if model == modelName {
					color.Green("  • %s (current)", model)
				} else {
					fmt.Printf("  • %s\n", model)
				}
			}
		} else {
			engine, exists := config.CheckAndGetEngine(engineName)
			if !exists {
				color.Yellow("  No models available for engine: %s", engineName)
				return
			}

			for _, model := range engine.Models {
				if model == modelName {
					color.Green("  • %s (current)", model)
				} else {
					fmt.Printf("  • %s\n", model)
				}
			}
		}

		// Print helpful commands
		fmt.Println(strings.Repeat("─", 50))
		color.HiBlue("Helpful Commands:")
		fmt.Println("• Switch engine: genie switch")
		fmt.Println("• List models:   genie switch --list-models")
		fmt.Println("• Change model:  genie switch --model <model-name>")
	},
}
