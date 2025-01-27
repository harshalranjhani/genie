package cmd

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// addEmoji returns an appropriate emoji for each command
func addEmoji(cmdName string) string {
	emojiMap := map[string]string{
		"chat":       "💬",
		"completion": "🔄",
		"do":         "🎯",
		"docs":       "📚",
		"document":   "📝",
		"engine":     "⚙️",
		"generate":   "🎨",
		"greet":      "👋",
		"help":       "❓",
		"image":      "🖼️",
		"init":       "🔧",
		"music":      "🎵",
		"readme":     "📖",
		"reset":      "🔄",
		"scrape":     "🕸️",
		"status":     "📊",
		"summarize":  "📊",
		"support":    "❤️",
		"switch":     "🔀",
		"tell":       "💭",
		"use":        "🎨",
		"verify":     "✅",
		"version":    "📌",
		"bug":        "🐛",
	}

	if emoji, ok := emojiMap[cmdName]; ok {
		return emoji
	}
	return "•"
}

// Define template functions
var templateFuncs = template.FuncMap{
	"addEmoji":                addEmoji,
	"rightPad":                rightPad,
	"trimTrailingWhitespaces": trimTrailingWhitespaces,
}

const helpTemplate = `
✨ Genie AI Assistant
────────────────────
Your AI-powered CLI companion for daily tasks

🎯 Usage
────────
  genie [command]

📚 Available Commands
───────────────────{{range .Commands}}{{if (and .IsAvailableCommand (not .IsAdditionalHelpTopicCommand))}}
  {{.Name | printf "%-12s"}} {{addEmoji .Name}} {{.Short}}{{end}}{{end}}

🛠️  Common Examples
─────────────────
  • Start a chat session:
    $ genie chat

  • Initialize your API keys:
    $ genie init

  • Reset your configuration:
    $ genie reset

  • Get command help:
    $ genie [command] --help

{{if .HasAvailableLocalFlags}}🚩 Flags
─────────
{{.LocalFlags.FlagUsages}}{{end}}
💡 Tips
───────
  • Use arrow keys to navigate chat history
  • Press Ctrl+C to exit any command
  • Type 'clear' to clear the chat history
  • Commands are case-sensitive

📖 Learn More
────────────
  • Documentation: https://docs.genie.harshalranjhani.in
  • Report issues: https://github.com/harshalranjhani/genie/issues
`

// Add this constant after the helpTemplate
const subcommandHelpTemplate = `
{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}
🎯 Usage
────────
  {{if .Runnable}}{{.CommandPath}} {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}

{{if .HasAvailableSubCommands}}📚 Available Commands
───────────────────{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name " " 15}} {{addEmoji .Name}} {{.Short}}{{end}}{{end}}

{{end}}{{if .HasAvailableLocalFlags}}🚩 Flags
─────────
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

{{end}}{{if .HasAvailableInheritedFlags}}🔄 Global Flags
──────────────
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}

{{end}}{{if .HasExample}}📝 Examples
──────────
{{.Example}}

{{end}}💡 Additional Help
───────────────
  Use "genie [command] --help" for more information about a command.`

func init() {
	// Register template functions
	cobra.AddTemplateFuncs(templateFuncs)

	// Set the main help template
	rootCmd.SetHelpTemplate(helpTemplate)

	// Set the subcommand help template for all commands
	for _, cmd := range rootCmd.Commands() {
		cmd.SetHelpTemplate(subcommandHelpTemplate)
	}

	rootCmd.AddCommand(helpCmd)
}

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help and usage information",
	Long: color.YellowString(`
✨ Genie Help Center
───────────────────
Get detailed information about Genie commands and features.
Use 'genie [command] --help' for more details about a specific command.
`),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Root().Help()
			return
		}

		// Show help for specific command
		targetCmd, _, err := cmd.Root().Find(args)
		if err != nil {
			color.Red("\n❌ Unknown command: %s", args[0])
			fmt.Println(color.YellowString("\nAvailable commands:"))
			cmd.Root().Help()
			return
		}

		targetCmd.Help()
	},
}

// rightPad adds padding to the right of a string
func rightPad(s string, padStr string, overallLen int) string {
	if len(s) >= overallLen {
		return s
	}
	padCount := 1 + ((overallLen - len(s)) / len(padStr))
	retStr := s + strings.Repeat(padStr, padCount)
	return retStr[:overallLen]
}

// trimTrailingWhitespaces removes trailing whitespace
func trimTrailingWhitespaces(s string) string {
	return strings.TrimRight(s, " \t\n")
}
