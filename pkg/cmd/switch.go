package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

const (
	accountName     = "engineName"
	modelAccountKey = "modelName"
	GPTEngine       = "GPT"
	GeminiEngine    = "Gemini"
)

var (
	listModels bool
	modelName  string

	// Available models for each engine
	gptModels = []string{
		"gpt-4",
		"gpt-4-turbo-preview",
		"gpt-3.5-turbo",
		// "o1-preview",
		// "o1-mini",
		"gpt-4o",
		"gpt-4o-2024-11-20",
		"gpt-4o-2024-08-06",
		"gpt-4o-mini",
		"gpt-4o-mini-2024-07-18",
	}

	geminiModels = []string{
		"gemini-1.5-pro",
		"gemini-1.5-flash",
		"gemini-1.5-flash-8b",
	}
)

func init() {
	switchCmd.Flags().BoolVar(&listModels, "list-models", false, "List available models for the current engine")
	switchCmd.Flags().StringVar(&modelName, "model", "", "Switch to specified model")
	rootCmd.AddCommand(switchCmd)
}

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch between different engines (Gemini, GPT) and models",
	Long:  `This command allows you to switch between different engines (Gemini, GPT) and their respective models.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get current engine
		engineName, err := keyring.Get(serviceName, accountName)
		if err != nil {
			color.Red("No engine set. Please run `genie init` to set the engine. (default: Gemini)")
			color.Cyan("Automatically switching to Gemini engine.")
			err := keyring.Set(serviceName, accountName, "Gemini")
			if err != nil {
				fmt.Println("Something went wrong. Please try again.")
				return
			}
			// Set default model for Gemini
			err = keyring.Set(serviceName, modelAccountKey, "gemini-1.5-pro")
			if err != nil {
				fmt.Println("Failed to set default model.")
				return
			}
			color.Green("✓ Switched to Gemini engine with default model (gemini-1.5-pro).")
			return
		}

		// Handle --list-models flag
		if listModels {
			printAvailableModels(engineName)
			return
		}

		// Handle model switching
		if modelName != "" {
			if err := switchModel(engineName, modelName); err != nil {
				color.Red("✗ %s", err.Error())
				return
			}
			return
		}

		// Switch engine
		fmt.Println(color.HiMagentaString(" Switching Engine"))
		fmt.Println(strings.Repeat("─", 50))

		color.Cyan("Current Engine:")
		color.Green("  • %s", engineName)

		newEngine := GeminiEngine
		if engineName == GeminiEngine {
			newEngine = GPTEngine
		}

		if err := keyring.Set(serviceName, accountName, newEngine); err != nil {
			color.Red("✗ Something went wrong. Please try again.")
			return
		}

		// Set default model for the new engine
		defaultModel := "gemini-1.5-pro"
		if newEngine == GPTEngine {
			defaultModel = "gpt-4"
		}
		if err := keyring.Set(serviceName, modelAccountKey, defaultModel); err != nil {
			color.Red("✗ Failed to set default model.")
			return
		}

		fmt.Println(strings.Repeat("─", 50))
		color.Cyan("New Engine:")
		color.Green("  • %s", newEngine)
		color.Green("  • Default model: %s", defaultModel)

		fmt.Println(strings.Repeat("─", 50))
		color.HiBlue("Helpful Commands:")
		fmt.Println("• List models:   genie switch --list-models")
		fmt.Println("• Change model:  genie switch --model <model-name>")
	},
}

func printAvailableModels(engine string) {
	currentModel, _ := keyring.Get(serviceName, modelAccountKey)

	fmt.Println(color.HiMagentaString("📋 Available Models"))
	fmt.Println(strings.Repeat("─", 50))

	color.Cyan("Current Engine:")
	color.Green("  • %s", engine)

	color.Cyan("\nCurrent Model:")
	color.Green("  • %s", currentModel)

	color.Cyan("\nAvailable Models:")
	if engine == GPTEngine {
		for _, model := range gptModels {
			if model == currentModel {
				color.Green("  • %s (current)", model)
			} else {
				fmt.Printf("  • %s\n", model)
			}
		}
	} else {
		for _, model := range geminiModels {
			if model == currentModel {
				color.Green("  • %s (current)", model)
			} else {
				fmt.Printf("  • %s\n", model)
			}
		}
	}

	fmt.Println(strings.Repeat("─", 50))
	color.HiBlue("Helpful Commands:")
	fmt.Println("• Switch engine: genie switch")
	fmt.Println("• Change model: genie switch --model <model-name>")
}

func switchModel(engine, model string) error {
	fmt.Println(color.HiMagentaString("🔄 Switching Model"))
	fmt.Println(strings.Repeat("─", 50))

	color.Cyan("Current Engine:")
	color.Green("  • %s", engine)

	var validModels []string
	if engine == GPTEngine {
		validModels = gptModels
	} else {
		validModels = geminiModels
	}

	// Validate model name
	isValid := false
	for _, validModel := range validModels {
		if strings.EqualFold(validModel, model) {
			isValid = true
			model = validModel // Use the correct case from the valid models list
			break
		}
	}

	if !isValid {
		return fmt.Errorf("Invalid model name for %s engine. Use --list-models to see available models", engine)
	}

	currentModel, _ := keyring.Get(serviceName, modelAccountKey)
	color.Cyan("\nCurrent Model:")
	color.Green("  • %s", currentModel)

	if err := keyring.Set(serviceName, modelAccountKey, model); err != nil {
		return fmt.Errorf("Failed to switch model: %v", err)
	}

	color.Cyan("\nNew Model:")
	color.Green("  • %s", model)

	fmt.Println(strings.Repeat("─", 50))
	return nil
}
