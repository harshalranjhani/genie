package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/pkg/prompts"
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

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(jsonData))
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

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(jsonData))
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

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(jsonData))
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
