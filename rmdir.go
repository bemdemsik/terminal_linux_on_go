package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	removeAllFlag := flag.Bool("R", false, "Remove all directories, with not empty")
	requestFlag := flag.Bool("Q", false, "Confirmation on deleting")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	allArgs := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: rmdir [options] directory")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var args []string

	var paths []string

	for _, arg := range allArgs {

		if strings.HasPrefix(arg, "-") {
			args = append(args, arg)
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) != 1 {
		fmt.Println("Should be only one path")
		return
	}

	path := filepath.Clean("/" + paths[0] + "/")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = alreadyPath(filepath.Clean("/" + path + "/"))
		if path == "" {
			fmt.Println("Path " + path + " does not exist")
			return
		}
	}

	if *requestFlag {
		fmt.Println("Do you really want to delete this directory? (y/n)")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		if input == "y\n" {
			if *removeAllFlag {
				removeAllDirectories(path)
			} else {
				removeOnlyEmptyDirectories(path)
			}
		}

	} else {
		if *removeAllFlag {
			removeAllDirectories(path)
		} else {
			removeOnlyEmptyDirectories(path)
		}
	}
}

func removeAllDirectories(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println("Error removing directory")
	}
}

func alreadyPath(path string) string {
	workDir, _ := os.Getwd()
	absPath := filepath.Clean(workDir + "/" + path)

	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return ""
	}

	return absPath

}

func removeOnlyEmptyDirectories(path string) {
	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dir, err := os.Open(currentPath)
			if err != nil {
				return err
			}
			defer dir.Close()

			_, err = dir.Readdir(1)
			if err != nil {
				err = os.Remove(currentPath)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

func isRealPath(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
