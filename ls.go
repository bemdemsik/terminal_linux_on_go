package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	recursiveFlag := flag.Bool("R", false, "Recursively list files and directories")
	allFilesFlag := flag.Bool("a", false, "Output all files, with hide")
	longFormatFlag := flag.Bool("l", false, "To get long format")
	customFormatFlag := flag.Bool("H", false, "Custom format for user")
	sortedFlag := flag.Bool("r", false, "Sorted by desc")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	inputArray := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: ls [options] [file]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var files []string

	for i := 1; i < len(inputArray); i++ {
		if strings.HasPrefix(inputArray[i], "-") {
			continue
		}

		files = append(files, inputArray[i])
	}

	if len(files) > 1 {
		fmt.Println("Too many values")
		return
	}

	var dir string

	if len(files) < 1 {
		dir, _ = filepath.Abs(".")
	} else {
		_, err := os.Stat(files[0])
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Should be a path")
				return
			}
		}
		dir = files[0]
	}

	listFiles(dir, *recursiveFlag, " ", *allFilesFlag, *longFormatFlag, *customFormatFlag, *sortedFlag)
}

func listFiles(dir string, recursive bool, indent string, withHide bool, longFormat bool, customFormat bool, sorted bool) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error read files in directory")
		return
	}

	if sorted {
		files = reverseFiles(files)
	}

	for _, file := range files {

		if isHideFileWithAccess(withHide, file) {
			continue
		}

		fileInfo, err := os.Lstat(filepath.Join(dir, file.Name()))
		if err != nil {
			fmt.Println("Error reading information about directory")
			return
		}

		mode := fileInfo.Mode()
		size := fileInfo.Size()
		modTime := fileInfo.ModTime().Format("Jan _2 15:04")
		name := file.Name()
		updateSize := humanReadableSize(size)

		if longFormat && customFormat {
			fmt.Printf("%s %s %s %s %s\n", indent, mode, updateSize, modTime, name)
		} else if longFormat {
			fmt.Printf("%s %s %6d %s %s\n", indent, mode, size, modTime, name)
		} else if customFormat {
			fmt.Printf("%s %s %s\n", indent, name, updateSize)
		} else {
			fmt.Printf("%s %s\n", indent, name)
		}

		if isRecursive(recursive, file) {
			subdirectory := filepath.Join(dir, file.Name())
			listFiles(subdirectory, recursive, indent+"  ", withHide, longFormat, customFormat, false)
		}
	}
}

func reverseFiles(files []os.DirEntry) []os.DirEntry {
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}

	return files
}

func humanReadableSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func isRecursive(recursive bool, file os.DirEntry) bool {
	if file.IsDir() && recursive {
		return true
	}

	return false
}

func isHideFileWithAccess(allFiles bool, file os.DirEntry) bool {
	if !allFiles && len(file.Name()) > 0 && file.Name()[0] == '.' {
		return true
	}

	return false
}
