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
	createFlag := flag.Bool("c", true, "Creating archive")
	extractFlag := flag.Bool("x", false, "Extract files from archive")
	informationFlag := flag.Bool("f", false, "Full information about unpacking")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	inputArray := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: tar [options] [file ...]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var args []string

	var paths []string
	for _, arg := range inputArray {

		if strings.HasPrefix(arg, "-") {
			args = append(args, arg)
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) < 2 {
		fmt.Println("Needs path and archive name")
		return
	}

	var currentPaths []string
	var folderName string
	hasPath := false
	var pathToArchive string
	var archiveName string
	hasArchive := false
	for _, path := range paths {
		candidateArchive := filepath.Base(path)
		pathToArchive = path
		extensionFile := strings.Split(candidateArchive, ".")
		if len(extensionFile) > 1 && extensionFile[1] == "zip" {
			archiveName = candidateArchive
			folderName = extensionFile[0]
			hasArchive = true
		} else {
			currentPaths = append(currentPaths, path)
			hasPath = true
		}
	}

	if !hasArchive {
		fmt.Println("Needs to input archive")
		return
	}

	if !hasPath {
		fmt.Println("Needs to input path(s) ('.' - this directory)")
		return
	}

	if *extractFlag && len(currentPaths) > 1 {
		fmt.Println("When unzipping need to specify one path")
		return
	}

	hasError := false

	if !*extractFlag {
		for i, currentPath := range currentPaths {
			var path string
			if currentPath == "." {
				path, _ = os.Getwd()
			} else {
				path = existPath(filepath.Clean("/" + currentPath + "/"))
				if path == "" {
					path = filepath.Clean(currentPath)
				}
			}

			if path == "" {
				fmt.Println("Incorrect path")
				hasError = true
				continue
			} else {
				currentPaths = append(currentPaths[:i], currentPaths[i+1:]...)
				currentPaths = append(currentPaths, path)
			}

			if !isCorrectArchive(archiveName) {
				fmt.Println("Need the archive name")
				hasError = true
				continue
			}
		}
	}

	if hasError {
		return
	}

	if *extractFlag {
		var path string
		if currentPaths[0] == "." {
			workingDirectory, _ := os.Getwd()
			path = filepath.Join(workingDirectory, folderName)
		} else {
			path = filepath.Join(currentPaths[0], folderName)
		}

		if uncompressArchive(archiveName, path, *informationFlag) != nil {
			fmt.Println("Error unpacking archive")
		}
	} else {
		if *createFlag {
			if createNewArchive(archiveName, currentPaths, *informationFlag) != nil {
				fmt.Println("Error packing archive")
			}
		} else {
			if !isExistArchive(pathToArchive) {
				fmt.Println("Archive not found")
				return
			}

			if createNewArchive(archiveName, currentPaths, *informationFlag) != nil {
				fmt.Println("Error packing archive")
			}
		}
	}
}

func uncompressArchive(src, dest string, needInformation bool) error {
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
		err = uncompressingFile(f, dest, needInformation)
		if err != nil {
			return err
		}
	}

	return nil
}

func uncompressingFile(f *zip.File, dest string, needInformation bool) error {
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

		if needInformation {
			fmt.Print("Unpacking " + dstFile.Name() + " 0%...")
		}

		_, err = io.Copy(dstFile, rc)
		if err != nil {
			return err
		}

		if needInformation {
			fmt.Println("100%")
		}
	}

	return nil
}

func existPath(path string) string {
	workDir, _ := os.Getwd()
	absPath := filepath.Clean(workDir + "/" + path)

	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return ""
	}

	return absPath

}

func isCorrectArchive(archiveName string) bool {
	inputArray := strings.Split(archiveName, ".")
	if len(inputArray) != 2 {
		return false
	}

	if inputArray[1] != "zip" {
		return false
	}

	return true
}

func isExistArchive(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func createNewArchive(archiveName string, sourceFolders []string, needInformation bool) error {
	var archiveFile *os.File
	if _, err := os.Stat(archiveName); os.IsNotExist(err) {
		archiveFile, err = os.Create(archiveName)
		if err != nil {
			return err
		}
		defer archiveFile.Close()
	} else {
		archiveFile, err = os.OpenFile(archiveName, os.O_RDWR, 0644)
		if err != nil {
			return err
		}
	}

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for _, folder := range sourceFolders {
		err := filepath.Walk(folder, func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relativePath := strings.TrimPrefix(filePath, folder)

			if needInformation {
				fmt.Print("Archiving " + relativePath + " 0%...")
			}

			if fileInfo.IsDir() {
				_, err := zipWriter.CreateHeader(&zip.FileHeader{
					Name:     relativePath + "/",
					Method:   zip.Store,
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
					Method:   zip.Store,
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
			if needInformation {
				fmt.Println("100%")
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
