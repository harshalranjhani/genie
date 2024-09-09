package cmd

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/spf13/cobra"
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape [url]",
	Short: "Scrape data from a URL, supports pagination!",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		element, _ := cmd.Flags().GetString("element")
		output, _ := cmd.Flags().GetString("output")
		pagination, _ := cmd.Flags().GetString("pagination")
		limit, _ := cmd.Flags().GetInt("limit")

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Prefix = color.HiCyanString("Scraping: ")
		s.Start()

		data, err := helpers.ScrapeURL(url, element, pagination, limit)
		if err != nil {
			s.Stop()
			fmt.Printf("Error scraping URL: %v\n", err)
			return
		}

		s.Stop()

		if output != "" {
			err := ioutil.WriteFile(output, []byte(data), 0644)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
			fmt.Printf("Data saved to %s\n", output)
		} else {
			fmt.Println("Scraped Data:")
			fmt.Println(data)
		}
	},
}

func init() {
	scrapeCmd.Flags().StringP("element", "e", "", "HTML element to extract (e.g., h1, p, a)")
	scrapeCmd.Flags().StringP("output", "o", "", "File to save scraped data")
	scrapeCmd.Flags().StringP("pagination", "p", "", "CSS selector for pagination links")
	scrapeCmd.Flags().IntP("limit", "l", 0, "Maximum number of pages to scrape")
	rootCmd.AddCommand(scrapeCmd)
}
