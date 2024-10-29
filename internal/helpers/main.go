package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func GetGitInfo(path string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		color.Red("Error: Git repository not found in current directory")
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	var info strings.Builder
	hasError := false

	// Get current branch
	head, err := repo.Head()
	if err != nil {
		color.Yellow("Warning: Could not get current branch information")
		hasError = true
	} else {
		info.WriteString(fmt.Sprintf("Current Branch: %s\n", head.Name().Short()))
	}

	// Get status using git command
	status, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		color.Yellow("Warning: Could not get git status information")
		hasError = true
	} else if len(status) > 0 {
		info.WriteString("\nUncommitted Changes:\n")
		info.WriteString(string(status))

		// Get diff with size limit and histogram
		diffStats, err := exec.Command("git", "diff", "--stat").Output()
		if err != nil {
			color.Yellow("Warning: Could not get git diff statistics")
		} else {
			info.WriteString("\nDiff Statistics:\n")
			info.WriteString(string(diffStats))
		}

		diff, err := exec.Command("git", "diff").Output()
		if err != nil {
			color.Yellow("Warning: Could not get git diff information")
		} else {
			diffStr := string(diff)
			if len(diffStr) > 12000 {
				info.WriteString("\nDiff (truncated - too large):\n")
				info.WriteString(diffStr[:12000])
				info.WriteString("\n... (diff truncated, see statistics above for full change overview)")
			} else {
				info.WriteString("\nDiff:\n")
				info.WriteString(diffStr)
			}
		}
	}

	// Get recent commits only if we have head reference
	if head != nil {
		commits, err := repo.Log(&git.LogOptions{
			From:  head.Hash(),
			Order: git.LogOrderCommitterTime,
			All:   false,
		})
		if err != nil {
			color.Yellow("Warning: Could not get commit history")
			hasError = true
		} else {
			info.WriteString("\nRecent Commits:\n")
			count := 0
			err = commits.ForEach(func(c *object.Commit) error {
				if count >= 5 {
					return fmt.Errorf("done")
				}
				info.WriteString(fmt.Sprintf("- %s: %s\n", c.Hash.String()[:7], c.Message))
				count++
				return nil
			})
			if err != nil && err.Error() != "done" {
				color.Yellow("Warning: Error while reading commit history")
				hasError = true
			}
		}
	}

	if hasError {
		return info.String(), fmt.Errorf("completed with some errors")
	}
	return info.String(), nil
}
