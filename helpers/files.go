package helpers

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/harshalranjhani/genie/structs"
)

func GetCurrentDirectoriesAndFiles(root string) (structs.Directory, error) {
	rootDir := structs.Directory{Name: root}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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
