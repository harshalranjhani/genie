package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/helpers"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(doCmd)
	doCmd.PersistentFlags().Bool("safe", false, "Set this to true if you wish to enable safe mode.")
}

var doCmd = &cobra.Command{
	Use:   "do `prompt in quotes`",
	Short: "command the genie to do something",
	Long:  `This is a command to instruct the genie to do something.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		safeSettings, _ := cmd.Flags().GetBool("safe")
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		rootDir, err := helpers.GetCurrentDirectoriesAndFiles(dir)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		var sb strings.Builder
		helpers.PrintData(&sb, rootDir, 0)
		var prompt string = fmt.Sprintf("Context: You are an intelligent CLI tool named Genie, designed to understand and execute file system operations based on the current state of the user's directory and explicit instructions provided. Your responses must strictly be executable commands suitable for a Unix-like shell, without any additional explanations, comments, or output.\n\nCurrent Directory Snapshot:\n---------------------------\n%s\n\nTask:\n-----\nBased on the above directory snapshot, execute the operation specified by the user's request encapsulated in 'args[0]'. 'args[0]' contains the explicit instruction for a file system operation that needs to be performed on the current directory or its contents.\n\nNote: The command you provide will be run directly in a Unix-like shell environment. Ensure your command is syntactically correct and contextually appropriate for the operation described in 'args[0]'. Your response should consist only of the command necessary to perform the operation, with no additional text.\n\nRequested Operation: %s\nProvide the Command, if you can't match the context or find a similar command, just echo that to the terminal. The Operating System of the User is: %s", sb.String(), args[0], runtime.GOOS)

		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		if safeSettings {
			color.Green("Safety settings are on.")
			if engineName == GPTEngine {
				color.Red("Safety settings are low by default for GPT engine.")
			}
		} else {
			color.Red("Safety settings are off.")
		}

		switch engineName {
		case GPTEngine:
			err := helpers.GetGPTCmdResponse(prompt, true)
			if err != nil {
				log.Fatal(err)
			}
		case GeminiEngine:
			err := helpers.GetGeminiCmdResponse(prompt, safeSettings)
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("Unknown engine name: ", engineName)
		}

		if err != nil {
			log.Fatal(err)
		}

	},
}
