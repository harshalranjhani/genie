package helpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
	"github.com/zalando/go-keyring"
)

func GetGPTGeneralResponse(prompt string, includeDir bool) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing: ")
	s.Start()

	openAIKey, err := keyring.Get("genie", "openai_api_key")
	if err != nil {
		s.Stop()
		fmt.Println("OpenAI API key not found in keyring. Please run `genie init` to store the key.")
		return
	}
	c := openai.NewClient(openAIKey)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		s.Stop()
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	s.Stop()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println()
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf(formatMarkdownToPlainText(response.Choices[0].Delta.Content))
	}
}

func GetGPTCmdResponse(prompt string, safeOn bool) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing: ")
	s.Start()

	openAIKey, err := keyring.Get("genie", "openai_api_key")
	if err != nil {
		s.Stop()
		fmt.Println("OpenAI API key not found in keyring. Please run `genie init` to store the key.")
		return err
	}
	client := openai.NewClient(openAIKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		s.Stop()
		return err
	}

	command := resp.Choices[0].Message.Content

	if safeOn {
		isSafe, err := checkModeration(openAIKey, command)
		if err != nil {
			s.Stop()
			return err
		}
		if !isSafe {
			s.Stop()
			fmt.Println("The generated command contains inappropriate content.")
			return errors.New("inappropriate content detected")
		}
	}

	s.Stop()
	fmt.Println("Running the command: ", command)
	RunCommand(command)

	return nil
}

func checkModeration(apiKey, content string) (bool, error) {
	client := openai.NewClient(apiKey)
	resp, err := client.Moderations(context.Background(), openai.ModerationRequest{
		Input: content,
	})

	if err != nil {
		return false, err
	}

	for _, result := range resp.Results {
		if result.Flagged {
			return false, nil
		}
	}

	return true, nil
}

func formatMarkdownToPlainText(mdText string) string {
	// Regular expressions to replace Markdown formatting
	reStrong := regexp.MustCompile(`\*\*(.*?)\*\*`)
	reEmphasis := regexp.MustCompile(`\*(.*?)\*`)
	reCode := regexp.MustCompile("([^])" + "`" + "(.*?)" + "`" + "([^`])")
	reHeaders := regexp.MustCompile(`\n#+\s(.*?)\n`)

	// Replace Markdown syntax with plain text formatting
	plainText := reStrong.ReplaceAllString(mdText, "$1")
	plainText = reEmphasis.ReplaceAllString(plainText, "$1")
	plainText = reCode.ReplaceAllString(plainText, "$1")
	plainText = reHeaders.ReplaceAllString(plainText, "\n$1\n")

	return plainText
}
