package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func checkError(err error) bool {
	if err != nil {
		return true
	}

	return false
}

func main() {
	allDirectoriesFlag := flag.Bool("p", false, "Create all directories in path")
	changeModeFlag := flag.Bool("m", false, "Change mod directories")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	args := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: mkdir [options] directory")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var onlyArgs []string

	var paths []string
	var mode string
	haveMode := false

	for i, arg := range args {
		if i > 0 && args[i-1] == "-m" {
			mode = args[i]
			haveMode = true
			continue
		}
		if strings.HasPrefix(arg, "-") {
			onlyArgs = append(onlyArgs, arg)
		} else {
			paths = append(paths, arg)
		}
	}

	if haveMode && strings.TrimSpace(mode) == "" {
		fmt.Println("Needs to input mode")
		return
	}

	if len(paths) != 1 {
		fmt.Println("Incorrect path")
		return
	}

	path := filepath.Clean("/" + paths[0] + "/")
	if !*allDirectoriesFlag {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			path = pathExist(filepath.Clean("/" + path + "/"))
		}
	} else {
		workDir, _ := os.Getwd()
		path = filepath.Join(workDir, path)
	}

	createDirectory(*changeModeFlag, *allDirectoriesFlag, path, mode)
}

func pathExist(path string) string {
	workDir, _ := os.Getwd()
	absPath := filepath.Clean(workDir + "/" + path)

	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return ""
	}

	return absPath

}

func createDirectory(withMode bool, withRecursive bool, value string, mode string) {
	permission := int64(os.ModePerm)
	var err error
	if withMode {
		permission, err = strconv.ParseInt(mode, 8, 32)
		if checkError(err) {
			fmt.Println("Incorrect mode")
		}
	}
	if withRecursive {
		if err := os.MkdirAll(value, os.FileMode(permission)); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := os.Mkdir(value, os.FileMode(permission)); err != nil {
			fmt.Println(err)
		}
	}
}

func isPath(path string) bool {
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
