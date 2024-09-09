package llm

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/google/generative-ai-go/genai"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/internal/structs"
	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"
	"google.golang.org/api/option"
)

func GetGeminiCmdResponse(prompt string, safeOn bool) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing: ")
	s.Start()

	godotenv.Load()
	ctx := context.Background()

	geminiKey, err := keyring.Get("genie", "gemini_api_key")
	if err != nil {
		s.Stop()
		fmt.Println("Gemini API key not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	if err != nil {
		s.Stop()
		return err
	}
	defer client.Close()

	// For text-only input, use the gemini-1.5-pro model
	model := client.GenerativeModel("gemini-1.5-pro")
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
		s.Stop()
		return err
	}
	respJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		s.Stop()
		return err
	}

	// Unmarshal the JSON response into the struct
	var genResp structs.GenResponse
	err = json.Unmarshal(respJSON, &genResp)
	if err != nil {
		s.Stop()
		return err
	}

	if len(genResp.Candidates) > 0 && len(genResp.Candidates[0].Content.Parts) > 0 {
		generatedText := genResp.Candidates[0].Content.Parts[0]
		// The generatedText is the command to be executed, so we need to run it
		s.Stop()
		fmt.Println("Running the command: ", generatedText)
		helpers.RunCommand(generatedText)
		return nil
	} else {
		s.Stop()
		fmt.Println("No generated text found")
		return nil
	}
}

func GetGeminiGeneralResponse(prompt string, safeOn bool, includeDir bool) (string, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing: ")
	s.Start()

	godotenv.Load()
	ctx := context.Background()

	geminiKey, err := keyring.Get("genie", "gemini_api_key")
	if err != nil {
		s.Stop()
		fmt.Println("Gemini API key not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	if err != nil {
		s.Stop()
		return "", err
	}
	defer client.Close()

	// For text-only input, use the gemini-1.5-pro model
	model := client.GenerativeModel("gemini-1.5-pro")
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
		s.Stop()
		return "", err
	}
	respJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		s.Stop()
		return "", err
	}

	var genResp structs.GenResponse
	err = json.Unmarshal(respJSON, &genResp)
	if err != nil {
		s.Stop()
		return "", err
	}

	if len(genResp.Candidates) > 0 && len(genResp.Candidates[0].Content.Parts) > 0 {
		generatedText := genResp.Candidates[0].Content.Parts[0]
		s.Stop()
		return generatedText, nil
	} else {
		s.Stop()
		return "No generated text found.", nil
	}
}

//go:embed scripts/generate.py
var generatePy []byte

func GenerateGeminiImage(prompt string) (string, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing: ")
	s.Start()

	// get ssid from keyring
	ssid, keyRingError := keyring.Get("genie", "ssid")
	if keyRingError != nil {
		s.Stop()
		fmt.Println("SSID not found in keyring. Please run `genie init` to store the key.")
		os.Exit(1)
	}

	tmpFile, err := os.CreateTemp("", "generate-*.py")
	if err != nil {
		s.Stop()
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(generatePy); err != nil {
		s.Stop()
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		s.Stop()
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	cmd := exec.Command("python", tmpFile.Name(), prompt, ssid)

	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command
	err = cmd.Run()
	if err != nil {
		s.Stop()
		return "", fmt.Errorf("failed to execute python script: %w", err)
	}

	filename := out.String()
	s.Stop()

	return filename, nil
}