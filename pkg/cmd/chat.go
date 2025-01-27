package cmd

import (
	"log"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/config"
	"github.com/harshalranjhani/genie/internal/helpers/llm"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.PersistentFlags().Bool("safe", false, "Set this to true if you wish to enable safe mode.")
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an interactive chat session",
	Long:  `Start an interactive chat session with the AI model.`,
	Run: func(cmd *cobra.Command, args []string) {
		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		engine, exists := config.CheckAndGetEngine(engineName)
		if !exists {
			log.Fatal("Unknown engine name: ", engineName)
		}

		if !engine.Features.SupportsChat {
			color.Red("%s engine does not support chat yet. Check back soon!", engineName)
			return
		}

		safeSettings, _ := cmd.Flags().GetBool("safe")

		if safeSettings && engine.Features.SupportsSafeMode {
			color.Green("Safety settings are on.")
		}

		switch engineName {
		case config.GPTEngine:
			if safeSettings {
				color.Yellow("Note: Safety settings in GPT are managed through OpenAI's content moderation.")
			}
			llm.StartGPTChat()
		case config.GeminiEngine:
			llm.StartGeminiChat(safeSettings)
		case config.DeepSeekEngine:
			llm.StartDeepSeekChat()
		case config.OllamaEngine:
			model, err := keyring.Get(serviceName, modelAccountKey)
			if err != nil {
				log.Fatal("Error retrieving model name from keyring: ", err)
			}
			llm.StartOllamaChat(model)
		}
	},
}
