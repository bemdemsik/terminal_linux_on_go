package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	names := os.Args
	flags := []string{"-o", "-h", "-p", "-n"}

	if names[0] != "./uname" {
		fmt.Println("Неправильное использование команды")
		return
	}
	found := false
	if len(names) > 1 {
		for _, v1 := range names {
			for _, v2 := range flags {
				if v1 == v2 {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	} else {
		found = true
	}
	if !found {
		fmt.Println("Ошибка в использовании команды")
	} else {
		if contains(names, "-o") || len(names) == 1 {
			fmt.Println("операционная система: ", runtime.GOOS)
		}
		if contains(names, "-n") || len(names) == 1 {
			hostname, err := os.Hostname()
			if err == nil {
				fmt.Println("имя компьютера: ", hostname)
			}
		}
		if contains(names, "-p") || len(names) == 1 {
			fmt.Println("количество CPUs: ", runtime.NumCPU())
		}
		if contains(names, "-h") {
			fmt.Println("Использование: uname [ПАРАМЕТР]\nПечатает определенные сведения о системе\n-n - отобразить имя машины\n-o - отобразить название операционной системы\n-p - отобразить количество ядер процессора\n")
		}
		if len(names) == 1 {
			fmt.Println("архитектура: ", runtime.GOARCH)
		}
	}
}
func contains(l []string, element string) bool {
	for _, item := range l {
		if item == element {
			return true
		}
	}
	return false
}
