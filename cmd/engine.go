package cmd

import (
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(engineCmd)
}

var engineCmd = &cobra.Command{
	Use:   "engine",
	Short: "Get the current engine being used by genie",
	Long:  `This command will return the current engine being used by genie.`,
	Run: func(cmd *cobra.Command, args []string) {
		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}
		color.Green("Current engine: %s", engineName)
	},
}
