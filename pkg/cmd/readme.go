package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/harshalranjhani/genie/internal/helpers/llm"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var templateName string
var readmeFileName string

func init() {
	readmeCmd.Flags().StringVarP(&templateName, "template", "t", "default", "Specify the template to use for generating the README (default, minimal, detailed)")
	readmeCmd.Flags().StringVarP(&readmeFileName, "filename", "f", "README.md", "Specify the name of the README file to generate")
	rootCmd.AddCommand(readmeCmd)
}

var readmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Generate README.md for the current directory",
	Long:  `This command generates a README file for your project. You can select from a list of templates or use the default template, and specify the output filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}

		readmePath := filepath.Join(cwd, readmeFileName)

		switch engineName {
		case GPTEngine:
			err := llm.GenerateReadmeWithGPT(readmePath, templateName)
			if err != nil {
				log.Fatalf("Failed to generate README with GPT: %v", err)
			}
		case GeminiEngine:
			err := llm.GenerateReadmeWithGemini(readmePath, templateName)
			if err != nil {
				log.Fatalf("Failed to generate README with Gemini: %v", err)
			}
		default:
			log.Fatal("Unknown engine name: ", engineName)
		}

		fmt.Printf("%s generated successfully!\n", readmeFileName)
	},
}
