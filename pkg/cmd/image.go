package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/helpers/llm"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "image",
	Short: "generate an image from a prompt",
	Long:  `This command generates an image from a prompt using a Python script.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a prompt")
			return
		}

		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		prompt := args[0]
		switch engineName {
		case GPTEngine:
			filePath, err := llm.GenerateGPTImage(prompt)
			if err != nil {
				return
			}
			fmt.Println("Image generated:", filePath)
		case GeminiEngine:
			color.Red("Gemini engine is currently not supported for image generation. Check back soon!")
			os.Exit(1)
		case DeepSeekEngine:
			color.Red("DeepSeek engine is currently not supported for image generation. Check back soon!")
			os.Exit(1)
		default:
			log.Fatal("Unknown engine name: ", engineName)
		}

	},
}
