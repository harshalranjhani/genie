package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/config"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/internal/helpers/llm"
	"github.com/harshalranjhani/genie/pkg/prompts"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(tellCmd)
	tellCmd.PersistentFlags().Bool("include-dir", false, "Option to include the current directory snapshot in the request.")
	tellCmd.PersistentFlags().Bool("include-git-changes", false, "Option to include git repository information in the request.")
}

var tellCmd = &cobra.Command{
	Use:   "tell",
	Short: "This is a command to seek help from the genie",
	Long:  `This is a command to seek help from the genie. For example: 'genie tell "what is docker?"'`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			color.Red("Please provide a question for the genie")
			return
		}

		prompt := args[0]
		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		_, exists := config.GetEngine(engineName)
		if !exists {
			log.Fatal("Unknown engine name: ", engineName)
		}

		includeDir, _ := cmd.Flags().GetBool("include-dir")
		includeGit, _ := cmd.Flags().GetBool("include-git-changes")

		var sb strings.Builder

		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if includeDir {
			rootDir, err := helpers.GetCurrentDirectoriesAndFiles(dir)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			helpers.PrintData(&sb, rootDir, 0)
		}

		if includeGit {
			gitInfo, err := helpers.GetGitInfo(dir)
			if err != nil {
				color.Red("Warning: Could not get git information: %v", err)
			} else {
				sb.WriteString("\nGit Repository Information:\n")
				sb.WriteString(gitInfo)
			}
		}

		prompt = prompts.GetTellPrompt(prompt, sb)

		switch engineName {
		case config.GPTEngine:
			llm.GetGPTGeneralResponse(prompt, includeDir)
		case config.GeminiEngine:
			strResp, err := llm.GetGeminiGeneralResponse(prompt, true, includeDir)
			if err != nil {
				log.Fatal("Error getting response from Gemini: ", err)
				os.Exit(1)
			}
			fmt.Println(formatMarkdownToPlainText(strResp))
		case config.DeepSeekEngine:
			err := llm.GetDeepSeekGeneralResponse(prompt, true, includeDir)
			if err != nil {
				log.Fatal("Error getting response from DeepSeek: ", err)
				os.Exit(1)
			}
		case config.OllamaEngine:
			model, err := keyring.Get(serviceName, modelAccountKey)
			if err != nil {
				log.Fatal("Error retrieving model name from keyring: ", err)
			}
			err = llm.GetOllamaGeneralResponse(prompt, model, includeDir)
			if err != nil {
				log.Fatal("Error getting response from Ollama: ", err)
				os.Exit(1)
			}
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
