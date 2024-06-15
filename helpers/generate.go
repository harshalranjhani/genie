package helpers

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
	"github.com/zalando/go-keyring"
)

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

func GenerateGPTImage(prompt string) (string, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = color.HiCyanString("Analyzing: ")
	s.Start()

	openAIKey, err := keyring.Get("genie", "openai_api_key")
	if err != nil {
		s.Stop()
		fmt.Println("OpenAI API key not found in keyring. Please run `genie init` to store the key.")
		return "", err
	}
	c := openai.NewClient(openAIKey)
	ctx := context.Background()

	// reqUrl := openai.ImageRequest{
	// 	Prompt:         prompt,
	// 	Size:           openai.CreateImageSize256x256,
	// 	ResponseFormat: openai.CreateImageResponseFormatURL,
	// 	N:              1,
	// }
	// respUrl, err := c.CreateImage(ctx, reqUrl)
	// if err != nil {
	// 	s.Stop()
	// 	fmt.Printf("Image creation error: %v\n", err)
	// 	return "", err
	// }
	// fmt.Println(respUrl.Data[0].URL)

	reqBase64 := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	respBase64, err := c.CreateImage(ctx, reqBase64)
	if err != nil {
		s.Stop()
		fmt.Printf("Image creation error: %v\n", err)
		return "", err
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		s.Stop()
		fmt.Printf("Base64 decode error: %v\n", err)
		return "", err
	}

	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		s.Stop()
		fmt.Printf("PNG decode error: %v\n", err)
		return "", err
	}

	prompt = strings.Replace(prompt, " ", "_", -1)
	filename := fmt.Sprintf("images/%s.png", prompt)

	file, err := os.Create(filename)
	if err != nil {
		s.Stop()
		fmt.Printf("File creation error: %v\n", err)
		return "", err
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		s.Stop()
		fmt.Printf("PNG encode error: %v\n", err)
		return "", err
	}

	s.Stop()

	return filename, nil
}
