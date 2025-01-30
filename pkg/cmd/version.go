package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "get the current version of genie",
	Long:  `All software has versions. This is genie's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("genie", Version)
	},
}

// Version represents the current version of genie
const Version = "v2.9.4"
