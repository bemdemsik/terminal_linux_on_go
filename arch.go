package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	names := os.Args

	if len(names) == 1 {
		hostname, err := os.Hostname()
		if err == nil {
			fmt.Println("имя компьютера: ", hostname)
		}
		fmt.Println("архитектура: ", runtime.GOARCH)
	} else {
		if names[1] == "-h" {
			fmt.Println("Использование: arch [ПАРАМЕТР]\nПечатает машинную архитектуру\n-h - показать справку\n")
		} else {
			fmt.Println("Ошибка использования ")
		}
	}
}
