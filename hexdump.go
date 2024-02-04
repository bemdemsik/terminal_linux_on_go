package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var enterFlags []string
var commandFlags []string
var value []string
var activeDirectory string
var err error

const bytesPerLine = 16

func main() {
	hexdump(os.Args[1:])
}
func hexdump(arguments []string) {
	commandFlags = append(commandFlags, "v", "C", "d", "h")
	if !fill(arguments) {
		return
	}
	if contains(enterFlags, "h") {
		if len(enterFlags) > 1 || len(value) > 0 {
			fmt.Println("Invalid key -h")
		} else {
			fmt.Println("-v - показать версию утилиты hexdump")
			fmt.Println("-C - каноническое отображение hex + ASCII.")
			fmt.Println("-d - двухбайтовое десятичнео представление")
			fmt.Println("-h - помощь")
		}
		return
	}
	if contains(enterFlags, "v") {
		if len(enterFlags) > 1 || len(value) > 0 {
			fmt.Println("Invalid key -v")
		} else {
			fmt.Println("hexdump version 2.9.20")
		}
		return
	}
	if len(value) == 0 {
		fmt.Println("Usage: hexdump [flags] <files>")
		return
	}
	for _, filename := range value {
		if _, err = os.Stat(activeDirectory + "/" + filename); err == nil {
			filename = activeDirectory + filename
		} else if _, err = os.Stat(filename); err != nil {
			fmt.Println("Данный файл на найден: ", filename)
			continue
		}
		if contains(enterFlags, "d") {
			showInDecimal(filename)
		} else {
			show(filename)
		}
	}
}
func show(filename string) {
	data, er := ioutil.ReadFile(filename)
	if er != nil {
		fmt.Println("Ошибка при чтении файла: ", filename)
	}

	fmt.Printf("Hex dump of %s:\n", filename)
	const bytesPerLine = 16
	totalBytes := len(data)

	for i := 0; i < totalBytes; i += bytesPerLine {
		end := i + bytesPerLine
		if end > totalBytes {
			end = totalBytes
		}
		hexDump := hex.Dump(data[i:end])
		hexd := ""
		if contains(enterFlags, "C") {
			hexd = strings.Replace(hexDump, "00000000 ", "", -1)
		} else {
			hexd = strings.Split(strings.Replace(hexDump, "00000000 ", "", -1), "|")[0] + "\n"
		}
		fmt.Printf("%08x  %s", i, hexd)
	}
}
func showInDecimal(filename string) {
	data, er := ioutil.ReadFile(filename)
	if er != nil {
		log.Fatal(er)
	}

	fmt.Printf("Decimal dump of %s:\n", filename)
	totalBytes := len(data)

	for i := 0; i < totalBytes; i += bytesPerLine {
		end := i + bytesPerLine
		if end > totalBytes {
			end = totalBytes
		}
		line := data[i:end]
		ascii := make([]string, len(line))
		decimal := make([]string, len(line))
		for idx, val := range line {
			ascii[idx] = string(val)
			decimal[idx] = fmt.Sprintf("%d", val)
		}
		fmt.Printf("%08d  %-24s", i, strings.Join(decimal, " "))
		if len(line) < bytesPerLine {
			padding := bytesPerLine - len(line)
			fmt.Print(strings.Repeat("   ", padding))
		}
		if contains(enterFlags, "C") {
			fmt.Printf("  %s\n", strings.Replace(strings.Join(ascii, ""), "\n", "..", -1))
		} else {
			fmt.Println()
		}
	}
}
func fill(arguments []string) bool {
	enterFlags = []string{}
	value = []string{}
	for i := 0; i < len(arguments); i++ {
		if arguments[i][0] == '-' {
			if !contains(commandFlags, string(arguments[i][1])) {
				fmt.Println("Invalid key " + arguments[i])
				return false
			} else {
				enterFlags = append(enterFlags, string(arguments[i][1]))
			}
		} else {
			value = append(value, arguments[i])
		}
	}
	return true
}
func contains(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}
