package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"syscall"
)

var enterFlags []string
var commandFlags []string
var value []string
var activeDirectory string
var totalSize float64
var totalSizeDirectory float64
var specifiedSize int64
var stat syscall.Statfs_t

func main() {
	activeDirectory, _ = os.Getwd()
	du(os.Args[1:])
}
func du(arguments []string) {
	commandFlags = append(commandFlags, "t", "M", "k", "c", "h")
	if !fill(arguments) {
		return
	}

	if contains(enterFlags, "h") {
		if len(enterFlags) > 1 || len(value) > 0 {
			fmt.Println("Invalid key -h")
		} else {
			fmt.Println("-t - не учитывать файлы и папки с размером меньше указанного(в байтах)")
			fmt.Println("-M - размер в мегабайтах")
			fmt.Println("-k - размер в килобайтах")
			fmt.Println("-c - выводить в конце общий размер всех папок;")
			fmt.Println("-h - помощь")
		}
		return
	}
	if contains(enterFlags, "M") && contains(enterFlags, "k") {
		fmt.Println("Invalid key -M and -k")
		return
	}
	if contains(enterFlags, "t") {
		index := indexOf(arguments, "-t") + 1
		if index == -1 {
			fmt.Println("Invalid value key -t: ", arguments[index])
			return
		}
		if reflect.TypeOf(index).Kind() == reflect.Float64 {
			check, er := strconv.ParseInt(arguments[index], 10, 64)
			if er != nil {
				fmt.Println("Invalid value key -t: ", arguments[index])
				return
			} else {
				specifiedSize = check
				value = append(value[:index], value[index+1:]...)
			}
		}
	}
	if len(value) == 0 {
		fmt.Println("Usage: du [flags] <directory>")
		return
	}
	for _, item := range value {
		totalSizeDirectory = 0
		dirOrFile := ""
		if _, err := os.Stat(activeDirectory + "/" + item); err == nil {
			dirOrFile = activeDirectory + item
		} else if _, err = os.Stat(item); err == nil {
			dirOrFile = item
		} else {
			fmt.Println("Данной директории не существует", item)
			continue
		}
		err := filepath.Walk(dirOrFile, func(path string, info fs.FileInfo, err error) error {
			if dirOrFile != path {
				syscall.Statfs(path, &stat)
				if contains(enterFlags, "t") {
					if info.Size() > specifiedSize {
						printDirOrFile(info.Size(), path)
					}
				} else {
					printDirOrFile(info.Size(), path)
				}
			} else {
				totalSizeDirectory = float64(info.Size())
				if contains(enterFlags, "M") {
					totalSizeDirectory /= 1024 * 1024
				}
				if contains(enterFlags, "k") {
					totalSizeDirectory /= 1024
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
		if contains(enterFlags, "c") {
			fmt.Printf("Общий размер %s: %.3f\n", dirOrFile, totalSizeDirectory)
		}
	}
}
func printDirOrFile(size int64, path string) {
	totalSize = float64(size)
	changeSize()
	totalSizeDirectory += totalSize
	fmt.Printf("%-20.3f%s\n", totalSize, path)
}
func changeSize() {
	if contains(enterFlags, "M") {
		totalSize /= 1024 * 1024
	}
	if contains(enterFlags, "k") {
		totalSize /= 1024
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
func indexOf(list []string, element string) int {
	for i := 0; i < len(list); i++ {
		if list[i] == element {
			return i
		}
	}
	return -1
}
