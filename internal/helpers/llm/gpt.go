package llm

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/pkg/prompts"
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
	helpers.RunCommand(command)

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

	prompt := prompts.GetDocumentPrompt(string(content))

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

func GenerateGPTImage(prompt string) (string, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Generating Image: ")
	s.Start()

	openAIKey, err := keyring.Get("genie", "openai_api_key")
	if err != nil {
		s.Stop()
		fmt.Println("OpenAI API key not found in keyring. Please run `genie init` to store the key.")
		return "", err
	}
	c := openai.NewClient(openAIKey)
	ctx := context.Background()

	// reqUrl := openai.ImageRequest{
	// 	Prompt:         prompt,
	// 	Size:           openai.CreateImageSize256x256,
	// 	ResponseFormat: openai.CreateImageResponseFormatURL,
	// 	N:              1,
	// }
	// respUrl, err := c.CreateImage(ctx, reqUrl)
	// if err != nil {
	// 	s.Stop()
	// 	fmt.Printf("Image creation error: %v\n", err)
	// 	return "", err
	// }
	// fmt.Println(respUrl.Data[0].URL)

	reqBase64 := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	respBase64, err := c.CreateImage(ctx, reqBase64)
	if err != nil {
		s.Stop()
		fmt.Printf("Image creation error: %v\n", err)
		return "", err
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		s.Stop()
		fmt.Printf("Base64 decode error: %v\n", err)
		return "", err
	}

	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		s.Stop()
		fmt.Printf("PNG decode error: %v\n", err)
		return "", err
	}

	prompt = strings.Replace(prompt, " ", "_", -1)
	filename := fmt.Sprintf("%s.png", prompt)

	file, err := os.Create(filename)
	if err != nil {
		s.Stop()
		fmt.Printf("File creation error: %v\n", err)
		return "", err
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		s.Stop()
		fmt.Printf("PNG encode error: %v\n", err)
		return "", err
	}

	s.Stop()

	return filename, nil
}

func GenerateReadmeWithGPT(readmePath string, templateName string) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	rootDir, err := helpers.GetCurrentDirectoriesAndFiles(cwd)
	if err != nil {
		return fmt.Errorf("failed to get directory structure: %w", err)
	}
	s.Prefix = color.HiCyanString("Generating README: ")
	s.Start()
	defer s.Stop()

	var repoData strings.Builder
	helpers.PrintData(&repoData, rootDir, 0)

	sanitizedRepoData := helpers.SanitizeUTF8(repoData.String())

	// get project name from root folder name
	projectName := filepath.Base(cwd)

	prompt := prompts.GetReadmePrompt(sanitizedRepoData, templateName, projectName)

	openAIKey, err := keyring.Get("genie", "openai_api_key")
	if err != nil {
		return fmt.Errorf("OpenAI API key not found in keyring: please run `genie init` to store the key: %w", err)
	}

	client := openai.NewClient(openAIKey)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant who generates README files.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return errors.New("no response from OpenAI API")
	}

	generatedText := resp.Choices[0].Message.Content

	if err := helpers.ProcessTemplateResponse(templateName, generatedText, readmePath); err != nil {
		return fmt.Errorf("failed to process template response: %w", err)
	}

	return nil
}

func GenerateBugReportGPT(description, severity, category, assignee, priority string) (string, error) {
	openAIKey, err := keyring.Get("genie", "openai_api_key")
	if err != nil {
		return "", fmt.Errorf("OpenAI API key not found: %w", err)
	}

	client := openai.NewClient(openAIKey)
	ctx := context.Background()

	prompt := prompts.GetBugReportPrompt(description, severity, category, assignee, priority)

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: "You are a helpful software engineer who writes clear, detailed bug reports.",
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error generating bug report: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return resp.Choices[0].Message.Content, nil
}
