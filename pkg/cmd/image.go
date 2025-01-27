package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/config"
	"github.com/harshalranjhani/genie/internal/helpers/llm"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(imageCmd)
}

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate images using AI",
	Long:  `Generate images using AI models. Currently supports DALL-E through GPT.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			color.Red("Please provide a prompt for image generation")
			return
		}

		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		engine, exists := config.GetEngine(engineName)
		if !exists {
			log.Fatal("Unknown engine name: ", engineName)
		}

		if !engine.Features.SupportsImageGen {
			color.Red("%s engine does not support image generation yet. Check back soon!", engineName)
			os.Exit(1)
			return
		}

		prompt := args[0]
		filePath, err := llm.GenerateGPTImage(prompt)
		if err != nil {
			return
		}
		fmt.Println("Image generated:", filePath)
	},
}
