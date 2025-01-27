package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
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
			"temperature": 0.7,
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
