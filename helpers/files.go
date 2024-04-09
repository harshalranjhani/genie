package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/harshalranjhani/genie/structs"
	"github.com/zalando/go-keyring"
)

func GetCurrentDirectoriesAndFiles(root string) (structs.Directory, error) {
	rootDir := structs.Directory{Name: root}
	ignoreListPath, err := keyring.Get("genie", "ignore_list_path")
	if err != nil {
		return structs.Directory{}, fmt.Errorf("error getting ignore list path: %w", err)
	}
	fmt.Println("Ignore List Path: ", ignoreListPath)
	ignorePatterns, err := readIgnorePatterns(ignoreListPath)
	if err != nil {
		return structs.Directory{}, fmt.Errorf("error reading ignore patterns: %w", err)
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if shouldIgnore(path, ignorePatterns) {
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

func readIgnorePatterns(filename string) ([]string, error) {
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

func shouldIgnore(path string, patterns []string) bool {
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
