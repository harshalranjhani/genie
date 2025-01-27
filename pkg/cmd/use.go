package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/config"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var (
	engineFlag string
	modelFlag  string
)

func init() {
	useCmd.Flags().StringVar(&engineFlag, "engine", "", "Specify the engine (GPT, Gemini, DeepSeek, Ollama)")
	useCmd.Flags().StringVar(&modelFlag, "model", "", "Specify the model to use")
	useCmd.MarkFlagRequired("engine")
	useCmd.MarkFlagRequired("model")
	rootCmd.AddCommand(useCmd)
}

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Directly switch to a specific engine and model combination",
	Long:  `Switch to a specific engine and model combination in one command. Example: genie use --engine GPT --model gpt-4`,
	Run: func(cmd *cobra.Command, args []string) {
		// Normalize engine name to match config
		engineName := strings.Title(strings.ToLower(engineFlag))

		// Validate engine
		engine, exists := config.CheckAndGetEngine(engineName)
		if !exists {
			color.Red("Invalid engine. Available engines: GPT, Gemini, DeepSeek, Ollama")
			return
		}

		// Validate model based on engine
		var validModels []string
		if engineName == config.OllamaEngine {
			models, err := getRunningOllamaModels()
			if err != nil {
				color.Red("Error: %v", err)
				return
			}
			if len(models) == 0 {
				color.Yellow("No running Ollama models found. Please start some models first.")
				return
			}
			validModels = models
		} else {
			validModels = engine.Models
		}

		isValidModel := false
		var correctModelName string
		for _, validModel := range validModels {
			if strings.EqualFold(validModel, modelFlag) {
				isValidModel = true
				correctModelName = validModel
				break
			}
		}

		if !isValidModel {
			color.Red("Invalid model for %s engine. Available models:", engineName)
			for _, model := range validModels {
				fmt.Printf("  • %s\n", model)
			}
			return
		}

		// Set engine and model
		if err := keyring.Set(serviceName, accountName, engineName); err != nil {
			color.Red("Failed to set engine: %v", err)
			return
		}

		if err := keyring.Set(serviceName, modelAccountKey, correctModelName); err != nil {
			color.Red("Failed to set model: %v", err)
			return
		}

		// Print success message
		fmt.Println(color.HiMagentaString("🧞 Configuration Updated"))
		fmt.Println(strings.Repeat("─", 50))
		color.Cyan("Engine:")
		color.Green("  • %s", engine.Name)
		color.Cyan("\nModel:")
		color.Green("  • %s", correctModelName)
	},
}
