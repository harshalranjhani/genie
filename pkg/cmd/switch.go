package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/config"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

const (
	accountName     = "engineName"
	modelAccountKey = "modelName"
	ollamaAPIURL    = "http://localhost:11434/api/ps"
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

	// Add DeepSeek models
	deepseekModels = []string{
		"deepseek-chat",
		"deepseek-reasoner",
	}
)

type OllamaModel struct {
	Name    string `json:"name"`
	Model   string `json:"model"`
	Details struct {
		ParameterSize string `json:"parameter_size"`
	} `json:"details"`
}

type OllamaResponse struct {
	Models []OllamaModel `json:"models"`
}

func init() {
	switchCmd.Flags().BoolVar(&listModels, "list-models", false, "List available models for the current engine")
	switchCmd.Flags().StringVar(&modelName, "model", "", "Switch to specified model")
	rootCmd.AddCommand(switchCmd)
}

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch between different engines (Gemini, GPT, DeepSeek) and models",
	Long:  `This command allows you to switch between different engines (Gemini, GPT, DeepSeek) and their respective models.`,
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

		// Modified engine switching logic
		var newEngine string
		currentEngine, exists := config.CheckAndGetEngine(engineName)
		if !exists {
			newEngine = "Gemini"
		} else {
			newEngine = config.GetNextEngine(currentEngine.Name)
		}

		if err := keyring.Set(serviceName, accountName, newEngine); err != nil {
			color.Red("✗ Something went wrong. Please try again.")
			return
		}

		// Set default model for the new engine
		newEngineConfig, _ := config.CheckAndGetEngine(newEngine)
		defaultModel := newEngineConfig.DefaultModel
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
	switch engine {
	case config.GPTEngine:
		for _, model := range gptModels {
			if model == currentModel {
				color.Green("  • %s (current)", model)
			} else {
				fmt.Printf("  • %s\n", model)
			}
		}
	case config.GeminiEngine:
		for _, model := range geminiModels {
			if model == currentModel {
				color.Green("  • %s (current)", model)
			} else {
				fmt.Printf("  • %s\n", model)
			}
		}
	case config.DeepSeekEngine:
		for _, model := range deepseekModels {
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

func getRunningOllamaModels() ([]string, error) {
	resp, err := http.Get(ollamaAPIURL)
	if err != nil {
		return nil, fmt.Errorf("Ollama server not running. Please ensure your Ollama server is running")
	}
	defer resp.Body.Close()

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("Failed to parse Ollama response: %v", err)
	}

	var models []string
	for _, model := range ollamaResp.Models {
		models = append(models, model.Name)
	}
	return models, nil
}

func switchModel(engine, model string) error {
	var validModels []string

	switch engine {
	case config.OllamaEngine:
		models, err := getRunningOllamaModels()
		if err != nil {
			color.Red("\nError: %v", err)
			return nil
		}
		if len(models) == 0 {
			color.Yellow("\nNo running Ollama models found. Please start some models first.")
			return nil
		}
		validModels = models
	case config.GPTEngine:
		validModels = gptModels
	case config.GeminiEngine:
		validModels = geminiModels
	case config.DeepSeekEngine:
		validModels = deepseekModels
	}

	// If listing models is requested
	if listModels {
		color.Cyan("\nAvailable models for %s:", engine)
		for _, m := range validModels {
			fmt.Printf("  • %s\n", m)
		}
		fmt.Println(strings.Repeat("─", 50))
		return nil
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
