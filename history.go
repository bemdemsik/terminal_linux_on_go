package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {

	flag.Parse()
	current, _ := user.Current()

	filePath := filepath.Join(current.HomeDir, ".bash_history")
	file, _ := os.ReadFile(filePath)

	fileContent := strings.Split(string(file), "\n")

	for index, command := range fileContent {
		if command != "" {
			fmt.Printf("%5d  %s\n", index+1, command)
		}
	}
}
