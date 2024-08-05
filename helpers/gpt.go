package helpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
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

func DocumentCodeWithGPT(filePath string) error {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing code: ")
	s.Start()
	defer s.Stop()
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	openAIKey, err := keyring.Get("genie", "openai_api_key")
	if err != nil {
		return err
	}
	client := openai.NewClient(openAIKey)
	ctx := context.Background()

	prompt := fmt.Sprintf(`Document the following code with Genie comments. 
Genie comments provide clear, structured headings and subheadings for the code to enhance readability and provide detailed documentation. 
Use Genie comments to explain the purpose, functionality, and usage of each part of the code. The format for genie comments is as follows:
In python:

# genie:heading: This is a heading
# genie:subheading: This is a subheading

or in javascript:

// genie:heading: This is a heading
// genie:subheading: This is a subheading

Make sure to match the exact format for the comments to be detected correctly. The format is genie:heading: for headings and genie:subheading: for subheadings. Remember to add a space after the colon and before the text. Also add a space after the comment marker (# or //) and before the genie keyword. Remember there cannot be multiple genie headings in one file, but there can be multiple genie subheadings under one heading.
Here is the code:
%s\nRemember to output the whole code including all imports, exports, functions, tests, etc. You are supposed to add genie comments wherever necessary and then return the whole code. Give the output as code only, no other text is required.`, content)

	req := openai.ChatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant who documents code.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return err
	}

	if len(resp.Choices) == 0 {
		return errors.New("no response from OpenAI API")
	}

	documentedContent := resp.Choices[0].Message.Content

	re := regexp.MustCompile("(?s)```.*?\n(.*?)\n```")
	matches := re.FindStringSubmatch(documentedContent)
	if len(matches) > 1 {
		documentedContent = matches[1]
	}

	err = ioutil.WriteFile(filePath, []byte(""), 0644)

	err = ioutil.WriteFile(filePath, []byte(strings.TrimSpace(documentedContent)), 0644)
	if err != nil {
		return err
	}

	return nil
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
