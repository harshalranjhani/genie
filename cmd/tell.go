package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/harshalranjhani/genie/helpers"
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
			prompt = fmt.Sprintf(
				"Context: You are an intelligent CLI tool named Genie, designed to understand and execute file system operations based on the current state of the user's directory and explicit instructions provided. Please provide assistance strictly related to command-line interface (CLI) issues and queries within UNIX or any other shell environment and any other thing related to the field of Computer Science. Focus on troubleshooting, script writing, command explanations, and system configurations. Avoid discussing unrelated topics.\n\n"+
					"Also, if someone asks about what all you can do other than this, here is the help command for genie:\n"+
					"Usage:\n"+
					"  genie [command]\n\n"+
					"Available Commands:\n"+
					"  [chat]       Start a chat with the genie and maintain a conversation.\n"+
					"  [completion] Generate the autocompletion script for the specified shell.\n"+
					"  [do]        Command the genie to do something.\n"+
					"  [docs]      Open the documentation.\n"+
					"  [document] Document your code with genie.\n"+
					"  [engine]    Get the current engine being used by genie.\n"+
					"  [generate]  Generate an image from a prompt.\n"+
					"  [greet]     Invoke the wise Genie for CLI guidance.\n"+
					"  [init]      Store your API keys securely in the system keychain.\n"+
					"  [music]     Generate music from text!\n"+
					"  [reset]     Reset your API keys.\n"+
					"  [scrape]    Scrape data from a URL, supports pagination!\n"+
					"  [summarize] Get a markdown summary of the current directory comments.\n"+
					"  [support]   Support the tool by donating to the project.\n"+
					"  [switch]    Switch between different engines (Gemini, GPT).\n"+
					"  [tell]      This is a command to seek help from the genie.\n"+
					"  [verify]    Verify your support status and get access to extra features.\n"+
					"  [version]   Get the current version of genie.\n\n"+
					"Flags:\n"+
					"  -h, --help   help for genie\n\n"+
					"Use \"genie [command] --help\" for more information about a command.\n"+
					"Additionally, you can visit https://genie.harshalranjhani.in/docs for a detailed documentation.\n\n"+
					"Here's what the user is asking: %s\n"+
					"The user has also provided the current directory's snapshot:\n\n"+
					"Current Directory Snapshot:\n"+
					"---------------------------\n"+
					"%s. Use it if required.", args[0], sb.String())
		} else {
			prompt = fmt.Sprintf(
				"Context: You are an intelligent CLI tool named Genie, designed to understand and execute file system operations based on the current state of the user's directory and explicit instructions provided. Please provide assistance strictly related to command-line interface (CLI) issues and queries within UNIX or any other shell environment and any other thing related to the field of Computer Science. Focus on troubleshooting, script writing, command explanations, and system configurations. Avoid discussing unrelated topics.\n\n"+
					"Also, if someone asks about what all you can do other than this, here is the help command for genie:\n"+
					"Usage:\n"+
					"  genie [command]\n\n"+
					"Available Commands:\n"+
					"  [chat]       Start a chat with the genie and maintain a conversation.\n"+
					"  [completion] Generate the autocompletion script for the specified shell.\n"+
					"  [do]        Command the genie to do something.\n"+
					"  [docs]      Open the documentation.\n"+
					"  [document] Document your code with genie.\n"+
					"  [engine]    Get the current engine being used by genie.\n"+
					"  [generate]  Generate an image from a prompt.\n"+
					"  [greet]     Invoke the wise Genie for CLI guidance.\n"+
					"  [init]      Store your API keys securely in the system keychain.\n"+
					"  [music]     Generate music from text!\n"+
					"  [reset]     Reset your API keys.\n"+
					"  [scrape]    Scrape data from a URL, supports pagination!\n"+
					"  [summarize] Get a markdown summary of the current directory comments.\n"+
					"  [support]   Support the tool by donating to the project.\n"+
					"  [switch]    Switch between different engines (Gemini, GPT).\n"+
					"  [tell]      This is a command to seek help from the genie.\n"+
					"  [verify]    Verify your support status and get access to extra features.\n"+
					"  [version]   Get the current version of genie.\n\n"+
					"Flags:\n"+
					"  -h, --help   help for genie\n\n"+
					"Use \"genie [command] --help\" for more information about a command.\n"+
					"Additionally, you can visit https://genie.harshalranjhani.in/docs for detailed documentation.\n\n"+
					"Here's what the user is asking: %s", args[0])
		}

		switch engineName {
		case GPTEngine:
			helpers.GetGPTGeneralResponse(prompt, includeDir)
		case GeminiEngine:
			strResp, err := helpers.GetGeminiGeneralResponse(prompt, true, includeDir)
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
