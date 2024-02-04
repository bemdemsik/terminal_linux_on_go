package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type wcCount struct {
	FileName   string
	Characters int
	Lines      int
	Words      int
}

func main() {

	// Определение флага
	charactersFlag := flag.Bool("m", false, "Characters")
	linesFlag := flag.Bool("l", false, "Lines")
	wordsFlag := flag.Bool("w", false, "Words")
	helpFlag := flag.Bool("h", false, "Help")

	totalCharacters := 0
	totalLines := 0
	totalWords := 0

	flag.Parse()
	array := os.Args[1:]

	var onlyArgs []string
	var values []string

	for _, arg := range array {
		if arg == "wc" || arg == "./wc" {
			continue
		}

		if strings.HasPrefix(arg, "-") {
			onlyArgs = append(onlyArgs, arg)
		} else {
			values = append(values, arg)
		}
	}
	if len(values) < 1 && !*helpFlag {
		fmt.Println("Неверный синтаксис. Укажите файл")
		return
	}
	flag.CommandLine.Parse(onlyArgs)
	//символы
	if *charactersFlag {
		fmt.Printf("%s\n", "Characters")
		for _, fileName := range values {
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Printf("Error opening file %s: %s\n", fileName, err)
				continue
			}
			defer file.Close()
			characters := countCharacters(file)
			fmt.Printf("%d\t%s\n", characters, fileName)
		}
	}
	//строки
	if *linesFlag {
		fmt.Printf("%s\n", "Lines")
		for _, fileName := range values {
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Printf("Error opening file %s: %s\n", fileName, err)
				continue
			}
			defer file.Close()
			lines := countLines(file)
			fmt.Printf("%d\t%s\n", lines, fileName)
		}
	}
	//слова
	if *wordsFlag {
		fmt.Printf("%s\n", "Words")
		for _, fileName := range values {
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Printf("Error opening file %s: %s\n", fileName, err)
				continue
			}
			defer file.Close()
			words := countWords(file)
			fmt.Printf("%d\t%s\n", words, fileName)
		}

	}
	//просто wc
	if *charactersFlag == false && *linesFlag == false && *wordsFlag == false && !*helpFlag {

		wcCounts := make([]wcCount, len(values))
		for i, fileName := range values {
			file, err := os.Open(fileName)
			if err != nil {
				fmt.Printf("Error opening file \n")
				continue
			}
			defer file.Close()

			wcCounts[i] = wc(fileName)
			totalCharacters += wcCounts[i].Characters
			totalLines += wcCounts[i].Lines
			totalWords += wcCounts[i].Words
		}

		fmt.Printf("%s\t\t%s\t\t%s\n", "Lines", "Words", "Characters")

		for _, wcCount := range wcCounts {
			fmt.Printf("%d\t\t%d\t\t%d\t\t%s\n", wcCount.Lines, wcCount.Words, wcCount.Characters, wcCount.FileName)
		}
		fmt.Printf("%d\t\t%d\t\t%d\t\ttotal\n", totalLines, totalWords, totalCharacters)
	}
	if *helpFlag && len(onlyArgs) == 1 {
		fmt.Print("Утилита может обрабатывать файлы. Стандартная инструкция выглядит так:\nwc file\n" +
			"wc — имя утилиты;\nfile — название обрабатываемого файла.\n\n")
		fmt.Print("-m\t--count\tПоказать количесто символов в объекте\n-l\t--lines\tВывести количество строк в объекте" +
			"\n-w\t--words\tОтобразить количество слов в объекте\n")

	}

}
func wc(filename string) wcCount {
	file, err := os.Open(filename)
	defer file.Close()

	characters := countCharacters(file)
	_, err = file.Seek(0, 0) // Reset file pointer to the beginning
	if err != nil {
		log.Fatal(err)
	}

	lines := countLines(file)
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	words := countWords(file)

	return wcCount{
		FileName:   filename,
		Characters: characters,
		Lines:      lines,
		Words:      words,
	}
}

// считаем символы
func countCharacters(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	count := 0

	for scanner.Scan() {
		count++
	}

	return count
}

// строки
func countLines(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	count := 0

	for scanner.Scan() {
		count++
	}
	return count
}

// слова
func countWords(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	count := 0

	for scanner.Scan() {
		count++
	}
	return count
}
