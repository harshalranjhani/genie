package assets

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var TextFilePath string = "assets/genie-ascii.txt"

// print the contents from a text file to the console
func PrintTextFileContents(textFilePath string) {
	file, err := os.Open(textFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		fmt.Println(err)
		return
	}

	color.Cyan(string(bs))
}
