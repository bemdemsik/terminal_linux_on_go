package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	currentTimeFlag := flag.Bool("a", false, "Change access time to current")
	overwritingFlag := flag.Bool("o", false, "Overwrite file, if is exists")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	args := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: command touch [options] [file(s)]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var onlyArgs []string
	var files []string

	for _, arg := range args {
		if arg == "touch" || arg == "./touch" {
			continue
		}
		if strings.HasPrefix(arg, "-") {
			onlyArgs = append(onlyArgs, arg)
		} else {
			files = append(files, arg)
		}
	}

	for _, file := range files {

		fileName := filepath.Base(file)
		path := existingPath(filepath.Dir(filepath.Clean(file)))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = filepath.Dir(filepath.Clean(file))
		}

		path = filepath.Join(path, fileName)

		employmentFile(path, *currentTimeFlag, *overwritingFlag)
	}
}

func existingPath(path string) string {
	workDir, _ := os.Getwd()
	absPath := filepath.Clean(workDir + "/" + path)

	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return ""
	}

	return absPath

}

func employmentFile(fileName string, withChangeTime bool, withOverwrite bool) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if !createFile(fileName) {
			return
		}
	} else {
		if withOverwrite {
			if !deleteFile(fileName) {
				return
			}
			if !createFile(fileName) {
				return
			}
		} else {
			fmt.Println("File " + fileName + " already exists")
		}
	}

	if withChangeTime {
		err := os.Chtimes(fileName, time.Now(), time.Now())
		if err != nil {
			fmt.Println(fmt.Println("Error updating access time file"))
			return
		}
	}
}

func createFile(fileName string) bool {
	file, createErr := os.Create(fileName)
	if createErr != nil {
		fmt.Println("Error creating file " + fileName)
		return false
	}
	defer file.Close()

	return true
}

func deleteFile(file string) bool {
	if _, err := os.Stat(file); err == nil {
		if err := os.Remove(file); err != nil {
			fmt.Println("Error removing file")
			return false
		}
	}

	return true
}

func isCorrectPath(path string) bool {
	pathArray := strings.Split(path, "/")

	var cleanParts []string
	for _, part := range pathArray {
		if strings.TrimSpace(part) == "" {
			continue
		}
		cleanParts = append(cleanParts, part)
	}

	if len(cleanParts) > 1 {
		return true
	}

	return false
}
