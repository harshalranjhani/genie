package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "genie",
	Short: "genie is an AI powered CLI tool to help you with your daily tasks.",
	Long:  `genie is an AI powered CLI tool to help you with your daily tasks.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// Do Stuff Here
	// },
}

func Execute() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		funcMap := template.FuncMap{
			"rpad":                    rightPad,
			"trimTrailingWhitespaces": trimTrailingWhitespaces,
			"addEmoji":                addEmoji,
		}

		tmpl := helpTemplate
		if cmd != rootCmd {
			tmpl = subcommandHelpTemplate
		}

		t, err := template.New("help").Funcs(funcMap).Parse(tmpl)
		if err != nil {
			fmt.Println(err)
			return
		}

		var out bytes.Buffer
		err = t.Execute(&out, cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		if cmd == rootCmd {
			myFigure := figure.NewFigure("genie", "", true)
			asciiArt := myFigure.String()
			lines := bytes.Split([]byte(asciiArt), []byte("\n"))
			for _, line := range lines {
				c := color.New(color.FgCyan)
				c.Println(string(line))
			}
		}

		fmt.Println(out.String())
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
