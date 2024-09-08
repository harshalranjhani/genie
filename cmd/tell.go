package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/harshalranjhani/genie/helpers"
	"github.com/harshalranjhani/genie/helpers/llm"
	"github.com/harshalranjhani/genie/helpers/prompts"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(tellCmd)
	tellCmd.PersistentFlags().Bool("include-dir", false, "Option to include the current directory snapshot in the request.")
}

var tellCmd = &cobra.Command{
	Use:   "tell `prompt in quotes`",
	Short: "This is a command to seek help from the genie.",
	Long:  `Ask the genie about anything you need help with related to CLI issues and queries within UNIX or any other shell environment. Focus on troubleshooting, script writing, command explanations, and system configurations. Avoid discussing unrelated topics. The Operating System of the User is: ` + runtime.GOOS,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		includeDir, _ := cmd.Flags().GetBool("include-dir")

		var prompt string

		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		if includeDir {
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
			prompt = prompts.GetTellPrompt(args[0], sb)
		} else {
			prompt = prompts.GetTellPrompt(args[0], strings.Builder{})
		}

		switch engineName {
		case GPTEngine:
			llm.GetGPTGeneralResponse(prompt, includeDir)
		case GeminiEngine:
			strResp, err := llm.GetGeminiGeneralResponse(prompt, true, includeDir)
			if err != nil {
				log.Fatal("Error getting response from Gemini: ", err)
				os.Exit(1)
			}
			fmt.Println(formatMarkdownToPlainText(strResp))
		default:
			log.Fatal("Unknown engine name: ", engineName)
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
