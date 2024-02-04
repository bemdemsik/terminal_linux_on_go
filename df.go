package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
)

var enterFlags []string
var commandFlags []string
var value []string
var totalSize float64
var free float64
var used float64
var usedIn float64
var stat syscall.Statfs_t

func main() {
	df(os.Args[1:])
}
func df(arguments []string) {
	commandFlags = append(commandFlags, "a", "H", "k", "h")
	if !fill(arguments) {
		return
	}

	if contains(enterFlags, "h") {
		if len(enterFlags) > 1 || len(value) > 0 {
			fmt.Println("Invalid key -h")
		} else {
			fmt.Println("-a - отобразить все файловые системы, в том числе виртуальные, псевдо")
			fmt.Println("-H - размер в мегабайтах")
			fmt.Println("-k - размер в килобайтах")
			fmt.Println("-h - помощь")
		}
		return
	}
	if contains(enterFlags, "H") && contains(enterFlags, "k") {
		fmt.Println("Invalid key -H and -k")
		return
	}
	linesAll := []string{""}
	if contains(enterFlags, "a") {
		files, err := ioutil.ReadFile("/proc/mounts")
		if err != nil {
			fmt.Println("Ошибка при чтении /proc/mounts")
			return
		}
		linesAll = strings.Split(string(files), "\n")
	} else {
		err := syscall.Statfs("/", &stat)
		if err != nil {
			fmt.Println("Ошибка при получении информации о файловой системе:", err)
			return
		}
		getSize()
		fmt.Printf("Размер файловой системы: %.3f\n", totalSize)
		fmt.Printf("Использовано: %.3f\n", used)
		fmt.Printf("Доступно: %.3f\n", free)
		return
	}

	fmt.Printf("%-30s%-20s%-20s%-20s%-20s%-20s\n", "Файловая система", "1К-блоков", "Использовано",
		"Доступно", "Использовано%", "Смонтированно в")
	for _, item := range linesAll {
		fields := strings.Fields(item)
		if len(fields) < 6 {
			continue
		}
		device := fields[0]
		mountPoint := fields[1]
		syscall.Statfs(mountPoint, &stat)
		if getSize() == 0 {
			continue
		}
		fmt.Printf("%-30s%-20.3f%-20.3f%-20.3f%-20.3f%-20s\n", device, totalSize, used, free, usedIn, mountPoint)
	}
}
func getSize() int {
	totalSize = float64(stat.Blocks * uint64(stat.Bsize))
	free = float64(stat.Bfree * uint64(stat.Bsize))
	used = totalSize - free
	if totalSize == 0 {
		totalSize, free, used, usedIn = 0, 0, 0, 0
		return 0
	}
	usedIn = 100 * used / totalSize
	if contains(enterFlags, "H") {
		totalSize /= 1024 * 1024
		free /= 1024 * 1024
		used /= 1024 * 1024
	}
	if contains(enterFlags, "k") {
		totalSize /= 1024
		free /= 1024
		used /= 1024
	}
	return 1
}
func fill(arguments []string) bool {
	enterFlags = []string{}
	value = []string{}
	for i := 0; i < len(arguments); i++ {
		if arguments[i][0] == '-' {
			if !contains(commandFlags, string(arguments[i][1])) {
				fmt.Print("Invalid key " + arguments[i])
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
