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
	outFileNameFlag := flag.Bool("v", false, "Out name of the file before text")
	specifyNumbersLineFlag := flag.Bool("n", false, "Allows to specify the number of lines")
	bytesFlag := flag.Bool("c", false, "Reading file in bytes")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	args := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: head [options] [file]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var files []string
	var onlyArgs []string
	var numberLines string
	var bytes string
	hasBytes := false
	hasLines := false

	for i, arg := range args {
		if i > 0 && args[i-1] == "-n" {
			hasLines = true
			numberLines = args[i]
			continue
		}

		if i > 0 && args[i-1] == "-c" {
			hasBytes = true
			bytes = arg
			continue
		}

		if strings.HasPrefix(arg, "-") {
			onlyArgs = append(onlyArgs, arg)
		} else {
			files = append(files, arg)
		}
	}

	if hasLines && hasBytes {
		fmt.Println("Arguments are incompatible")
		return
	}

	var linesOfFile int
	var err error
	if hasLines {
		linesOfFile, err = strconv.Atoi(numberLines)
		if err != nil {
			fmt.Println("Not correct lines")
			return
		}
	}

	var bytesNumber int

	if hasBytes {
		bytesNumber, err = strconv.Atoi(bytes)
		if err != nil {
			fmt.Println("Not correct bytes")
			return
		}

	}

	for _, file := range files {
		path := filepath.Clean("/" + file + "/")
		if _, err := os.Stat(file); os.IsNotExist(err) {
			path = pathIsExist(filepath.Clean("/" + file + "/"))
			if path == "" {
				fmt.Println("File " + file + " does not exist")
				continue
			}
		}

		lines, err := readLines(path)
		if err != nil {
			fmt.Printf("Error reading file")
			return
		}

		if *outFileNameFlag {
			fmt.Println(file)
		}

		bytesCount := 0

		for i, line := range lines {
			if *specifyNumbersLineFlag {
				if i == linesOfFile {
					break
				}
			} else if *bytesFlag {
				bytesCount += len(line) + 1
				if bytesCount > bytesNumber {
					break
				}
			} else {
				if i == 10 {
					break
				}
			}

			fmt.Printf("%s\n", line)
		}
	}
}

func readLines(filePath string) ([]string, error) {
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

func pathIsExist(path string) string {
	workDir, _ := os.Getwd()
	absPath := filepath.Clean(workDir + "/" + path)

	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return ""
	}

	return absPath

}
