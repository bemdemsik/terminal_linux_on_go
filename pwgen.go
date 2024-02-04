package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var enterFlags []string
var commandFlags []string
var value []string

var lenPassword int
var countPassword int

var er error

const (
	useSymbol      = "qwerafs7dgtyzh5nxcm9vk4ioulj6pQWP3OW8EQNDJZM2XC0VCZL1AKSDHR"
	specialSymbols = "!@#$%^&*()"
)

func main() {
	pwgen(os.Args[1:])
}

func pwgen(arguments []string) {
	commandFlags = append(commandFlags, "l", "y", "n", "h")
	if !fill(arguments) {
		return
	}

	if contains(enterFlags, "h") {
		if len(enterFlags) > 1 || len(value) > 0 {
			fmt.Println("Invalid key -h")
		} else {
			fmt.Println("-l - длина пароля")
			fmt.Println("-y - пароль со спец.символом")
			fmt.Println("-n - иколичество паролей")
			fmt.Println("-h - помощь")
		}
		return
	}

	if contains(enterFlags, "l") {
		lenPassword, er = strconv.Atoi(arguments[indexOf(arguments, "-l")+1])
		if er != nil {
			fmt.Println("Неверный формат длины пароля: ", arguments[indexOf(arguments, "-l")+1])
			return
		}
	} else {
		lenPassword = 8
	}
	if contains(enterFlags, "n") {
		if indexOf(arguments, "-n")+1 == len(arguments) {
			fmt.Println("Неверный формат количества паролей: ", arguments[indexOf(arguments, "-n")])
			return
		}
		countPassword, er = strconv.Atoi(arguments[indexOf(arguments, "-n")+1])
		if er != nil {
			fmt.Println("Неверный формат количества паролей: ", arguments[indexOf(arguments, "-n")+1])
			return
		}
	} else {
		countPassword = 160
	}
	passwords, e := generatePasswords()
	if e != nil {
		fmt.Println(e)
	} else {
		for i, item := range passwords {
			fmt.Print(item, "\t")

			if (i+1)%8 == 0 {
				fmt.Println("")
			}
		}

		if len(passwords)%8 != 0 {
			fmt.Println("")
		}
	}
}
func generatePasswords() ([]string, error) {
	pswd := []string{}
	for i := 0; i < countPassword; i++ {
		password := make([]byte, lenPassword)
		for j := 0; j < lenPassword; j++ {
			source := rand.NewSource(time.Now().UnixNano())
			generator := rand.New(source)
			index := generator.Intn(len(useSymbol) - 1)
			password[j] = useSymbol[index]
		}
		if contains(enterFlags, "y") {
			source := rand.NewSource(time.Now().UnixNano())
			generator := rand.New(source)
			indexP := generator.Intn(len(password) - 1)
			indexS := generator.Intn(len(specialSymbols) - 1)
			password[indexP] = specialSymbols[indexS]
		}
		pswd = append(pswd, string(password))
	}
	return pswd, nil
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
func indexOf(list []string, element string) int {
	for i := 0; i < len(list); i++ {
		if list[i] == element {
			return i
		}
	}
	return -1
}
