package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	helpFlag := flag.Bool("h", false, "Help")
	sizeFlag := flag.Bool("s", false, "Show the size of file")
	flag.Parse()
	args := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: file [options] <file-path>")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var onlyArgs []string
	var files []string

	for _, arg := range args {
		if arg == "./file" || arg == "file" {
			continue
		}
		if strings.HasPrefix(arg, "-") {
			onlyArgs = append(onlyArgs, arg)
		} else {
			files = append(files, arg)
		}
	}

	if len(files) != 1 {
		fmt.Println("Need the input one file")
		return
	}
	file := files[0]

	contentType := detectFileType(file)

	fmt.Println(contentType)

	if *sizeFlag {
		fmt.Println(getSizeOfFile(file))
	}

}

func getSizeOfFile(file string) string {
	fileInfo, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}

	fileSize := fileInfo.Size()
	return fmt.Sprintf("%d bytes", fileSize)
}

func detectFileType(filePath string) string {
	fileExtension := getFileExtension(filePath)
	switch strings.ToLower(fileExtension) {
	case ".txt":
		return "Текстовый документ"
	case ".jpg", ".jpeg", ".png":
		return "Изображение"
	case ".mp3", ".wav":
		return "Аудиофайл"
	case ".mp4", ".avi", ".mkv":
		return "Видеофайл"
	case ".pdf":
		return "PDF-документ"
	default:
		return "Неизвестный тип файла"
	}
}

func getFileExtension(filePath string) string {
	lastDotIndex := strings.LastIndex(filePath, ".")
	if lastDotIndex != -1 {
		return filePath[lastDotIndex:]
	}
	return ""
}
