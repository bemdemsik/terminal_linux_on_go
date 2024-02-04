package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	numberCommand, er := strconv.ParseUint(os.Args[1], 10, 64)
	if er != nil {
		fmt.Println("Usage: ex [номер команды]")
		return
	}
	command := ""
	filePath := "/home/user/.bash_history"
	file, _ := os.ReadFile(filePath)
	fileContent := strings.Split(string(file), "\n")
	for i, item := range fileContent {
		if numberCommand == uint64(i+1) {
			command = item
			break
		}
	}
	if command == "" {
		fmt.Println("Команда с номером", numberCommand, "не найдена")
		return
	}
	fmt.Println(command)
	cmd := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
}
