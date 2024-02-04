package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	clearScreen()
}

func clearScreen() {
	cmd := exec.Command(clearCommand())
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing clear")
	}
}

func clearCommand() string {
	switch runtime.GOOS {
	case "windows":
		return "cmd"
	default:
		return "clear"
	}
}
