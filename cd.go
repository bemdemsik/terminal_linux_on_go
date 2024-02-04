package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	showErrorFlag := flag.Bool("e", false, "Show error if directory does not exist")
	helpFlag := flag.Bool("h", false, "Help")
	flag.Parse()

	args := os.Args[1:]

	if *helpFlag {
		fmt.Println("cd [OPTIONS] [DIRECTORY]")
		fmt.Println("")
		fmt.Println("OPTIONS:")
		flag.PrintDefaults()
		return
	}

	parentDirectory, _ := os.Getwd()

	err := executeChangeDirectory(args, parentDirectory, *showErrorFlag)
	if err != nil {
		fmt.Println("Error to change directory")
	}
}

func executeChangeDirectory(args []string, parentDirectory string, showErrorFlag bool) error {
	if len(args) < 1 {
		if err := os.Chdir("/home/user"); err != nil {
			return err
		}
		os.Setenv("PWD", "/home/user")
		fmt.Println("/home/user")
		return nil
	}

	if len(args) > 2 {
		fmt.Println("Usage cd <directory>")
		return nil
	}

	directory := args[0]

	curDirectory := ""

	switch directory {
	case "-":
		if err := os.Chdir(parentDirectory); err != nil {
			return fmt.Errorf("Failed to change directory to previous directory")
		}
		fmt.Println(parentDirectory)
		os.Setenv("PWD", parentDirectory)
		return nil
	case "..":
		parentDirectory = filepath.Dir(parentDirectory)
		curDirectory = parentDirectory
	default:
		curDirectory = filepath.Clean("/" + directory + "/")
		if _, err := os.Stat(curDirectory); os.IsNotExist(err) {
			curDirectory = checkContinuePath(curDirectory)
			if curDirectory == "" {
				if showErrorFlag {
					return fmt.Errorf("Directory does not exist")
				} else {
					return nil
				}

			}
		}
		parentDirectory = curDirectory

	}

	if err := os.Chdir(parentDirectory); err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return fmt.Errorf("Failed to change directory")
	}

	os.Setenv("PWD", parentDirectory)
	fmt.Println(parentDirectory)
	return nil
}

func checkContinuePath(path string) string {
	workDir, _ := os.Getwd()
	absPath := filepath.Clean(workDir + "/" + path)

	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return ""
	}

	return absPath
}
