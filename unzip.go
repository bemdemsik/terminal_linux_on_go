package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	listFilesFlag := flag.Bool("l", false, "Show all files in archive")
	testErrorFlag := flag.Bool("t", false, "Checking the archive for errors")
	commentFlag := flag.Bool("z", false, "Show comment after uncompressed")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	inputArray := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: command unzip [options] [archive] [path]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var args []string

	var paths []string
	for _, arg := range inputArray {
		if arg == "unzip" || arg == "./unzip" {
			continue
		}

		if strings.HasPrefix(arg, "-") {
			args = append(args, arg)
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) != 2 {
		fmt.Println("Needs one path and one archive")
		return
	}

	var currentPath string
	hasPath := false
	var archiveName string
	hasArchive := false
	for _, path := range paths {
		candidateArchive := filepath.Base(path)
		extensionFile := strings.Split(candidateArchive, ".")
		if len(extensionFile) > 1 && extensionFile[1] == "zip" {
			archiveName = candidateArchive
			hasArchive = true
		} else {
			currentPath = path
			hasPath = true
		}
	}

	if !hasArchive {
		fmt.Println("Needs to input archive")
		return
	}

	if !hasPath {
		fmt.Println("Needs to input path ('.' - this directory)")
		return
	}

	if currentPath == "." {
		currentPath, _ = os.Getwd()
	} else {
		currentPath = findAbsolutePath(currentPath)
	}

	if currentPath == "" {
		fmt.Println("Incorrect path")
		return
	}

	if !isRealArchive(archiveName) {
		fmt.Println("Need the archive name")
		return
	}

	if uncompressedArchive(archiveName, currentPath, *listFilesFlag) != nil {
		fmt.Println("Error unpacking archive")
	} else {
		if *testErrorFlag {
			fmt.Println("Errors not found")
		}
		if *commentFlag {
			fmt.Println("Unpacking successfully")
		}
	}

}

func uncompressedArchive(src, dest string, listFiles bool) error {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		err := os.MkdirAll(dest, 0755)
		if err != nil {
			return err
		}
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if listFiles {
			if f.Name == "/" {
				continue
			}
			fmt.Println(f.Name)
		}
		err = uncompressFile(f, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func uncompressFile(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	path := filepath.Join(dest, f.Name)

	if f.FileInfo().IsDir() {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		directory := filepath.Dir(path)
		err = os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return err
		}

		dstFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

func findAbsolutePath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return ""
	}

	if !strings.HasPrefix(absPath, "/") {
		return ""
	}

	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return ""
	}

	return absPath
}

func isRealArchive(archiveName string) bool {
	inputArray := strings.Split(archiveName, ".")
	if len(inputArray) != 2 {
		return false
	}

	if inputArray[1] != "zip" {
		return false
	}

	return true
}
