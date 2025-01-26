package llm

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"
)

func GetDeepSeekCmdResponse(prompt string, safeOn bool) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing: ")
	s.Start()

	godotenv.Load()
	ctx := context.Background()

	deepseekKey, err := keyring.Get("genie", "deepseek_api_key")
	if err != nil {
		s.Stop()
		fmt.Println("DeepSeek API key not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}
	client := deepseek.NewClient(deepseekKey)

	// Get the selected model from keyring
	modelName := deepseek.DeepSeekChat
	if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
		switch selectedModel {
		case "deepseek-chat":
			modelName = deepseek.DeepSeekChat
		case "deepseek-reasoner":
			modelName = deepseek.DeepSeekReasoner
		default:
			modelName = deepseek.DeepSeekChat
		}
	}

	// Create chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: modelName,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    constants.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.7,
		ResponseFormat: &deepseek.ResponseFormat{
			Type: "text",
		},
	}

	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to get response from DeepSeek: %v", err)
	}

	s.Stop()

	// Print the response
	if len(response.Choices) > 0 {
		fmt.Println("Running the command: ", response.Choices[0].Message.Content)
		helpers.RunCommand(response.Choices[0].Message.Content)
	} else {
		return fmt.Errorf("no response received from DeepSeek")
	}

	return nil
}
