package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/config"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/internal/helpers/llm"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var filePathToConnect string

func init() {
	documentCmd.Flags().StringVarP(&filePathToConnect, "file", "f", "", "Path to the file to be documented")
	documentCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(documentCmd)
}

var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Document your code with genie",
	Long:  `Transform your code with genie comments with great documentation which can be later used easily to get summaries of your code.`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath := filePathToConnect
		if !filepath.IsAbs(filePathToConnect) {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("Failed to get current working directory: %v", err)
			}
			filePath = filepath.Join(cwd, filePathToConnect)
		}

		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		engine, exists := config.GetEngine(engineName)
		if !exists {
			log.Fatal("Unknown engine name: ", engineName)
		}

		if !engine.Features.SupportsDocumentation {
			color.Yellow("%s engine does not support documentation generation yet. Check back soon!", engineName)
			return
		}

		switch engineName {
		case config.GPTEngine:
			err := llm.DocumentCodeWithGPT(filePath)
			if err != nil {
				log.Fatalf("Failed to document code: %v", err)
			}
			color.Green("Code documented successfully!")
			pathToOpen := fmt.Sprintf("code %s", filePath)
			helpers.RunCommand(pathToOpen)
		case config.DeepSeekEngine:
			err := llm.DocumentCodeWithDeepSeek(filePath)
			if err != nil {
				log.Fatalf("Failed to document code: %v", err)
			}
			color.Green("Code documented successfully!")
		}
	},
}
