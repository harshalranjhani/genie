package cmd

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/config"
	"github.com/harshalranjhani/genie/internal/helpers/llm"
	"github.com/harshalranjhani/genie/pkg/prompts"
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

		_, exists := config.CheckAndGetEngine(engineName)
		if !exists {
			log.Fatal("Unknown engine name: ", engineName)
		}

		c := color.New(color.FgRed)
		switch engineName {
		case config.GPTEngine:
			llm.GetGPTGeneralResponse(prompt, false)
		case config.GeminiEngine:
			strResp, err := llm.GetGeminiGeneralResponse(prompt, true, false)
			if err != nil {
				log.Fatal("Error getting response from Gemini: ", err)
				os.Exit(1)
			}
			c.Println(formatMarkdownToPlainText(strResp))
		case config.DeepSeekEngine:
			err := llm.GetDeepSeekGeneralResponse(prompt, true, false)
			if err != nil {
				log.Fatal("Error getting response from DeepSeek: ", err)
				os.Exit(1)
			}
		case config.OllamaEngine:
			model, err := keyring.Get(serviceName, modelAccountKey)
			if err != nil {
				log.Fatal("Error retrieving model name from keyring: ", err)
			}
			llm.GetOllamaGeneralResponse(prompt, model, false)
		}
	},
}
