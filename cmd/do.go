package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/helpers"
	"github.com/harshalranjhani/genie/structs"
	"github.com/spf13/cobra"
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
		}

		safeSettings, _ := cmd.Flags().GetBool("safe")
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		rootDir, err := helpers.GetCurrentDirectoriesAndFiles(dir)
		if err != nil {
			panic(err)
		}
		var sb strings.Builder
		helpers.PrintData(&sb, rootDir, 0)

		var prompt string = fmt.Sprintf("Context: You are an intelligent CLI tool named Genie, designed to understand and execute file system operations based on the current state of the user's directory and explicit instructions provided. Your responses must strictly be executable commands suitable for a Unix-like shell, without any additional explanations, comments, or output.\n\nCurrent Directory Snapshot:\n---------------------------\n%s\n\nTask:\n-----\nBased on the above directory snapshot, execute the operation specified by the user's request encapsulated in 'args[0]'. 'args[0]' contains the explicit instruction for a file system operation that needs to be performed on the current directory or its contents.\n\nNote: The command you provide will be run directly in a Unix-like shell environment. Ensure your command is syntactically correct and contextually appropriate for the operation described in 'args[0]'. Your response should consist only of the command necessary to perform the operation, with no additional text.\n\nRequested Operation: %s\nProvide the Command, if you can't match the context or find a similar command, just echo that to the terminal", sb.String(), args[0])

		if safeSettings {
			color.Green("Safety settings are on.")
		} else {
			color.Red("Safety settings are off.")
		}

		resp, err := helpers.GetResponse(prompt, safeSettings)
		if err != nil {
			log.Fatal(err)
		}
		respJSON, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			log.Fatal("Error marshalling response to JSON:", err)
		}

		// Unmarshal the JSON response into the struct
		var genResp structs.GenResponse
		err = json.Unmarshal(respJSON, &genResp)
		if err != nil {
			log.Fatal("Error unmarshalling response JSON:", err)
		}

		if len(genResp.Candidates) > 0 && len(genResp.Candidates[0].Content.Parts) > 0 {
			generatedText := genResp.Candidates[0].Content.Parts[0]
			// the generatedText is the command to be executed, so we need to run it
			fmt.Println("Running the command: ", generatedText)
			helpers.RunCommand(generatedText)
		} else {
			fmt.Println("No generated text found")
		}
	},
}
