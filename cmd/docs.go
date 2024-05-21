package cmd

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(docsCmd)
}

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "open the documentation",
	Long:  `Open the documentation in your default web browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		url := "https://genie.harshalranjhani.in/setup"
		err := browser.OpenURL(url)
		if err != nil {
			fmt.Println("Failed to open URL:", err)
		} else {
			fmt.Println("Documentation: https://genie.harshalranjhani.in/setup")
		}
	},
}
