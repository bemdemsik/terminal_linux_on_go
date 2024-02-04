package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	filePath := "/home/user/.bash_history"
	file, _ := os.ReadFile(filePath)
	fileContent := strings.Split(string(file), "\n")
	command := fileContent[len(fileContent)-2]
	fmt.Println(command)
	if len(strings.Split(command, " ")) == 1 {
		cmd := exec.Command(strings.Split(command, " ")[0])
		output, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(output))
	} else {
		cmd := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)
		output, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(output))
	}
}
