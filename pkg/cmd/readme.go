package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/config"
	"github.com/harshalranjhani/genie/internal/helpers/llm"
	"github.com/harshalranjhani/genie/internal/middleware"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var templateName string
var readmeFileName string

func init() {
	readmeCmd.Flags().StringVarP(&templateName, "template", "t", "default", "Specify the template to use for generating the README (default, minimal, detailed, animated, interactive)")
	readmeCmd.Flags().StringVarP(&readmeFileName, "filename", "f", "README.md", "Specify the name of the README file to generate")
	rootCmd.AddCommand(readmeCmd)
}

var readmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Generate README.md for the current directory",
	Long:  `This command generates a README file for your project. You can select from a list of templates or use the default template, and specify the output filename.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if templateName == "animated" || templateName == "interactive" {
			if err := middleware.VerifySubscriptionMiddleware(cmd, args); err != nil {
				log.Fatal(err)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		_, exists := config.GetEngine(engineName)
		if !exists {
			log.Fatal("Unknown engine name: ", engineName)
		}

		cwd, err := os.Getwd()
		if err != nil {
			color.Red("Failed to get current working directory: %v", err)
			os.Exit(1)
		}

		readmePath := filepath.Join(cwd, readmeFileName)
		switch engineName {
		case config.GPTEngine:
			err := llm.GenerateReadmeWithGPT(readmePath, templateName)
			if err != nil {
				log.Fatalf("Failed to generate README with GPT: %v", err)
			}
		case config.GeminiEngine:
			err := llm.GenerateReadmeWithGemini(readmePath, templateName)
			if err != nil {
				log.Fatalf("Failed to generate README with Gemini: %v", err)
			}
		case config.DeepSeekEngine:
			err := llm.GenerateReadmeWithDeepSeek(readmePath, templateName)
			if err != nil {
				log.Fatalf("Failed to generate README with DeepSeek: %v", err)
			}
		case config.OllamaEngine:
			model, err := keyring.Get(serviceName, modelAccountKey)
			if err != nil {
				log.Fatal("Error retrieving model name from keyring: ", err)
			}
			err = llm.GenerateReadmeWithOllama(readmePath, templateName, model)
			if err != nil {
				log.Fatalf("Failed to generate README with Ollama: %v", err)
			}
		}

		fmt.Printf("%s generated successfully!\n", readmeFileName)
	},
}
