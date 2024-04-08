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
			prompt = "You are a genie, someone is expecting a greeting from you. They are saying" + args[0] + "Respond to them. No limited wishes. You are a CLI genie you will help the person to use CLI in a smarter way. Match the context. This is just a greeting so make it a one-liner."
		} else {
			prompt = "You are a genie, someone is expecting a greeting from you. Respond to them. No limited wishes. You are a CLI genie you will help the person to use CLI in a smarter way. Match the context. This is just a greeting so make it a one-liner."
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
