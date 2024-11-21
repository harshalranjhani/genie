package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/helpers/llm"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(bugCmd)
	bugCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringP("severity", "s", "medium", "Bug severity (low, medium, high, critical)")
	reportCmd.Flags().StringP("category", "c", "", "Bug category (ui, backend, security, performance, etc.)")
	reportCmd.Flags().StringP("assignee", "a", "", "Who should be assigned to this bug")
	reportCmd.Flags().StringP("priority", "p", "medium", "Bug priority (low, medium, high)")

	bugCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return fmt.Errorf("unknown subcommand %q for %q\nDid you mean 'bug report'?", args[0], cmd.CommandPath())
	}
}

var bugCmd = &cobra.Command{
	Use:   "bug",
	Short: "Manage bug reports",
	Long: `Create and manage detailed bug reports with AI assistance.

Available Commands:
  report      Create a new bug report with description

Usage:
  genie bug report "bug description" [flags]

Example:
  genie bug report "Login button not working on Firefox"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var reportCmd = &cobra.Command{
	Use:   "report \"bug description\"",
	Short: "Create a new bug report",
	Long: `Create a detailed bug report from your description. 
Genie will analyze the issue and generate a structured report with:
- Steps to reproduce
- Expected vs actual behavior
- Potential fixes
- Severity and category classification

All the information is saved in a markdown file in the bugs directory categorized by priority order.`,
	Args: cobra.ExactArgs(1),
	Run:  runBugReport,
}

func runBugReport(cmd *cobra.Command, args []string) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing bug report: ")
	s.Start()

	description := args[0]
	severity, _ := cmd.Flags().GetString("severity")
	category, _ := cmd.Flags().GetString("category")
	assignee, _ := cmd.Flags().GetString("assignee")
	priority, _ := cmd.Flags().GetString("priority")

	// Get current time in local timezone instead of UTC
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05 MST")

	// Add timestamp to the beginning of the bug report template
	bugReportPrefix := fmt.Sprintf("# Bug Report Created: %s\n\n", formattedTime)

	engineName, err := keyring.Get(serviceName, "engineName")
	if err != nil {
		s.Stop()
		color.Red("Error retrieving engine name: %v", err)
		return
	}

	// Create bugs directory and priority subdirectory
	bugsDir := filepath.Join(".", "bugs")
	priorityDir := filepath.Join(bugsDir, strings.ToLower(priority))
	if err := os.MkdirAll(priorityDir, 0755); err != nil {
		s.Stop()
		color.Red("Error creating directory structure: %v", err)
		return
	}

	var bugReport string
	switch engineName {
	case GPTEngine:
		bugReport, err = llm.GenerateBugReportGPT(description, severity, category, assignee, priority)
	case GeminiEngine:
		bugReport, err = llm.GenerateBugReportGemini(description, severity, category, assignee, priority)
	default:
		s.Stop()
		color.Red("Unknown engine: %s", engineName)
		return
	}

	if err != nil {
		s.Stop()
		color.Red("Error generating bug report: %v", err)
		return
	}

	// Combine the timestamp with the generated report
	fullBugReport := bugReportPrefix + bugReport

	// Generate filename based on timestamp and category
	timestamp := time.Now().Format("20060102-150405")
	sanitizedCategory := strings.ToLower(strings.ReplaceAll(category, " ", "-"))
	if sanitizedCategory == "" {
		sanitizedCategory = "general"
	}

	filename := fmt.Sprintf("%s-%s-%s.md", timestamp, sanitizedCategory, severity)
	filepath := filepath.Join(priorityDir, filename)

	if err := os.WriteFile(filepath, []byte(fullBugReport), 0644); err != nil {
		s.Stop()
		color.Red("Error writing bug report: %v", err)
		return
	}

	s.Stop()
	color.Green("\n✓ Bug report generated successfully!")
	fmt.Printf("\nLocation: %s\n", color.CyanString(filepath))
}
