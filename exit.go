package main

import (
	"fmt"
	"syscall"
)

func main() {

	if err := syscall.Kill(syscall.Getppid(), syscall.SIGKILL); err != nil {
		fmt.Println("Failed to kill process:", err)
		return
	}
}
