package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(switchCmd)
}

const accountName = "engineName"

const GPTEngine = "GPT"
const GeminiEngine = "Gemini"

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch between different engines (Gemini, GPT)",
	Long:  `This command allows you to switch between different engines (Gemini, GPT).`,
	Run: func(cmd *cobra.Command, args []string) {
		engineName, err := keyring.Get(serviceName, accountName)

		if err != nil {
			color.Red("No engine set. Please run `genie init` to set the engine. (default: Gemini)")
			color.Cyan("Automatically switching to Gemini engine.")
			err := keyring.Set(serviceName, accountName, "Gemini")
			if err != nil {
				fmt.Println("Something went wrong. Please try again.")
				return
			}
			color.Green("Switched to Gemini engine.")
			return
		}

		color.Cyan("Current engine: %s", engineName)

		if engineName == "Gemini" {
			err := keyring.Set(serviceName, accountName, "GPT")
			if err != nil {
				fmt.Println("Something went wrong. Please try again.")
				return
			}
			color.Green("Switched to GPT engine.")
		} else {
			err := keyring.Set(serviceName, accountName, "Gemini")
			if err != nil {
				fmt.Println("Something went wrong. Please try again.")
				return
			}
			color.Green("Switched to Gemini engine.")
		}
	},
}
