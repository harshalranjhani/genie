package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

func rightPad(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func trimTrailingWhitespaces(s string) string {
	return strings.TrimRight(s, " \t\n")
}

const helpTemplate = `genie is an AI powered CLI tool to help you with your daily tasks.

Usage:
  genie [command]

Available Commands:
{{range .Commands}}{{if (and .IsAvailableCommand (not .IsAdditionalHelpTopicCommand))}}  [{{.Name}}] {{.Short}}
{{end}}{{end}}
{{if .HasAvailableLocalFlags}}Flags:
{{.LocalFlags.FlagUsages}}

{{end}}Use "genie [command] --help" for more information about a command.
`

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Help command",
	Long:  `Help about genie.`,
	// Run: func(cmd *cobra.Command, args []string) {

	// },
}
