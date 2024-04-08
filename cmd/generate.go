package cmd

import (
	"fmt"

	"github.com/harshalranjhani/genie/helpers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate an image from a prompt",
	Long:  `This command generates an image from a prompt using a Python script.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a prompt")
			return
		}

		prompt := args[0]

		filePath, err := helpers.GenerateImage(prompt)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Image generated:", filePath)
	},
}
