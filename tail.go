package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Определение флага
	numBytes := flag.Int("c", 0, "Number of bytes to print")
	numLines := flag.Bool("n", false, "Number of lines to print")
	quiet := flag.Bool("q", false, "No file name")
	helpFlag := flag.Bool("h", false, "Help")

	numberLines := 10
	numberBytes := 0
	flag.Parse()
	array := os.Args[1:]

	var onlyArgs []string
	var values []string
	var er error
	for i, arg := range array {
		if arg == "tail" || arg == "./tail" {
			continue
		}

		if i > 0 && array[i-1] == "-n" {
			numberLines, er = strconv.Atoi(arg)
			if er != nil {
				fmt.Println("Неверный синтаксис. Ввели не число")
				return
			}
			continue
		}
		if i > 0 && array[i-1] == "-c" {
			numberBytes, er = strconv.Atoi(arg)
			if er != nil {
				fmt.Println("Неверный синтаксис. Введите кол-во байтов")
				return
			}
			continue
		}
		if strings.HasPrefix(arg, "-") {
			onlyArgs = append(onlyArgs, arg)
		} else {
			values = append(values, arg)
		}
	}
	if len(values) < 1 {
		fmt.Println("Неверный синтаксис. Укажите файл")
		return
	}

	for _, fileName := range values {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("Error opening file %s: %s\n", fileName, err)
			continue
		}
		if *numBytes > 0 {

			defer file.Close()
			stat, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}
			fileSize := stat.Size()
			if fileSize < int64(numberBytes) {
				numberBytes = int(fileSize)
			}
			file.Seek(-int64(*numBytes), io.SeekEnd)
			io.Copy(os.Stdout, file)

		} else if *helpFlag && len(onlyArgs) == 1 {
			fmt.Print("tail опции файл\nПо умолчанию утилита выводит десять последних строк из файла. Опции:" +
				"\n\t-c - выводить указанное количество байт с конца файла;\n" +
				"\t-n - выводить указанное количество строк из конца файла;\n" +
				"\t-q - не выводить имена файлов;\n" +
				"\t-h - выводить информацию о команде;\n")
		} else {
			if !*quiet && len(array) > 1 {
				fmt.Println("==> ", fileName, " <==")
			}
			if *numLines {
				printLinesFromFile(file, numberLines)
			} else {
				printLinesFromFile(file, 10)
			}
		}

	}

}
func printLinesFromFile(file *os.File, numLines int) {

	scanner := bufio.NewScanner(file)

	// Создаем срез для хранения строк
	lines := make([]string, 0)

	// Считываем строки из файла
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Если количество строк в файле меньше, чем запрошенное, выводим все строки
	if len(lines) <= numLines {
		for _, line := range lines {
			fmt.Println(line)
		}
		return
	}

	// Выводим заданное количество строк с конца
	for i := len(lines) - numLines; i < len(lines); i++ {
		fmt.Println(lines[i])
	}
}

func readLinesFromBuffer(buffer []byte) []string {
	lines := make([]string, 0)
	for i := len(buffer) - 1; i >= 0; i-- {
		if buffer[i] == '\n' {
			lines = append(lines, string(buffer[i+1:]))
		}
	}
	return lines
}
