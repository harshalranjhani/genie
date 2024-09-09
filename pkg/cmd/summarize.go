package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/helpers"
	"github.com/harshalranjhani/genie/internal/structs"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func init() {
	rootCmd.AddCommand(summarizeCmd)
	summarizeCmd.PersistentFlags().String("email", "", "The email to send the markdown summary to.")
	summarizeCmd.PersistentFlags().Bool("support", false, "Lists down the supported file types for summarization.")
	summarizeCmd.PersistentFlags().String("filename", "summary", "The name of the markdown file to be generated.")
}

var headings []structs.Heading

var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Get a markdown summary of the current directory comments",
	Long:  "Whenever you start comments with genie:heading: or genie:subheading: in your files, you can use this command to get a markdown summary of the comments in the current directory.",
	Run: func(cmd *cobra.Command, args []string) {

		email, _ := cmd.Flags().GetString("email")
		supportFlag, _ := cmd.Flags().GetBool("support")
		fileName, _ := cmd.Flags().GetString("filename")

		if supportFlag {
			color.Yellow("Supported file types for summarization:")
			for key := range commentMarkers {
				color.Yellow(key)
			}
			return
		}

		// Get the cwd root
		root, err := os.Getwd()

		if err != nil {
			color.Red("Error getting current working directory: %v", err)
			return
		}

		c := color.New(color.BgHiBlue).Add(color.Underline)
		c.Printf("Generating markdown summary for directory: %s\n", root)

		ignoreListPath, err := keyring.Get("genie", "ignore_list_path")
		if err != nil {
			color.Red("Error getting ignore list path: %v", err)
			return
		}
		c = color.New(color.FgCyan).Add(color.Underline)
		c.Println("Ignore List Path: ", ignoreListPath)
		ignorePatterns, err := helpers.ReadIgnorePatterns(ignoreListPath)
		if err != nil {
			color.Red("Error reading ignore patterns: %v", err)
			return
		}

		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if helpers.ShouldIgnore(path, ignorePatterns) {
				// Debug: Print each ignored file
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			return processFileForSummarization(path, info, err)
		})
		if err != nil {
			color.Red("Error walking through the directory: %v", err)
			return
		}

		if email != "" {
			// Send the markdown to the email
			helpers.SendMarkdownFileToEmail(email, headings)
		} else {
			helpers.GenerateMarkdown(headings, fileName)
		}
	},
}

var commentMarkers = map[string]string{
	".go":    "//",
	".js":    "//",
	".ts":    "//",
	".jsx":   "//",
	".tsx":   "//",
	".py":    "#",
	".java":  "//",
	".c":     "//",
	".cpp":   "//",
	".cs":    "//",
	".rb":    "#",
	".rs":    "//",
	".swift": "//",
	".kt":    "//",
	".php":   "//",
	".sh":    "#",
	".pl":    "#",
	".r":     "#",
	".scala": "//",
	".hs":    "--",
}

func processFileForSummarization(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}

	extension := filepath.Ext(path)
	commentMarker, supported := commentMarkers[extension]
	if !supported {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1
	var currentHeading *structs.Heading
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Use the comment marker to identify headings and subheadings
		if strings.HasPrefix(line, fmt.Sprintf("%s genie:heading:", commentMarker)) {
			content := strings.TrimPrefix(line, fmt.Sprintf("%s genie:heading:", commentMarker))
			currentHeading = &structs.Heading{
				FilePath: path,
				LineNum:  lineNum,
				Content:  strings.TrimSpace(content),
			}
			headings = append(headings, *currentHeading)
		} else if strings.HasPrefix(line, fmt.Sprintf("%s genie:subheading:", commentMarker)) {
			if currentHeading == nil {
				color.Yellow("Found subheading without a heading in file: %s, line: %d", path, lineNum)
			} else {
				content := strings.TrimPrefix(line, fmt.Sprintf("%s genie:subheading:", commentMarker))
				subheading := structs.Subheading{
					LineNum: lineNum,
					Content: strings.TrimSpace(content),
				}
				// Add subheading to the last heading
				headings[len(headings)-1].Subheadings = append(headings[len(headings)-1].Subheadings, subheading)
			}
		}
		lineNum++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
