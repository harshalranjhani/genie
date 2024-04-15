package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(musicCmd)
	musicCmd.PersistentFlags().String("d", "", "The duration for the song to be generated. (default: 8 seconds, maximum: 15 seconds)")
	musicCmd.PersistentFlags().Bool("logs", false, "Set this to true if you wish to see the logs of the music being generated.")
}

var musicCmd = &cobra.Command{
	Use:   "music",
	Short: "Generate music from text!",
	Long:  `This command generates music based on the provided text prompt.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		durationFlagValue, _ := cmd.Flags().GetString("d")
		if durationFlagValue == "" {
			durationFlagValue = "8"
		}
		durationNum, err := strconv.Atoi(durationFlagValue)
		if err != nil || durationNum > 15 || durationNum < 8 {
			fmt.Println("Invalid duration. Duration must be between 8 and 15 seconds.")
			os.Exit(1)
		}

		replicateApiKey, err := keyring.Get("genie", "replicate_api_key")
		if err != nil {
			fmt.Println("Replicate API key not found in keyring. Please run `genie init` to store the key.")
			os.Exit(1)
		}

		color.Cyan("Using Replicate API key")
		fmt.Println()
		url := "https://api.replicate.com/v1/predictions"
		contentType := "application/json"
		authHeader := "Bearer " + replicateApiKey

		showLogs, _ := cmd.Flags().GetBool("logs")

		jsonStr := fmt.Sprintf(`{
			"version": "671ac645ce5e552cc63a54a2bbff63fcf798043055d2dac5fc9e36a837eedcfb",
			"input": {
				"prompt": "%s",
				"model_version": "stereo-large",
				"output_format": "mp3",
				"normalization_strategy": "peak",
				"duration": %d
			}
		}`, prompt, durationNum)

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", contentType)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request to server:", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var response struct {
			Status string `json:"status"`
			Logs   string `json:"logs"`
			URLs   struct {
				Get string `json:"get"`
			} `json:"urls"`
			Output string `json:"output,omitempty"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Error parsing JSON")
			fmt.Println("Try checking your API key limits")
			return
		}

		fmt.Println("Waiting for the music to be generated...")
		getURL := response.URLs.Get
		for response.Status != "succeeded" {
			req, err := http.NewRequest("GET", getURL, nil)
			if err != nil {
				fmt.Println("Error creating polling request:", err)
				return
			}
			req.Header.Set("Authorization", authHeader)
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error checking status:", err)
				return
			}
			body, _ = ioutil.ReadAll(resp.Body)
			json.Unmarshal(body, &response)
			resp.Body.Close()

			fmt.Println("Raw json", string(body))
			fmt.Println("Current status:", response.Status)
			if len(response.Logs) > 0 && showLogs {
				fmt.Println(response.Logs)
			}

			if response.Status == "succeeded" {
				break
			}
			time.Sleep(5 * time.Second)
		}

		fmt.Println("Music generated successfully! Access it here:", response.Output)
	},
}
