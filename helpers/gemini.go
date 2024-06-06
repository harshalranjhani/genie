package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/harshalranjhani/genie/structs"
	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"
	"google.golang.org/api/option"
)

func GetGeminiCmdResponse(prompt string, safeOn bool) error {
	godotenv.Load()
	ctx := context.Background()

	geminiKey, err := keyring.Get("genie", "gemini_api_key")
	if err != nil {
		fmt.Println("Gemini API key not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	if err != nil {
		return err
	}
	defer client.Close()
	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	if safeOn {
		model.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: genai.HarmBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: genai.HarmBlockLowAndAbove,
			},
		}
	} else {
		model.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockNone,
			},
		}
	}
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		return err
	}
	respJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}

	// Unmarshal the JSON response into the struct
	var genResp structs.GenResponse
	err = json.Unmarshal(respJSON, &genResp)
	if err != nil {
		return err
	}

	if len(genResp.Candidates) > 0 && len(genResp.Candidates[0].Content.Parts) > 0 {
		generatedText := genResp.Candidates[0].Content.Parts[0]
		// the generatedText is the command to be executed, so we need to run it
		fmt.Println("Running the command: ", generatedText)
		RunCommand(generatedText)
		return nil
	} else {
		fmt.Println("No generated text found")
		return nil
	}
}

func GetGeminiGeneralResponse(prompt string, safeOn bool) (string, error) {
	godotenv.Load()
	ctx := context.Background()

	geminiKey, err := keyring.Get("genie", "gemini_api_key")
	if err != nil {
		fmt.Println("Gemini API key not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()
	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	if safeOn {
		model.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: genai.HarmBlockLowAndAbove,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: genai.HarmBlockLowAndAbove,
			},
		}
	} else {
		model.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockNone,
			},
		}
	}
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		return "", err
	}
	respJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", err
	}

	// Unmarshal the JSON response into the struct
	var genResp structs.GenResponse
	err = json.Unmarshal(respJSON, &genResp)
	if err != nil {
		return "", err
	}

	if len(genResp.Candidates) > 0 && len(genResp.Candidates[0].Content.Parts) > 0 {
		generatedText := genResp.Candidates[0].Content.Parts[0]
		return generatedText, nil
	} else {
		return "No generated text found.", nil
	}
}
