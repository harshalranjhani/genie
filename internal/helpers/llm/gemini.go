package llm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bufio"
	"log"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/constants"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/internal/middleware"
	"github.com/harshalranjhani/genie/pkg/prompts"
	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"
	"google.golang.org/genai"
)

func GetGeminiCmdResponse(prompt string, safeOn bool) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing: ")
	s.Start()

	godotenv.Load()
	ctx := context.Background()

	geminiKey, err := keyring.Get("genie", "gemini_api_key")
	if err != nil {
		s.Stop()
		fmt.Println("Gemini API key not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  geminiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		s.Stop()
		return err
	}

	modelName := "gemini-2.5-flash" // default model
	if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
		modelName = selectedModel
	}

	config := &genai.GenerateContentConfig{}
	if safeOn {
		config.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
		}
	} else {
		config.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThresholdBlockNone,
			},
		}
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, genai.Text(prompt), config)
	if err != nil {
		s.Stop()
		return err
	}

	generatedText := resp.Text()
	if generatedText == "" {
		s.Stop()
		fmt.Println("No generated text found")
		return nil
	}

	s.Stop()
	fmt.Println("Running the command: ", generatedText)
	helpers.RunCommand(generatedText)
	return nil
}

func GetGeminiGeneralResponse(prompt string, safeOn bool, includeDir bool) (string, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing: ")
	s.Start()

	godotenv.Load()
	ctx := context.Background()

	geminiKey, err := keyring.Get("genie", "gemini_api_key")
	if err != nil {
		s.Stop()
		return "", fmt.Errorf("gemini API key not found in keyring: please run `genie init` to store the key: %w", err)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  geminiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		s.Stop()
		return "", err
	}

	modelName := "gemini-2.5-flash" // default model
	if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
		modelName = selectedModel
	}

	config := &genai.GenerateContentConfig{}
	if safeOn {
		config.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
		}
	} else {
		config.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThresholdBlockNone,
			},
		}
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, genai.Text(prompt), config)
	if err != nil {
		s.Stop()
		return "", err
	}

	generatedText := resp.Text()
	s.Stop()

	if generatedText == "" {
		return "No generated text found.", nil
	}

	return generatedText, nil
}

func GenerateReadmeWithGemini(readmePath string, templateName string) error {
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

	geminiKey, err := keyring.Get("genie", "gemini_api_key")
	if err != nil {
		return fmt.Errorf("gemini API key not found in keyring: please run `genie init` to store the key: %w", err)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  geminiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}

	modelName := "gemini-2.5-flash" // default model
	if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
		modelName = selectedModel
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, genai.Text(prompt), nil)
	if err != nil {
		return err
	}

	generatedText := resp.Text()
	if generatedText != "" {
		if err := helpers.ProcessTemplateResponse(templateName, generatedText, readmePath); err != nil {
			return fmt.Errorf("failed to process template response: %w", err)
		}
	}

	return nil
}

func GenerateBugReportGemini(description, severity, category, assignee, priority string) (string, error) {
	ctx := context.Background()

	geminiKey, err := keyring.Get("genie", "gemini_api_key")
	if err != nil {
		return "", fmt.Errorf("gemini API key not found: %w", err)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  geminiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create gemini client: %w", err)
	}

	modelName := "gemini-2.5-flash" // default model
	if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
		modelName = selectedModel
	}

	prompt := prompts.GetBugReportPrompt(description, severity, category, assignee, priority)

	resp, err := client.Models.GenerateContent(ctx, modelName, genai.Text(prompt), nil)
	if err != nil {
		return "", err
	}

	generatedText := resp.Text()
	if generatedText == "" {
		return "", fmt.Errorf("no response generated")
	}

	return generatedText, nil
}

func StartGeminiChat(safeOn bool) {
	ctx := context.Background()
	geminiKey, err := keyring.Get("genie", "gemini_api_key")
	if err != nil {
		fmt.Println("Gemini API key not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  geminiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	modelName := "gemini-2.5-flash" // default model
	if selectedModel, err := keyring.Get("genie", "modelName"); err == nil {
		modelName = selectedModel
	}

	config := getSafetyConfig(safeOn)
	chat, err := client.Chats.Create(ctx, modelName, config, nil)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	setupChatStyles()
	startChatSession(ctx, chat, scanner)
}

func getSafetyConfig(safeOn bool) *genai.GenerateContentConfig {
	config := &genai.GenerateContentConfig{}
	if safeOn {
		config.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: genai.HarmBlockThresholdBlockLowAndAbove,
			},
		}
	} else {
		config.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThresholdBlockNone,
			},
		}
	}
	return config
}

var (
	style          lipgloss.Style
	promptStyle    lipgloss.Style
	aiStyle        lipgloss.Style
	multilineStyle lipgloss.Style
)

func setupChatStyles() {
	style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#9D4EDD"))

	promptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06D6A0")).
		Bold(true)

	aiStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#118AB2")).
		PaddingLeft(2)

	multilineStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)
}

func startChatSession(ctx context.Context, chat *genai.Chat, scanner *bufio.Scanner) {
	// Configure readline with multiline support
	config := &readline.Config{
		Prompt:                 promptStyle.Render("You 💭 > "),
		HistoryFile:            "/tmp/genie_gemini_history",
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
			// Clear terminal screen
			fmt.Print("\033[H\033[2J")
			// Reprint welcome message
			color.New(color.FgHiMagenta).Println("🧞 Chat session started!")
			fmt.Println(style.Render("Commands: 'exit' | 'clear' | '/history' | '/email'"))
			fmt.Println(multilineStyle.Render("Tip: For multiline input, type '\\' at end of line or use '---' on a new line to send."))
			fmt.Println(strings.Repeat("─", 50))
			continue
		case constants.HistoryCommand:
			exportGeminiChatHistory(chat.History(true))
			continue
		case constants.EmailCommand:
			emailChatHistory(chat.History(true))
			continue
		}

		handleChatMessage(ctx, chat, userInput)
	}
}

func handleChatMessage(ctx context.Context, chat *genai.Chat, userInput string) {
	s := spinner.New(spinner.CharSets[11], 80*time.Millisecond)
	s.Prefix = color.HiCyanString("🤔 Thinking: ")
	s.Suffix = " Please wait..."
	s.Start()

	genResp, err := chat.SendMessage(ctx, genai.Part{Text: userInput})
	s.Stop()

	if err != nil {
		log.Fatal(err)
	}

	response := genResp.Text()
	if response != "" {
		fmt.Print(color.HiCyanString("\n🤖 AI: "))

		paragraphs := strings.Split(response, "\n")
		for _, p := range paragraphs {
			if p != "" {
				fmt.Println(aiStyle.Render(p))
			}
		}
		fmt.Println(strings.Repeat("─", 50))
	}
}

func exportGeminiChatHistory(history []*genai.Content) {
	if len(history) == 0 {
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

	for _, msg := range history {
		var text string
		if len(msg.Parts) > 0 && msg.Parts[0].Text != "" {
			text = msg.Parts[0].Text
		}
		switch msg.Role {
		case "user":
			content.WriteString(fmt.Sprintf("### 💭 You\n%s\n\n", text))
		case "model":
			content.WriteString(fmt.Sprintf("### 🤖 AI\n%s\n\n", text))
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

func emailChatHistory(history []*genai.Content) {
	if len(history) == 0 {
		fmt.Printf("%s No chat history available to email.\n", color.RedString("❌"))
		return
	}

	// Create a divider for visual separation
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println(color.HiMagentaString("📧 Emailing Chat History"))
	fmt.Println(strings.Repeat("─", 50))

	// Get current model name
	modelName := "gemini-2.5-flash"
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

	var messages []map[string]string
	for _, msg := range history {
		var text string
		if len(msg.Parts) > 0 && msg.Parts[0].Text != "" {
			text = msg.Parts[0].Text
		}
		messages = append(messages, map[string]string{
			"role":    msg.Role,
			"content": text,
		})
	}

	payload := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"model":     modelName,
		"messages":  messages,
		"metadata": map[string]string{
			"sessionId": fmt.Sprintf("gemini-%d", time.Now().Unix()),
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
