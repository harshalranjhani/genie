package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/harshalranjhani/genie/helpers"
	"github.com/harshalranjhani/genie/structs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use:   "do `prompt in quotes`",
	Short: "command the genie to do something",
	Long:  `This is a command to instruct the genie to do something.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}

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

		prompt := "Here is a list of all my files and directories in the folder I am in right now:\n" + sb.String() + "\n I need you to " + args[0] + "\n Just give the command and nothing else since im going to run that directly."
		resp, err := helpers.GetResponse(prompt)
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
