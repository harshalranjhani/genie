package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/internal/helpers/llm"
	"github.com/harshalranjhani/genie/pkg/prompts"
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
		prompt := prompts.GetDoPrompt(sb, args[0])

		engineName, err := keyring.Get(serviceName, "engineName")
		if err != nil {
			log.Fatal("Error retrieving engine name from keyring:", err)
		}

		if safeSettings {
			color.Green("Safety settings are on.")
			if engineName == GPTEngine {
				color.Red("Safety settings are low by default for GPT engine.")
			} else if engineName == DeepSeekEngine {
				color.Red("Currently DeepSeek does not support safe mode. But we're still instructing it to be extra cautious for this particular request.")
				prompt += " Please ensure the command is safe and does not contain any destructive behavior like deleting files, directories, etc. If it does, please reject it and just echo why you rejected"
			}
		} else {
			color.Red("Safety settings are off.")
		}

		switch engineName {
		case GPTEngine:
			err := llm.GetGPTCmdResponse(prompt, true)
			if err != nil {
				log.Fatal(err)
			}
		case GeminiEngine:
			err := llm.GetGeminiCmdResponse(prompt, safeSettings)
			if err != nil {
				log.Fatal(err)
			}
		case DeepSeekEngine:
			err := llm.GetDeepSeekCmdResponse(prompt, safeSettings)
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
