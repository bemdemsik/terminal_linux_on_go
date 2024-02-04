package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	forceFlag := flag.Bool("f", false, "If file is exists, he is overwriting")
	recursiveFlag := flag.Bool("r", false, "Copying recursively directories")
	notForceFlag := flag.Bool("n", false, "If file is exists, he is not overwriting")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	inputArray := os.Args[1:]

	if *helpFlag {
		fmt.Println("The \"cp\" command is used to copy files and directories.")

		fmt.Println("Flags:")
		fmt.Println("-f    If set, when the destination file already exists, it will be overwritten.")
		fmt.Println("-r    If set, directories are copied recursively.")
		fmt.Println("-n    If set, existing files will not be overwritten.")

		fmt.Println("Example usage:")
		fmt.Println("cp -f file1.txt file2.txt   (Copy file1.txt to file2.txt, overwriting file2.txt if it exists)")
		return
	}

	var args []string

	var paths []string
	for _, arg := range inputArray {
		if arg == "cp" || arg == "./cp" {
			continue
		}

		if strings.HasPrefix(arg, "-") {
			args = append(args, arg)
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) != 2 {
		fmt.Println("Need a destination path and a departure path")
		return
	}

	if *forceFlag && *notForceFlag {
		fmt.Println("Arguments are incongruous")
		return
	}

	departurePath := paths[0]

	if *recursiveFlag && !isDirectory(departurePath) {
		fmt.Println("With recursion need the directory, not file")
		return
	}

	if !*recursiveFlag && isDirectory(departurePath) {
		fmt.Println("Without recursion need the file, not directory")
		return
	}

	if !*recursiveFlag && !isFile(departurePath) {
		fmt.Println("Need the input file")
		return
	}

	destinationPath := paths[1]

	if destinationPath == "." {
		destinationPath, _ = os.Getwd()
	}

	if *recursiveFlag && !isDirectory(destinationPath) {
		fmt.Println("With recursion, destination path should be a directory")
		return
	}

	if *recursiveFlag {
		if copyDirectory(departurePath, destinationPath, *forceFlag, *notForceFlag) != nil {
			fmt.Println("Error after recursive copying files and directories")
			return
		}
	} else {
		if copyFile(departurePath, destinationPath, *forceFlag, *notForceFlag) != nil {
			fmt.Println("Error after copying file")
			return
		}
	}
}

func copyDirectory(sourceDir, destDir string, isForce bool, isNotForce bool) error {
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	destSubDir := filepath.Join(destDir, filepath.Base(sourceDir))
	err = os.Mkdir(destSubDir, os.ModePerm)
	if err != nil {
		return err
	}

	dirEntries, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, entry := range dirEntries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		destinationPath := filepath.Join(destSubDir, entry.Name())
		if entry.IsDir() {
			err = copyDirectory(sourcePath, destinationPath, isForce, isNotForce)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(sourcePath, destinationPath, isForce, isNotForce)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(sourceFile, destFile string, isForce bool, isNotForce bool) error {
	_, err := os.Stat(destFile)
	if isNotForce && !os.IsNotExist(err) && !isForce {
		fmt.Println("File " + destFile + " already exist")
		return nil
	}

	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	destFileInfo, err := os.Stat(destFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if destFileInfo != nil && destFileInfo.IsDir() {
		destFile = filepath.Join(destFile, filepath.Base(filepath.Clean(sourceFile)))
	}

	destination, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	if !fileInfo.Mode().IsDir() {
		return false
	}

	return true
}

func isFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	if fileInfo.IsDir() {
		return false
	}

	return true
}
