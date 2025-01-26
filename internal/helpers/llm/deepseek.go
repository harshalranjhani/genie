package llm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/pkg/prompts"
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

func GetDeepSeekGeneralResponse(prompt string, safeOn bool, includeDir bool) error {
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

	// Create streaming chat completion request
	request := &deepseek.StreamChatCompletionRequest{
		Model: modelName,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    constants.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}

	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to create stream from DeepSeek: %v", err)
	}
	defer stream.Close()

	s.Stop()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println()
			return nil
		}

		if err != nil {
			return fmt.Errorf("stream error: %v", err)
		}

		for _, choice := range response.Choices {
			fmt.Print(choice.Delta.Content)
		}
	}
}

func DocumentCodeWithDeepSeek(filePath string) error {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing code: ")
	s.Start()
	defer s.Stop()

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	deepseekKey, err := keyring.Get("genie", "deepseek_api_key")
	if err != nil {
		return fmt.Errorf("DeepSeek API key not found in keyring: please run `genie init` to store the key: %w", err)
	}

	client := deepseek.NewClient(deepseekKey)
	ctx := context.Background()

	prompt := prompts.GetDocumentPrompt(string(content))

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

	req := &deepseek.ChatCompletionRequest{
		Model: modelName,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    constants.ChatMessageRoleSystem,
				Content: "You are a helpful assistant who documents code.",
			},
			{
				Role:    constants.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 1.0,
		ResponseFormat: &deepseek.ResponseFormat{
			Type: "text",
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return errors.New("no response from DeepSeek API")
	}

	documentedContent := resp.Choices[0].Message.Content

	re := regexp.MustCompile("(?s)```.*?\n(.*?)\n```")
	matches := re.FindStringSubmatch(documentedContent)
	if len(matches) > 1 {
		documentedContent = matches[1]
	}

	err = ioutil.WriteFile(filePath, []byte(""), 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, []byte(strings.TrimSpace(documentedContent)), 0644)
	if err != nil {
		return err
	}

	return nil
}
