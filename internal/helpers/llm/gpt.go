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
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/constants"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/internal/middleware"
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

	model := "gpt-4"
	if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
		model = selectedModel
	}

	req := openai.ChatCompletionRequest{
		Model: model,
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

	model := "gpt-4"
	if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
		model = selectedModel
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
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

	model := "gpt-4"
	if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
		model = selectedModel
	}

	req := openai.ChatCompletionRequest{
		Model: model,
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

	model := "gpt-4"
	if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
		model = selectedModel
	}

	req := openai.ChatCompletionRequest{
		Model: model,
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

func StartGPTChat() {
	ctx := context.Background()
	openAIKey, err := keyring.Get("genie", "openai_api_key")
	if err != nil {
		fmt.Println("OpenAI API key not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}
	client := openai.NewClient(openAIKey)

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#9D4EDD"))

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06D6A0")).
		Bold(true)

	aiStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#118AB2"))

	// Update the readline prompt with promptStyle color
	rl, err := readline.New(promptStyle.Render("You 💭 > "))
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	color.New(color.FgHiMagenta).Println("🧞 Chat session started!")
	fmt.Println(style.Render("Type your message and press Enter to send. Type 'exit' to end the session."))
	fmt.Println(style.Render("Type 'clear' to clear chat history."))
	fmt.Println(style.Render("Type '/history' to export chat history to markdown."))
	fmt.Println(style.Render("Type '/email' to email chat history."))
	fmt.Println(strings.Repeat("─", 50))

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant.",
		},
	}

	for {
		userInput, err := rl.Readline()
		if err != nil {
			break
		}

		userInput = strings.TrimSpace(userInput)

		switch strings.ToLower(userInput) {
		case constants.ExitCommand:
			fmt.Println(style.Render("\n👋 Ending chat session. Goodbye!"))
			return
		case constants.ClearCommand:
			// Clear message history
			messages = []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant.",
				},
			}
			// Clear terminal screen
			fmt.Print("\033[H\033[2J")
			// Reprint welcome message
			color.New(color.FgHiMagenta).Println("🧞 Chat session started!")
			fmt.Println(style.Render("Type your message and press Enter to send. Type 'exit' to end the session."))
			fmt.Println(style.Render("Type 'clear' to clear chat history."))
			fmt.Println(style.Render("Type '/history' to export chat history to markdown."))
			fmt.Println(style.Render("Type '/email' to email chat history."))
			fmt.Println(strings.Repeat("─", 50))
			continue
		case constants.HistoryCommand:
			exportChatHistory(messages)
			continue
		case constants.EmailCommand:
			emailGPTChatHistory(messages)
			continue
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})

		s := spinner.New(spinner.CharSets[11], 80*time.Millisecond)
		s.Prefix = color.HiCyanString("🤔 Thinking: ")
		s.Suffix = " Please wait..."
		s.Start()

		model := "gpt-4"
		if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
			model = selectedModel
		}

		stream, err := client.CreateChatCompletionStream(
			ctx,
			openai.ChatCompletionRequest{
				Model:    model,
				Messages: messages,
				Stream:   true,
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		defer stream.Close()

		s.Stop()
		fmt.Print(color.HiCyanString("\n🤖 AI: "))

		var fullResponse strings.Builder
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				break
			}

			content := response.Choices[0].Delta.Content
			fullResponse.WriteString(content)
			fmt.Print(aiStyle.Render(content))
		}

		// Add the complete response to message history
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: fullResponse.String(),
		})

		fmt.Println("\n" + strings.Repeat("─", 50))
	}
}

func exportChatHistory(messages []openai.ChatCompletionMessage) {
	if len(messages) <= 1 { // Check if there's only the system message or no messages
		fmt.Printf("%s No chat history available to export.\n", color.RedString("❌"))
		return
	}

	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("📝 Exporting chat history: ")
	s.Start()

	timestamp := time.Now().Format("2006-01-02-15-04-05")
	filename := filepath.Join(".", fmt.Sprintf("chat-history-%s.md", timestamp))

	var content strings.Builder
	content.WriteString("# Chat History\n\n")
	content.WriteString(fmt.Sprintf("Generated on: %s\n\n", time.Now().Format("January 2, 2006 15:04:05")))
	content.WriteString("---\n\n")

	for _, msg := range messages[1:] { // Skip the system message
		switch msg.Role {
		case openai.ChatMessageRoleUser:
			content.WriteString(fmt.Sprintf("### 💭 You\n%s\n\n", msg.Content))
		case openai.ChatMessageRoleAssistant:
			content.WriteString(fmt.Sprintf("### 🤖 AI\n%s\n\n", msg.Content))
		}
		content.WriteString("---\n\n")
	}

	err := os.WriteFile(filename, []byte(content.String()), 0644)
	s.Stop()

	if err != nil {
		fmt.Printf("%s Failed to export chat history: %v\n", color.RedString("❌"), err)
		return
	}

	successMsg := fmt.Sprintf("✨ Chat history exported to: %s", filename)
	fmt.Println(color.GreenString(successMsg))
}

func emailGPTChatHistory(messages []openai.ChatCompletionMessage) {
	if len(messages) <= 1 {
		fmt.Printf("%s No chat history available to email.\n", color.RedString("❌"))
		return
	}

	// Create a divider for visual separation
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println(color.HiMagentaString("📧 Emailing Chat History"))
	fmt.Println(strings.Repeat("─", 50))

	// Get current model name
	modelName := "gpt-4"
	if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
		modelName = selectedModel
	}

	// Get user status to check for verified email
	status, err := middleware.LoadStatus()
	var email string
	if err != nil || status == nil || status.Email == "" {
		fmt.Print(color.YellowString("Please enter your email address: "))
		fmt.Scanln(&email)
	} else {
		email = status.Email
	}

	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("📝 Sending to ") + color.CyanString(email) + color.HiCyanString(": ")
	s.Start()

	var chatMessages []map[string]string
	for _, msg := range messages[1:] {
		chatMessages = append(chatMessages, map[string]string{
			"role":    string(msg.Role),
			"content": msg.Content,
		})
	}

	payload := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"model":     modelName,
		"messages":  chatMessages,
		"metadata": map[string]string{
			"sessionId": fmt.Sprintf("gpt-%d", time.Now().Unix()),
			"format":    "markdown",
		},
	}

	if err := helpers.SendChatHistoryEmail(email, payload); err != nil {
		s.Stop()
		fmt.Printf("\n%s Failed to send chat history: %v\n", color.RedString("❌"), err)
		fmt.Println(strings.Repeat("─", 50))
		return
	}

	s.Stop()
	fmt.Printf("\n%s Chat history sent successfully to %s!\n",
		color.GreenString("✨"),
		color.CyanString(email))
	fmt.Println(strings.Repeat("─", 50))
}
