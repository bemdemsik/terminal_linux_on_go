package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var enterFlags []string
var commandFlags []string
var value []string
var activeDirectory string
var signal syscall.Signal = syscall.SIGKILL

func main() {
	kill(os.Args[1:])
}
func kill(arguments []string) {
	commandFlags = append(commandFlags, "SIGINT", "SIGTERM", "SIGQUIT", "SIGHUP", "SIGKILL", "h")
	if !fill(arguments) {
		return
	}
	if contains(enterFlags, "h") {
		if len(enterFlags) > 1 || len(value) > 0 {
			fmt.Println("Invalid key -h")
		} else {
			fmt.Println("-SIGINT - процесс правильно завершает все свои действия и возвращает управление")
			fmt.Println("-SIGQUIT - программа может выполнить корректное завершение или проигнорировать сигнал. " +
				"В отличие от предыдущего, она генерирует дамп памяти")
			fmt.Println("-SIGHUP - сообщает процессу, что соединение с управляющим терминалом разорвано")
			fmt.Println("-SIGTERM - немедленно завершает процесс, но обрабатывается программой")
			fmt.Println("-SIGKILL - тоже немедленно завершает процесс, но, в отличие от предыдущего варианта, он обрабатывается ядром. " +
				"Поэтому ресурсы и дочерние процессы остаются запущенными")
			fmt.Println("-h - помощь")
		}
		return
	}
	if len(value) == 0 {
		fmt.Println("Usage: kill [flags] <pid>")
		return
	}
	if len(enterFlags) > 0 {
		getSignal()
	}
	pid, err := strconv.Atoi(value[0])
	if err != nil {
		fmt.Println("Invalid PID:", value[1])
		return
	}

	if err := syscall.Kill(pid, signal); err != nil {
		fmt.Println("Failed to kill process:", err)
		return
	}

	fmt.Printf("Process with PID %d killed\n", pid)
}
func getSignal() {
	switch enterFlags[0] {
	case "SIGQUIT":
		signal = syscall.SIGQUIT
	case "SIGTERM":
		signal = syscall.SIGTERM
	case "SIGINT":
		signal = syscall.SIGINT
	default:
		signal = syscall.SIGKILL
	}
}
func fill(arguments []string) bool {
	enterFlags = []string{}
	value = []string{}
	for i := 0; i < len(arguments); i++ {
		if arguments[i][0] == '-' {
			if !contains(commandFlags, strings.Replace(arguments[i], "-", "", -1)) {
				fmt.Println("Invalid key " + arguments[i])
				return false
			} else {
				enterFlags = append(enterFlags, strings.Replace(arguments[i], "-", "", -1))
			}
		} else {
			value = append(value, arguments[i])
		}
	}
	return true
}
func contains(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}
