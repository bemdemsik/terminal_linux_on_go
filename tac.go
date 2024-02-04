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

	if *helpFlag {
		fmt.Println("Usage: tac [options] [file ...]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var onlyArgs []string
	var files []string

	for _, arg := range args {
		if arg == "./tac" || arg == "tac" {
			continue
		}
		if strings.HasPrefix(arg, "-") {
			onlyArgs = append(onlyArgs, arg)
		} else {
			files = append(files, arg)
		}
	}

	if *enumerateFlag && *notEmptyEnumerateFlag {
		fmt.Println("Invalid using arguments syntax")
		return
	}

	for _, file := range files {
		path := filepath.Clean("/" + file + "/")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = isExistingPath(filepath.Clean("/" + file + "/"))
			if path == "" {
				fmt.Println("File " + file + " does not exist")
				continue
			}
		}

		lines, err := readLinesReverse(path)
		if err != nil {
			fmt.Printf("Error reading file")
			return
		}

		lines = reverseArray(lines)
		index := 1

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

			fmt.Printf("%s", line)

			if *dollarLineFlag {
				fmt.Print(" $")
			}

			fmt.Print("\n")
		}
	}
}

func isExistingPath(path string) string {
	workDir, _ := os.Getwd()
	absPath := filepath.Clean(workDir + "/" + path)

	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return ""
	}

	return absPath

}

func reverseArray(files []string) []string {
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}

	return files
}

func readLinesReverse(filePath string) ([]string, error) {
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
