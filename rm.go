package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var enterFlags []string
var flags []string
var value []string
var enter string
var activeDirectory string

func main() {
	activeDirectory, _ = os.Getwd()
	rm(os.Args[1:])
}
func rm(arguments []string) {
	flags = append(flags, "d", "r", "v", "h", "f")
	if !fill(arguments) {
		return
	}
	if contains(enterFlags, "h") {
		if len(enterFlags) > 1 || len(value) > 0 {
			fmt.Println("Invalid key -h")
		} else {
			fmt.Println("-d - удалить все пустые каталоги")
			fmt.Println("-v - пояснять действия")
			fmt.Println("-f - игнорировать несуществующие файлы и аргументы, ни о чём не спрашивать")
			fmt.Println("-r - рекурсивное удаление")
			fmt.Println("-h - помощь")
		}
		return
	}
	if len(value) == 0 {
		fmt.Println("Usage: rm [flags] <directory> | <file>")
		return
	}
	deleteFiles()
}

func deleteFiles() {
	if contains(value, ".") {
		folder, _ := os.Open(activeDirectory)
		value, _ = folder.Readdirnames(0)
		folder.Close()
	}
	for _, file := range value {
		var err error
		dir := ""
		if _, err = os.Stat(activeDirectory + "/" + file); err == nil {
			dir = activeDirectory + "/" + file
		} else if _, err = os.Stat(file); err == nil {
			dir = file
		} else {
			fmt.Println("Данной директории или файла не существует", file)
			continue
		}
		if contains(enterFlags, "r") {
			err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if y, f := checkForce(path); y && f {
					err = os.RemoveAll(path)
				}
				if err != nil {
					return err
				}
				return nil
			})
		} else if contains(enterFlags, "d") {
			folder, _ := os.Open(dir)
			fileNames, _ := folder.Readdirnames(0)
			folder.Close()
			if len(fileNames) == 0 {
				if y, f := checkForce(dir); y && f {
					err = os.Remove(dir)
				}
			}
		} else {
			if y, f := checkForce(dir); y && f {
				err = os.Remove(dir)
			}
		}
		if err != nil {
			fmt.Println(err)
		} else if contains(enterFlags, "v") {
			fmt.Println("Удалено: ", dir)
		}
	}
}
func checkForce(path string) (bool, bool) {
	if contains(enterFlags, "f") {
		return false, false
	}
	simbol := ""
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Вы уверены что хотите удалить ", path, "? y/n: ")
		simbol, _ = reader.ReadString('\n')
		simbol = strings.Replace(simbol, "\n", "", -1)
		if simbol == "y" || simbol == "n" {
			break
		}
	}
	if simbol == "n" {
		return false, true
	}
	return true, true
}

func fill(arguments []string) bool {
	enterFlags = []string{}
	value = []string{}
	for i := 0; i < len(arguments); i++ {
		if arguments[i][0] == '-' {
			if !contains(flags, string(arguments[i][1])) {
				fmt.Println("Invalid key " + arguments[i])
				return false
			} else {
				enterFlags = append(enterFlags, string(arguments[i][1]))
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
