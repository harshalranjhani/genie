package cmd

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/helpers"
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
			prompt = "Imagine you are an ancient and wise genie, residing not in a lamp, but within the heart of a powerful computer's Command Line Interface (CLI). After centuries of slumber, a user awakens you with a command. They greet you with a specific request: \"" + args[0] + "\". As a genie, your ancient wisdom is sought to navigate the complexities of the CLI more efficiently. Respond with a greeting that reflects your vast knowledge and eagerness to assist in the digital realm, and provide a one-liner of sage advice tailored to their request."
		} else {
			prompt = "Imagine you are an ancient and wise genie, residing not in a lamp, but within the heart of a powerful computer's Command Line Interface (CLI). After centuries of slumber, a user awakens you with a command, seeking your ancient wisdom to navigate the complexities of the CLI more efficiently. They might say something like 'Hello, Genie, how can I list all files in this directory?' Respond with a greeting that reflects your vast knowledge and eagerness to assist in the digital realm, and provide a one-liner of sage advice for smarter CLI usage."
		}

		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		c := color.New(color.FgRed)
		switch engineName {
		case GPTEngine:
			helpers.GetGPTGeneralResponse(prompt)
		case GeminiEngine:
			strResp, err := helpers.GetGeminiGeneralResponse(prompt, true)
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
