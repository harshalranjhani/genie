package cmd

import (
	"log"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/helpers/llm"
	"github.com/harshalranjhani/genie/internal/middleware"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.PersistentFlags().Bool("safe", false, "Set this to true if you wish to enable safe mode.")
}

var chatCmd = &cobra.Command{
	Use:     "chat",
	Short:   "Start a chat with the genie and maintain a conversation.",
	Long:    `Use this command to chat with the genie and maintain a conversation directly from the terminal.`,
	PreRunE: middleware.VerifySubscriptionMiddleware,
	Run: func(cmd *cobra.Command, args []string) {
		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}
		safeSettings, _ := cmd.Flags().GetBool("safe")

		switch engineName {
		case GPTEngine:
			if safeSettings {
				color.Yellow("Note: Safety settings in GPT are managed through OpenAI's content moderation.")
			}
			llm.StartGPTChat()
		case GeminiEngine:
			llm.StartGeminiChat(safeSettings)
		case DeepSeekEngine:
			llm.StartDeepSeekChat()
		default:
			log.Fatal("Unknown engine name: ", engineName)
		}
	},
}
