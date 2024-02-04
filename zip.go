package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	deleteFromArchiveFlag := flag.Bool("d", false, "Delete archive")
	recursiveFlag := flag.Bool("r", false, "Archiving recursively")
	withoutCompress := flag.Bool("0", false, "Archiving without compression files")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	inputArray := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: command zip [options] [path]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var args []string

	var paths []string
	for _, arg := range inputArray {
		if arg == "zip" || arg == "./zip" {
			continue
		}

		if strings.HasPrefix(arg, "-") {
			args = append(args, arg)
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) != 1 {
		fmt.Println("Needs one path")
		return
	}

	var currentPath string

	currentPath = getAbsolutePath(paths[0])

	if currentPath == "" {
		fmt.Println("Incorrect path")
		return
	}

	archiveName := filepath.Base(currentPath)

	if *deleteFromArchiveFlag {
		if *recursiveFlag || *withoutCompress {
			fmt.Println("Flags are incompatible")
		}
	}

	if *deleteFromArchiveFlag {
		if !isArchive(archiveName) {
			fmt.Println("Needs to input archive")
			return
		}

		err := deleteArchive(currentPath)
		if err != nil {
			fmt.Println("Error deleting archive")
		}
		return
	} else {
		if isArchive(archiveName) {
			fmt.Println("Needs to input directory")
			return
		} else {
			archiveName += ".zip"
		}
	}

	if !*recursiveFlag {
		err := createArchiveWithoutRecursion(currentPath, archiveName, !*withoutCompress)
		if err != nil {
			fmt.Println("Error archive directory")
			return
		}
	} else {
		err := createArchive(archiveName, currentPath, !*withoutCompress)
		if err != nil {
			fmt.Println("Error archive directory")
			return
		}
	}

}

func getAbsolutePath(path string) string {
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

func isArchive(archiveName string) bool {
	inputArray := strings.Split(archiveName, ".")
	if len(inputArray) != 2 {
		return false
	}

	if inputArray[1] != "zip" {
		return false
	}

	return true
}

func deleteArchive(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func createArchiveWithoutRecursion(path string, archiveName string, withCompress bool) error {
	archiveFile, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	dirInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !dirInfo.IsDir() {
		return errors.New("provided path is not a directory")
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(path, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		var method uint16

		if withCompress {
			method = zip.Deflate
		} else {
			method = zip.Store
		}

		zipFile, err := zipWriter.CreateHeader(&zip.FileHeader{
			Name:     file.Name(),
			Method:   method,
			Modified: fileInfo.ModTime(),
		})
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFile, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func createArchive(archiveName string, sourceFolder string, withCompress bool) error {
	archiveFile, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	var method uint16

	if withCompress {
		method = zip.Deflate
	} else {
		method = zip.Store
	}

	err = filepath.Walk(sourceFolder, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath := strings.TrimPrefix(filePath, sourceFolder)

		if fileInfo.IsDir() {
			_, err := zipWriter.CreateHeader(&zip.FileHeader{
				Name:     relativePath + "/",
				Method:   method,
				Modified: fileInfo.ModTime(),
			})
			if err != nil {
				return err
			}
		} else {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			zipFile, err := zipWriter.CreateHeader(&zip.FileHeader{
				Name:     relativePath,
				Method:   method,
				Modified: fileInfo.ModTime(),
			})
			if err != nil {
				return err
			}

			_, err = io.Copy(zipFile, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
