package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

func init() {
	rootCmd.AddCommand(greetCmd)
}

type GenResponse struct {
	Candidates []struct {
		Content struct {
			Parts []string `json:"Parts"`
		} `json:"Content"`
	} `json:"Candidates"`
}

var greetCmd = &cobra.Command{
	Use:   "greet",
	Short: "fun greet genie command",
	Long:  `This is a fun greet genie command`,
	Run: func(cmd *cobra.Command, args []string) {
		godotenv.Load()
		ctx := context.Background()
		// fmt.Println("API Key:", os.Getenv("GEMINI_API_KEY"))
		client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		model := client.GenerativeModel("gemini-pro")

		prompt := ""
		if len(args) > 0 {
			prompt = "Imagine you are an ancient and wise genie, residing not in a lamp, but within the heart of a powerful computer's Command Line Interface (CLI). After centuries of slumber, a user awakens you with a command. If the user includes a specific greeting or request (" + (args[0]) + "), they're seeking your ancient wisdom to navigate the complexities of the CLI more efficiently. They might say something like 'Hello, Genie, how can I list all files in this directory?' Respond with a greeting that reflects your vast knowledge and eagerness to assist in the digital realm. No wish is too complex for you to grant, especially when it comes to mastering the CLI. Craft your response as a one-liner, demonstrating your readiness to offer sage advice and practical tips for smarter CLI usage. Remember, your goal is to match the context of their inquiry or greeting, providing a response that blends the mystical with the practical."
		} else {
			prompt = "Imagine you are an ancient and wise genie, residing not in a lamp, but within the heart of a powerful computer's Command Line Interface (CLI). After centuries of slumber, a user awakens you with a command. If the user includes a specific greeting or request, they're seeking your ancient wisdom to navigate the complexities of the CLI more efficiently. They might say something like 'Hello, Genie, how can I list all files in this directory?' Respond with a greeting that reflects your vast knowledge and eagerness to assist in the digital realm. No wish is too complex for you to grant, especially when it comes to mastering the CLI. Craft your response as a one-liner, demonstrating your readiness to offer sage advice and practical tips for smarter CLI usage. Remember, your goal is to match the context of their inquiry or greeting, providing a response that blends the mystical with the practical."
		}
		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			log.Fatal(err)
		}
		respJSON, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			log.Fatal("Error marshalling response to JSON:", err)
		}

		// Unmarshal the JSON response into the struct
		var genResp GenResponse
		err = json.Unmarshal(respJSON, &genResp)
		if err != nil {
			log.Fatal("Error unmarshalling response JSON:", err)
		}

		if len(genResp.Candidates) > 0 && len(genResp.Candidates[0].Content.Parts) > 0 {
			generatedText := genResp.Candidates[0].Content.Parts[0]
			fmt.Println(generatedText)
		} else {
			fmt.Println("No generated text found")
		}

	},
}
