package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	allNumerate := flag.Bool("b", false, "Enumerate all strings, with empty")
	dollarLineFlag := flag.Bool("E", false, "Show dollar in every end of line")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	allArgs := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: nl [options] file")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var args []string

	for _, arg := range allArgs {
		if strings.HasPrefix(arg, "-") {
			args = append(args, arg)
		}
	}

	file := getPathFromString(strings.Join(allArgs, ""))

	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("File " + file + " does not exist")
		return
	}

	lines, err := readLinesForNumbers(file)
	if err != nil {
		fmt.Printf("Error reading file")
		return
	}

	index := 1

	for _, line := range lines {
		if *allNumerate {
			fmt.Print(strconv.Itoa(index) + " ")
			index++
		} else {
			if strings.TrimSpace(line) != "" {
				fmt.Print(strconv.Itoa(index) + " ")
				index++
			}
		}

		fmt.Printf("%s", line)

		if *dollarLineFlag {
			fmt.Print(" $\n")
		} else {
			fmt.Println()
		}
	}
}

func getPathFromString(input string) string {
	nlIndex := strings.Index(input, "nl ")

	if nlIndex == -1 {
		return ""
	}

	trimmed := input[nlIndex+len("nl "):]

	return trimmed
}

func readLinesForNumbers(filePath string) ([]string, error) {
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
