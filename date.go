package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	names := os.Args
	flags := []string{"-s", "-h", "-w", "-d"}

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
		if contains(names, "-s") {
			currentTime := time.Now()
			fmt.Println(currentTime.Format("2006-01-02"))
		}
		if contains(names, "-w") {
			currentTime := time.Now()
			previousDate := currentTime.AddDate(0, 0, +7)
			fmt.Println(previousDate.Format("Jan Mon 2006-01-02"))
		}
		if contains(names, "-d") {
			currentTime := time.Now()
			dayYear := currentTime.YearDay()
			fmt.Println("Сегодняшний день в году: ", dayYear)
		}
		if contains(names, "-h") {
			fmt.Println("date [ПАРАМЕТР]\nВыводит текущее время, или изменяет время в системе\n-s - сокращенный формат даты\n-w - показывает дату, которая будет через неделю\n-d - показывает какой по счету сегодняшний день в году")
		}
		if len(names) == 1 {
			currentTime := time.Now()
			fmt.Println(currentTime.Format("Mon Jan 02 15:04:05 MST"))
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
