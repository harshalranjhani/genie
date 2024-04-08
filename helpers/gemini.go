package helpers

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func GetResponse(prompt string) (*genai.GenerateContentResponse, error) {
	godotenv.Load()
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	return resp, err
}
