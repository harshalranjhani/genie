package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/harshalranjhani/genie/pkg/assets"
)

func SanitizeUTF8(s string) string {
	return strings.Map(func(r rune) rune {
		if r == utf8.RuneError {
			return -1
		}
		return r
	}, s)
}

func ExtractKeyValuePairs(text string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}

	return result
}

func ProcessTemplateResponse(templateName string, generatedText string, readmePath string) error {
	lines := strings.Split(generatedText, "\n")
	if len(lines) > 2 && strings.HasPrefix(lines[0], "```") && strings.HasSuffix(lines[len(lines)-1], "```") {
		generatedText = strings.Join(lines[1:len(lines)-1], "\n")
	}

	var readmeData map[string]interface{}
	err := json.Unmarshal([]byte(generatedText), &readmeData)
	if err != nil {
		return fmt.Errorf("failed to parse generated JSON: %w", err)
	}

	tmpl, err := template.New("readme").Parse(assets.GetReadmeTemplate(templateName))
	if err != nil {
		return fmt.Errorf("failed to parse README template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, readmeData); err != nil {
		return fmt.Errorf("failed to execute README template: %w", err)
	}
	if err := os.WriteFile(readmePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write README file: %w", err)
	}
	return nil
}
