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
		}
		tmpl, err := template.New("help").Funcs(funcMap).Parse(helpTemplate)
		if err != nil {
			fmt.Println(err)
			return
		}
		var out bytes.Buffer
		err = tmpl.Execute(&out, cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		myFigure := figure.NewFigure("genie", "", true)
		asciiArt := myFigure.String()

		lines := bytes.Split([]byte(asciiArt), []byte("\n"))
		for _, line := range lines {
			c := color.New(color.FgCyan)
			c.Println(string(line))
		}

		fmt.Print(out.String())
		color.Cyan("Additionally, you can visit https://docs.genie.harshalranjhani.in for a detailed documentation.")
	})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
