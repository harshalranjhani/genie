package cmd

import (
	"fmt"
	"log"

	"github.com/harshalranjhani/genie/helpers"
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
			filePath, err := helpers.GenerateGPTImage(prompt)
			if err != nil {
				return
			}
			fmt.Println("Image generated:", filePath)
		case GeminiEngine:
			filePath, err := helpers.GenerateGeminiImage(prompt)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Image generated:", filePath)
		default:
			log.Fatal("Unknown engine name: ", engineName)
		}

	},
}
