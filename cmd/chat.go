package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/google/generative-ai-go/genai"
	"github.com/harshalranjhani/genie/helpers"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"google.golang.org/api/option"
)

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.PersistentFlags().Bool("safe", false, "Set this to true if you wish to enable safe mode.")
}

var chatCmd = &cobra.Command{
	Use:     "chat",
	Short:   "Start a chat with the genie and maintain a conversation.",
	Long:    `Use this command to chat with the genie and maintain a conversation directly from the terminal.`,
	PreRunE: helpers.VerifySubscriptionMiddleware,
	Run: func(cmd *cobra.Command, args []string) {
		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}
		switch engineName {
		case GPTEngine:
			color.Red("Chat is not supported for GPT engine. Please switch to Gemini engine to start a chat.")
		case GeminiEngine:
			safeSettings, _ := cmd.Flags().GetBool("safe")
			startChat(safeSettings)
		default:
			log.Fatal("Unknown engine name: ", engineName)
		}

	},
}

func startChat(safeOn bool) {
	ctx := context.Background()
	geminiKey, err := keyring.Get("genie", "gemini_api_key")
	if err != nil {
		fmt.Println("Gemini API key not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	if safeOn {
		model.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: genai.HarmBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: genai.HarmBlockLowAndAbove,
			},
		}
	} else {
		model.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockNone,
			},
		}
	}
	cs := model.StartChat()

	scanner := bufio.NewScanner(os.Stdin)
	color.Red("Chat session started. Type your message and press Enter to send. Type 'exit' to end the session.")

	for {
		fmt.Print("> ")
		scanner.Scan()
		userInput := scanner.Text()

		if strings.ToLower(userInput) == "exit" {
			fmt.Println("Ending chat session. Goodbye!")
			break
		}

		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(userInput),
			},
			Role: "user",
		})

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Prefix = color.HiCyanString("Analyzing: ")
		s.Start()

		genResp, err := cs.SendMessage(ctx, genai.Text(userInput))
		s.Stop()

		if err != nil {
			log.Fatal(err)
		}

		if len(genResp.Candidates) > 0 && len(genResp.Candidates[0].Content.Parts) > 0 {
			fmt.Println(genResp.Candidates[0].Content.Parts[0])
		}
	}
}
