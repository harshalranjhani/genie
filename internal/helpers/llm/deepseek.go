package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/chzyer/readline"
	"github.com/cohesion-org/deepseek-go"
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/constants"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/internal/middleware"
	"github.com/harshalranjhani/genie/pkg/prompts"
	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"
)

type DeepSeekStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role             string `json:"role,omitempty"`
			Content          string `json:"content,omitempty"`
			ReasoningContent string `json:"reasoning_content,omitempty"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason,omitempty"`
	} `json:"choices"`
}

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

func StartDeepSeekChat() {
	ctx := context.Background()
	deepseekKey, err := keyring.Get("genie", "deepseek_api_key")
	if err != nil {
		fmt.Println("DeepSeek API key not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#9D4EDD"))

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06D6A0")).
		Bold(true)

	aiStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#118AB2"))

	reasoningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFB703")).
		Italic(true)

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

	messages := []map[string]string{
		{
			"role":    "system",
			"content": "You are a helpful assistant.",
		},
	}

	client := &http.Client{}

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
			messages = []map[string]string{
				{
					"role":    "system",
					"content": "You are a helpful assistant.",
				},
			}
			fmt.Print("\033[H\033[2J")
			color.New(color.FgHiMagenta).Println("🧞 Chat session started!")
			fmt.Println(style.Render("Type your message and press Enter to send. Type 'exit' to end the session."))
			fmt.Println(style.Render("Type 'clear' to clear chat history."))
			fmt.Println(style.Render("Type '/history' to export chat history to markdown."))
			fmt.Println(style.Render("Type '/email' to email chat history."))
			fmt.Println(strings.Repeat("─", 50))
			continue
		case constants.HistoryCommand:
			exportDeepSeekChatHistory(messages)
			continue
		case constants.EmailCommand:
			emailDeepSeekChatHistory(messages)
			continue
		}

		messages = append(messages, map[string]string{
			"role":    "user",
			"content": userInput,
		})

		s := spinner.New(spinner.CharSets[11], 80*time.Millisecond)
		s.Prefix = color.HiCyanString("🤔 Thinking: ")
		s.Suffix = " Please wait..."
		s.Start()

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

		requestBody := map[string]interface{}{
			"model":    modelName,
			"messages": messages,
			"stream":   true,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+deepseekKey)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		s.Stop()
		fmt.Print(color.HiCyanString("\n🤖 AI: "))

		var reasoningContent strings.Builder
		var currentContent strings.Builder

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("\nError reading stream: %v\n", err)
				break
			}

			// Skip empty lines
			if len(line) <= 1 {
				continue
			}

			// Remove "data: " prefix
			data := bytes.TrimPrefix(line, []byte("data: "))
			if bytes.Equal(data, line) {
				continue // Line doesn't start with "data: "
			}

			// Parse the JSON response
			var streamResp DeepSeekStreamResponse
			if err := json.Unmarshal(data, &streamResp); err != nil {
				continue
			}

			for _, choice := range streamResp.Choices {
				if modelName == deepseek.DeepSeekReasoner && choice.Delta.ReasoningContent != "" {
					reasoningContent.WriteString(choice.Delta.ReasoningContent)
				}
				if choice.Delta.Content != "" {
					currentContent.WriteString(choice.Delta.Content)
					fmt.Print(aiStyle.Render(choice.Delta.Content))
				}
			}
		}

		// Print reasoning content only if using reasoner model
		if modelName == deepseek.DeepSeekReasoner && reasoningContent.Len() > 0 {
			fmt.Printf("\n%s\n", reasoningStyle.Render("💡 Reasoning:\n"+reasoningContent.String()))
		}

		// Update messages based on model
		if modelName == deepseek.DeepSeekReasoner {
			messages = append(messages, map[string]string{
				"role":              "assistant",
				"content":           currentContent.String(),
				"reasoning_content": reasoningContent.String(),
			})
		} else {
			messages = append(messages, map[string]string{
				"role":    "assistant",
				"content": currentContent.String(),
			})
		}

		fmt.Println("\n" + strings.Repeat("─", 50))
	}
}

func exportDeepSeekChatHistory(messages []map[string]string) {
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
	content.WriteString("# DeepSeek Chat History\n\n")
	content.WriteString(fmt.Sprintf("Generated on: %s\n\n", time.Now().Format("January 2, 2006 15:04:05")))
	content.WriteString("---\n\n")

	for _, msg := range messages[1:] { // Skip the system message
		switch msg["role"] {
		case "user":
			content.WriteString(fmt.Sprintf("### 💭 You\n%s\n\n", msg["content"]))
		case "assistant":
			aiResponse := msg["content"]
			if reasoning, ok := msg["reasoning_content"]; ok && reasoning != "" {
				aiResponse = fmt.Sprintf("%s\n\n💡 Reasoning:\n%s", aiResponse, reasoning)
			}
			content.WriteString(fmt.Sprintf("### 🤖 AI\n%s\n\n", aiResponse))
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

func emailDeepSeekChatHistory(messages []map[string]string) {
	if len(messages) <= 1 {
		fmt.Printf("%s No chat history available to email.\n", color.RedString("❌"))
		return
	}

	// Create a divider for visual separation
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println(color.HiMagentaString("📧 Emailing Chat History"))
	fmt.Println(strings.Repeat("─", 50))

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
	for _, msg := range messages[1:] { // Skip system message
		if msg["role"] == "assistant" {
			content := msg["content"]
			if reasoning, ok := msg["reasoning_content"]; ok && reasoning != "" {
				content = fmt.Sprintf("%s\n\n💡 Reasoning:\n%s", content, reasoning)
			}
			chatMessages = append(chatMessages, map[string]string{
				"role":    msg["role"],
				"content": content,
			})
		} else {
			chatMessages = append(chatMessages, map[string]string{
				"role":    msg["role"],
				"content": msg["content"],
			})
		}
	}

	payload := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"model":     "deepseek-reasoner",
		"messages":  chatMessages,
		"metadata": map[string]string{
			"sessionId": fmt.Sprintf("deepseek-%d", time.Now().Unix()),
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

func GenerateBugReportDeepSeek(description, severity, category, assignee, priority string) (string, error) {
	deepseekKey, err := keyring.Get("genie", "deepseek_api_key")
	if err != nil {
		return "", fmt.Errorf("DeepSeek API key not found: %w", err)
	}

	client := deepseek.NewClient(deepseekKey)
	ctx := context.Background()

	prompt := prompts.GetBugReportPrompt(description, severity, category, assignee, priority)

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

	request := &deepseek.ChatCompletionRequest{
		Model: modelName,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    constants.ChatMessageRoleSystem,
				Content: "You are a helpful software engineer who writes clear, detailed bug reports.",
			},
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
		return "", fmt.Errorf("error generating bug report: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return response.Choices[0].Message.Content, nil
}
