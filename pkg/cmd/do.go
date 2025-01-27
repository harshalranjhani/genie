package cmd

import (
	"log"
	"os"
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
	rootCmd.AddCommand(doCmd)
	doCmd.PersistentFlags().Bool("safe", false, "Set this to true if you wish to enable safe mode.")
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Command the genie to do something",
	Long:  `Command the genie to do something. For example: 'genie do "list all files in the current directory"'`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			color.Red("Please provide a command for the genie")
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

		engine, exists := config.GetEngine(engineName)
		if !exists {
			log.Fatal("Unknown engine name: ", engineName)
		}

		if safeSettings {
			color.Green("Safety settings are on.")
			if !engine.Features.SupportsSafeMode {
				color.Red("Currently %s does not support safe mode. But we're still instructing it to be extra cautious for this particular request.", engineName)
				prompt += " Please ensure the command is safe and does not contain any destructive behavior like deleting files, directories, etc. If it does, please reject it and just echo why you rejected"
			}
		} else {
			color.Red("Safety settings are off.")
		}

		switch engineName {
		case config.GPTEngine:
			err := llm.GetGPTCmdResponse(prompt, true)
			if err != nil {
				log.Fatal(err)
			}
		case config.GeminiEngine:
			err := llm.GetGeminiCmdResponse(prompt, safeSettings)
			if err != nil {
				log.Fatal(err)
			}
		case config.DeepSeekEngine:
			err := llm.GetDeepSeekCmdResponse(prompt, safeSettings)
			if err != nil {
				log.Fatal(err)
			}
		case config.OllamaEngine:
			model, err := keyring.Get(serviceName, modelAccountKey)
			if err != nil {
				log.Fatal("Error retrieving model name from keyring: ", err)
			}
			err = llm.GetOllamaCmdResponse(prompt, model, safeSettings)
			if err != nil {
				log.Fatal(err)
			}
		}

		if err != nil {
			log.Fatal(err)
		}

	},
}
