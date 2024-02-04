package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	notEmptyEnumerateFlag := flag.Bool("b", false, "Numerate only not empty lines")
	dollarLineFlag := flag.Bool("E", false, "Show dollar in every end of line")
	enumerateFlag := flag.Bool("n", false, "Enumerate strings in all files")
	helpFlag := flag.Bool("h", false, "Help")
	flag.Parse()
	args := os.Args[1:]

	var onlyArgs []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			onlyArgs = append(onlyArgs, arg)
		}
	}

	if *helpFlag {
		fmt.Println("cat [OPTIONS] [FILE...]")
		fmt.Println("")
		fmt.Println("OPTIONS:")
		fmt.Println("  -b  Numerate only not empty lines")
		fmt.Println("      Example: cat -b file.txt")
		fmt.Println("")
		fmt.Println("  -E  Show dollar in every end of line")
		fmt.Println("      Example: cat -E file.txt")
		fmt.Println("")
		fmt.Println("  -n  Enumerate strings in all files")
		fmt.Println("      Example: cat -n file1.txt file2.txt")
		return
	}

	if *enumerateFlag && *notEmptyEnumerateFlag {
		fmt.Println("Invalid using arguments syntax")
		return
	}

	inputArray := args

	var files []string

	for i := 0; i < len(inputArray); i++ {
		if inputArray[i] == "cat" || inputArray[i] == "./cat" {
			continue
		}
		if strings.HasPrefix(inputArray[i], "-") {
			continue
		}

		files = append(files, inputArray[i])
	}

	index := 1

	for _, file := range files {
		path := filepath.Clean("/" + file + "/")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = isExistPath(filepath.Clean("/" + file + "/"))
			if path == "" {
				fmt.Println("File " + file + " does not exist")
				continue
			}
		}

		lines, err := readAllLines(path)
		if err != nil {
			fmt.Println("Error reading file " + file)
			return
		}

		for _, line := range lines {
			if *enumerateFlag {
				fmt.Print(strconv.Itoa(index) + " ")
				index++
			}

			if *notEmptyEnumerateFlag {
				if strings.TrimSpace(line) != "" {
					fmt.Print(strconv.Itoa(index) + " ")
					index++
				}
			}

			fmt.Print(line)

			if *dollarLineFlag {
				fmt.Print(" $")
			}

			fmt.Print("\n")
		}
	}
}

func isExistPath(path string) string {
	workDir, _ := os.Getwd()
	absPath := filepath.Clean(workDir + "/" + path)

	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return ""
	}

	return absPath

}

func readAllLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
