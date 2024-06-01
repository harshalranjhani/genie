package helpers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/structs"
	"github.com/zalando/go-keyring"
)

func GetCurrentDirectoriesAndFiles(root string) (structs.Directory, error) {
	rootDir := structs.Directory{Name: root}
	ignoreListPath, err := keyring.Get("genie", "ignore_list_path")
	if err != nil {
		return structs.Directory{}, fmt.Errorf("error getting ignore list path: %w", err)
	}
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("Ignore List Path: ", ignoreListPath)
	ignorePatterns, err := ReadIgnorePatterns(ignoreListPath)
	if err != nil {
		return structs.Directory{}, fmt.Errorf("error reading ignore patterns: %w", err)
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ShouldIgnore(path, ignorePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		relativePath, _ := filepath.Rel(root, path)
		if info.IsDir() {
			if relativePath == "." {
				return nil
			}
			dir := structs.Directory{Name: relativePath}
			rootDir.Children = append(rootDir.Children, dir)
		} else {
			file := structs.File{Name: relativePath, Size: info.Size()}
			rootDir.Files = append(rootDir.Files, file)
		}
		return nil
	})

	if err != nil {
		return structs.Directory{}, err
	}

	return rootDir, nil
}

func ReadIgnorePatterns(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		patterns = append(patterns, scanner.Text())
	}
	return patterns, scanner.Err()
}

func ShouldIgnore(path string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			log.Printf("Invalid pattern %q: %v", pattern, err)
			continue
		}
		if matched {
			return true
		}
	}
	// Additional check to ignore hidden files/folders.
	if strings.HasPrefix(filepath.Base(path), ".") {
		return true
	}
	return false
}

func PrintData(sb *strings.Builder, root structs.Directory, level int) {
	indent := strings.Repeat("  ", level)
	sb.WriteString(fmt.Sprintf("%s[%s]\n", indent, root.Name))

	for _, file := range root.Files {
		sb.WriteString(fmt.Sprintf("%s- %s\n", indent+"  ", file.Name))
	}

	for _, child := range root.Children {
		PrintData(sb, child, level+1)
	}
}

func RunCommand(command string) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateMarkdown(headings []structs.Heading, fileName string) {

	// if headings is empty, return an error
	if len(headings) == 0 {
		color.Red("No genie headings found to generate markdown file.")
		return
	}

	outputPath := fileName + ".md"
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating markdown file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, heading := range headings {
		link := fmt.Sprintf("[%s:%d](%s#L%d)", filepath.Base(heading.FilePath), heading.LineNum, heading.FilePath, heading.LineNum)
		_, err := writer.WriteString(fmt.Sprintf("## %s: %s\n", link, heading.Content))
		if err != nil {
			fmt.Printf("Error writing to markdown file: %v\n", err)
			return
		}
		for _, subheading := range heading.Subheadings {
			subLink := fmt.Sprintf("[%s:%d](%s#L%d)", filepath.Base(heading.FilePath), subheading.LineNum, heading.FilePath, subheading.LineNum)
			_, err := writer.WriteString(fmt.Sprintf("  - %s: %s\n", subLink, subheading.Content))
			if err != nil {
				fmt.Printf("Error writing to markdown file: %v\n", err)
				return
			}
		}
	}
	writer.Flush()
	fmt.Println("Markdown file generated successfully!")
}

type MailObj struct {
	Email    string            `json:"email"`
	Headings []structs.Heading `json:"headings"`
}

type MailRequest struct {
	MailObj MailObj `json:"mailObj"`
}

func SendMarkdownFileToEmail(email string, headings []structs.Heading) error {

	// if headings is empty, return an error
	if len(headings) == 0 {
		color.Red("No genie headings found to send in the email.")
		return fmt.Errorf("no headings found to send in the email")
	}

	fmt.Println("Sending email...")
	// Prepare the request payload
	mailRequest := MailRequest{
		MailObj: MailObj{
			Email:    email,
			Headings: headings,
		},
	}

	// Marshal the payload to JSON
	jsonData, err := json.Marshal(mailRequest)
	if err != nil {
		return fmt.Errorf("error marshaling mail request: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", "https://api.harshalranjhani.in/mail/genie-summary", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from server: %s", body)
	}

	fmt.Println("Email sent successfully.")
	return nil
}
