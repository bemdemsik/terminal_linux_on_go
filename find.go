package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var enterFlags []string
var commandFlags []string
var value []string
var activeDirectory string
var searchDirectories []string
var dirOrFile string
var valuePercent []string

func main() {
	activeDirectory, _ = os.Getwd()
	searchDirectories = []string{}
	find(os.Args[1:])
}

func find(arguments []string) {
	commandFlags = append(commandFlags, "v", "p", "f", "d", "h")
	if !fill(arguments) {
		return
	}
	if contains(enterFlags, "h") {
		if len(enterFlags) > 1 || len(value) > 0 {
			fmt.Println("Invalid key -h")
		} else {
			fmt.Println("-v - показать версию утилиты find")
			fmt.Println("-p - выводить полные имена файлов")
			fmt.Println("-f - искать только файлы")
			fmt.Println("-d - поиск папки в Linux")
			fmt.Println("-h - помощь")
		}
		return
	}
	if contains(enterFlags, "v") {
		if len(enterFlags) > 1 || len(value) > 0 {
			fmt.Println("Invalid key -v")
		} else {
			fmt.Println("find version 2.9.20")
		}
		return
	}
	if contains(enterFlags, "d") && contains(enterFlags, "f") {
		fmt.Println("Invalid flags -d and -f")
		return
	}
	if len(value) > 0 {
		for _, item := range value {
			searchDirectories = append(searchDirectories, item)
		}
	} else {
		searchDirectories = append(searchDirectories, activeDirectory)
	}

	for _, dir := range searchDirectories {
		index := 0
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			dirOrFile = ""
			if index > 0 {
				if len(valuePercent) != 0 {
					for _, file := range valuePercent {
						if strings.Contains(file, "*") {
							if !info.IsDir() && strings.Split(info.Name(), ".")[1] == strings.Split(file, ".")[1] {
								dirOrFile = path
							}
						} else if file == info.Name() {
							dirOrFile = path
						}
					}
				} else {
					if contains(enterFlags, "d") {
						if info.IsDir() {
							dirOrFile = path
						}
					} else if contains(enterFlags, "f") {
						if !info.IsDir() {
							dirOrFile = path
						}
					} else {
						dirOrFile = path
					}
				}

				if dirOrFile != "" {
					if contains(enterFlags, "p") {
						fmt.Println(dirOrFile)
					} else {
						fmt.Println(strings.Split(dirOrFile, "/")[len(strings.Split(dirOrFile, "/"))-1])
					}
				}
			}
			index += 1
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func fill(arguments []string) bool {
	enterFlags = []string{}
	value = []string{}
	for i := 0; i < len(arguments); i++ {
		if arguments[i][0] == '-' {
			if !contains(commandFlags, string(arguments[i][1])) {
				fmt.Println("Invalid key " + arguments[i])
				return false
			} else {
				enterFlags = append(enterFlags, string(arguments[i][1]))
			}
		} else if arguments[i][0] == '%' {
			valuePercent = append(valuePercent, strings.Replace(arguments[i], "%", "", -1))
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
