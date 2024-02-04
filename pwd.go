package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	versionFlag := flag.Bool("v", false, "Showing version of command")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	allArgs := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: pwd [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var args []string

	for _, arg := range allArgs {
		if strings.HasPrefix(arg, "-") {
			args = append(args, arg)
		}
	}

	currentPath, _ := os.Getwd()

	fmt.Println(currentPath)

	if *versionFlag {
		fmt.Println("Version 0.12.3.54322")
	}
}
