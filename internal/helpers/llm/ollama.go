package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
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
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/constants"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/internal/middleware"
	"github.com/harshalranjhani/genie/pkg/prompts"
	"github.com/zalando/go-keyring"
)

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaRequest struct {
	Model    string                 `json:"model"`
	Messages []OllamaMessage        `json:"messages"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options"`
}

type OllamaResponse struct {
	Model   string        `json:"model"`
	Message OllamaMessage `json:"message"`
	Done    bool          `json:"done"`
}

func getOllamaURL() string {
	url, err := keyring.Get("genie", "ollama_url")
	if err != nil || url == "" {
		return "http://localhost:11434"
	}
	return url
}

func GetOllamaGeneralResponse(prompt string, model string, includeDir bool) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing: ")
	s.Start()
	// Prepare the request
	messages := []OllamaMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestBody := OllamaRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
		Options: map[string]interface{}{
			"temperature": 1.0,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		s.Stop()
		fmt.Printf("Failed to marshal request: %v\n", err)
		return err
	}

	resp, err := http.Post(getOllamaURL()+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.Stop()
		fmt.Printf("Failed to connect to Ollama: %v\n", err)
		fmt.Println("Make sure Ollama is running locally (http://localhost:11434)")
		return err
	}
	defer resp.Body.Close()

	s.Stop()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var streamResponse OllamaResponse
		if err := json.Unmarshal(scanner.Bytes(), &streamResponse); err != nil {
			fmt.Printf("Failed to parse response: %v\n", err)
			continue
		}

		if streamResponse.Message.Content != "" {
			fmt.Printf(formatMarkdownToPlainText(streamResponse.Message.Content))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading stream: %v\n", err)
	}

	return nil
}

func GetOllamaCmdResponse(prompt string, model string, safeOn bool) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing: ")
	s.Start()

	// Prepare the request
	messages := []OllamaMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestBody := OllamaRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
		Options: map[string]interface{}{
			"temperature": 1.0,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		s.Stop()
		fmt.Printf("Failed to marshal request: %v\n", err)
		return err
	}

	resp, err := http.Post(getOllamaURL()+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.Stop()
		fmt.Printf("Failed to connect to Ollama: %v\n", err)
		fmt.Println("Make sure Ollama is running locally (http://localhost:11434)")
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Stop()
		return err
	}

	var response OllamaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		s.Stop()
		return err
	}

	command := response.Message.Content

	s.Stop()
	fmt.Println("Running the command: ", command)
	helpers.RunCommand(command)

	return nil
}

func DocumentCodeWithOllama(filePath string, model string) error {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing code: ")
	s.Start()
	defer s.Stop()

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	prompt := prompts.GetDocumentPrompt(string(content))

	// Prepare the request
	messages := []OllamaMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant who documents code.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestBody := OllamaRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
		Options: map[string]interface{}{
			"temperature": 0.7,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(getOllamaURL()+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response OllamaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	documentedContent := response.Message.Content

	// Extract code from markdown code blocks if present
	re := regexp.MustCompile("(?s)```.*?\n(.*?)\n```")
	matches := re.FindStringSubmatch(documentedContent)
	if len(matches) > 1 {
		documentedContent = matches[1]
	}

	// Write the documented content back to the file
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

func GenerateReadmeWithOllama(readmePath string, templateName string, model string) error {
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

	// Prepare the request
	messages := []OllamaMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant who generates README files.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestBody := OllamaRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
		Options: map[string]interface{}{
			"temperature": 1.0,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(getOllamaURL()+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var response OllamaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	generatedText := response.Message.Content

	if err := helpers.ProcessTemplateResponse(templateName, generatedText, readmePath); err != nil {
		return fmt.Errorf("failed to process template response: %w", err)
	}

	return nil
}

func GenerateBugReportOllama(description, severity, category, assignee, priority, model string) (string, error) {
	// Prepare the request
	messages := []OllamaMessage{
		{
			Role:    "system",
			Content: "You are a helpful software engineer who writes clear, detailed bug reports.",
		},
		{
			Role:    "user",
			Content: prompts.GetBugReportPrompt(description, severity, category, assignee, priority),
		},
	}

	requestBody := OllamaRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
		Options: map[string]interface{}{
			"temperature": 0.7,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(getOllamaURL()+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var response OllamaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Message.Content == "" {
		return "", fmt.Errorf("no response generated")
	}

	return response.Message.Content, nil
}

func StartOllamaChat(model string) {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#9D4EDD"))

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06D6A0")).
		Bold(true)

	aiStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#118AB2"))

	multilineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)

	// Configure readline with multiline support
	config := &readline.Config{
		Prompt:                 promptStyle.Render("You 💭 > "),
		HistoryFile:            "/tmp/genie_ollama_history",
		HistoryLimit:           100,
		DisableAutoSaveHistory: false,
		InterruptPrompt:        "^C",
		EOFPrompt:              "exit",
		EnableMask:             false,
	}
	rl, err := readline.NewEx(config)
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	color.New(color.FgHiMagenta).Println("🧞 Chat session started!")
	fmt.Println(style.Render("Commands: 'exit' | 'clear' | '/history' | '/email'"))
	fmt.Println(multilineStyle.Render("Tip: For multiline input, type '\\' at end of line or use '---' on a new line to send."))
	fmt.Println(strings.Repeat("─", 50))

	messages := []OllamaMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant.",
		},
	}

	for {
		// Read input with multiline support
		var inputLines []string
		isMultiline := false

		for {
			line, err := rl.Readline()
			if err != nil {
				if err == readline.ErrInterrupt {
					if len(inputLines) > 0 {
						// Cancel current multiline input
						inputLines = nil
						isMultiline = false
						rl.SetPrompt(promptStyle.Render("You 💭 > "))
						fmt.Println(style.Render("Input cancelled."))
						break
					}
					continue
				}
				// EOF or other error - exit chat
				fmt.Println(style.Render("\n👋 Ending chat session. Goodbye!"))
				return
			}

			// Check for line continuation (backslash at end)
			if strings.HasSuffix(line, "\\") {
				inputLines = append(inputLines, strings.TrimSuffix(line, "\\"))
				isMultiline = true
				rl.SetPrompt(promptStyle.Render("  ... > "))
				continue
			}

			// Check for multiline end marker
			if isMultiline && strings.TrimSpace(line) == "---" {
				rl.SetPrompt(promptStyle.Render("You 💭 > "))
				break
			}

			inputLines = append(inputLines, line)

			// If not in multiline mode, break after first line
			if !isMultiline {
				break
			}
		}

		if len(inputLines) == 0 {
			continue
		}

		userInput := strings.TrimSpace(strings.Join(inputLines, "\n"))
		if userInput == "" {
			continue
		}

		switch strings.ToLower(userInput) {
		case constants.ExitCommand:
			fmt.Println(style.Render("\n👋 Ending chat session. Goodbye!"))
			return
		case constants.ClearCommand:
			messages = []OllamaMessage{
				{
					Role:    "system",
					Content: "You are a helpful assistant.",
				},
			}
			fmt.Print("\033[H\033[2J")
			color.New(color.FgHiMagenta).Println("🧞 Chat session started!")
			fmt.Println(style.Render("Commands: 'exit' | 'clear' | '/history' | '/email'"))
			fmt.Println(multilineStyle.Render("Tip: For multiline input, type '\\' at end of line or use '---' on a new line to send."))
			fmt.Println(strings.Repeat("─", 50))
			continue
		case constants.HistoryCommand:
			exportOllamaChatHistory(messages)
			continue
		case constants.EmailCommand:
			emailOllamaChatHistory(messages, model)
			continue
		}

		messages = append(messages, OllamaMessage{
			Role:    "user",
			Content: userInput,
		})

		s := spinner.New(spinner.CharSets[11], 80*time.Millisecond)
		s.Prefix = color.HiCyanString("🤔 Thinking: ")
		s.Suffix = " Please wait..."
		s.Start()

		requestBody := OllamaRequest{
			Model:    model,
			Messages: messages,
			Stream:   true,
			Options: map[string]interface{}{
				"temperature": 0.7,
			},
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Post(getOllamaURL()+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		s.Stop()
		fmt.Print(color.HiCyanString("\n🤖 AI: "))

		var fullResponse strings.Builder
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			var streamResponse OllamaResponse
			if err := json.Unmarshal(scanner.Bytes(), &streamResponse); err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				continue
			}

			content := streamResponse.Message.Content
			if content != "" {
				fullResponse.WriteString(content)
				fmt.Print(aiStyle.Render(content))
			}
		}

		messages = append(messages, OllamaMessage{
			Role:    "assistant",
			Content: fullResponse.String(),
		})

		fmt.Println("\n" + strings.Repeat("─", 50))
	}
}

func exportOllamaChatHistory(messages []OllamaMessage) {
	if len(messages) <= 1 {
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

	for _, msg := range messages[1:] {
		switch msg.Role {
		case "user":
			content.WriteString(fmt.Sprintf("### 💭 You\n%s\n\n", msg.Content))
		case "assistant":
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

func emailOllamaChatHistory(messages []OllamaMessage, model string) {
	if len(messages) <= 1 {
		fmt.Printf("%s No chat history available to email.\n", color.RedString("❌"))
		return
	}

	fmt.Println(strings.Repeat("─", 50))
	fmt.Println(color.HiMagentaString("📧 Emailing Chat History"))
	fmt.Println(strings.Repeat("─", 50))

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
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	payload := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"model":     model,
		"messages":  chatMessages,
		"metadata": map[string]string{
			"sessionId": fmt.Sprintf("ollama-%d", time.Now().Unix()),
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
