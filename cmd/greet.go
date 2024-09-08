package cmd

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/helpers/llm"
	"github.com/harshalranjhani/genie/helpers/prompts"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(greetCmd)
}

var greetCmd = &cobra.Command{
	Use:   "greet [prompt in quotes or nothing]",
	Short: "Invoke the wise Genie for CLI guidance",
	Long: `Invoke the ancient and wise Genie to assist you with your CLI needs. 
The Genie, residing within the heart of your powerful computer's Command Line Interface, 
is ready to provide sage advice and practical tips for smarter CLI usage. 
Whether you're navigating complex commands or seeking general guidance, the Genie is here to help.`,
	Run: func(cmd *cobra.Command, args []string) {
		var prompt string
		if len(args) > 0 {
			prompt = prompts.GetGreetPrompt(args[0])
		} else {
			prompt = prompts.GetGreetPrompt("")
		}

		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		c := color.New(color.FgRed)
		switch engineName {
		case GPTEngine:
			llm.GetGPTGeneralResponse(prompt, false)
		case GeminiEngine:
			strResp, err := llm.GetGeminiGeneralResponse(prompt, true, false)
			if err != nil {
				log.Fatal("Error getting response from Gemini: ", err)
				os.Exit(1)
			}
			c.Println(formatMarkdownToPlainText(strResp))
		default:
			log.Fatal("Unknown engine name: ", engineName)
		}
	},
}
