package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/harshalranjhani/genie/helpers"
	"github.com/harshalranjhani/genie/structs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tellCmd)
}

var tellCmd = &cobra.Command{
	Use:   "tell `prompt in quotes`",
	Short: "This is a command to seek help from the genie.",
	Long:  `Ask the genie about anything you need help with related to CLI issues and queries within UNIX or any other shell environment. Focus on troubleshooting, script writing, command explanations, and system configurations. Avoid discussing unrelated topics.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}

		var prompt string = fmt.Sprintf("Context: You are an intelligent CLI tool named Genie, designed to understand and execute file system operations based on the current state of the user's directory and explicit instructions provided. Please provide assistance strictly related to command-line interface (CLI) issues and queries within UNIX or any other shell environment. Focus on troubleshooting, script writing, command explanations, and system configurations. Avoid discussing unrelated topics.\nHere's what the user is asking %s", args[0])

		resp, err := helpers.GetResponse(prompt, true)
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
			fmt.Println(formatMarkdownToPlainText(generatedText))
		} else {
			fmt.Println("No generated text found")
		}
	},
}

func formatMarkdownToPlainText(mdText string) string {
	// Regular expressions to replace Markdown formatting
	reStrong := regexp.MustCompile(`\*\*(.*?)\*\*`)
	reEmphasis := regexp.MustCompile(`\*(.*?)\*`)
	reCode := regexp.MustCompile("([^])" + "`" + "(.*?)" + "`" + "([^`])")
	reHeaders := regexp.MustCompile(`\n#+\s(.*?)\n`)

	// Replace Markdown syntax with plain text formatting
	plainText := reStrong.ReplaceAllString(mdText, "$1")
	plainText = reEmphasis.ReplaceAllString(plainText, "$1")
	plainText = reCode.ReplaceAllString(plainText, "$1")
	plainText = reHeaders.ReplaceAllString(plainText, "\n$1\n")

	return plainText
}
