package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
)

func main() {
	megaFlag := flag.Bool("mega", false, "Display memory in megabytes")
	kiloFlag := flag.Bool("kilo", false, "Display memory in kilobytes")
	byteFlag := flag.Bool("b", false, "Display memory in bytes")
	helpFlag := flag.Bool("h", false, "Help")

	flag.Parse()
	inputArray := os.Args[1:]

	if *helpFlag {
		fmt.Println("Usage: free [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	var args []string

	for _, arg := range inputArray {
		if arg == "free" {
			continue
		}
		if strings.HasPrefix(arg, "-") {
			args = append(args, arg)
		} else {
			break
		}
	}

	if len(args) > 1 {
		fmt.Println("Arguments are incongruous")
		return
	}

	var stat syscall.Sysinfo_t
	err := syscall.Sysinfo(&stat)

	if err != nil {
		fmt.Println("Error reading memory")
		return
	}

	if *megaFlag {
		displayMemoryInfo("megabytes", stat)
	} else if *kiloFlag {
		displayMemoryInfo("kilobytes", stat)
	} else if *byteFlag {
		displayMemoryInfo("bytes", stat)
	} else {
		displayMemoryInfo("bytes", stat)
	}
}

func displayMemoryInfo(unit string, stat syscall.Sysinfo_t) {
	var total float64
	var used float64
	var free float64
	var shared float64
	var buffer float64
	var cached float64

	if unit == "megabytes" {
		total = float64(stat.Totalram) / 1024 / 1024
		used = (float64(stat.Totalram) - float64(stat.Freeram)) / 1024 / 1024
		free = float64(stat.Freeram) / 1024 / 1024
		shared = float64(stat.Sharedram) / 1024 / 1024
		buffer = float64(stat.Bufferram) / 1024 / 1024
		cached = calculateCachedMemory("megabytes")
	} else if unit == "kilobytes" {
		total = float64(stat.Totalram) / 1024
		used = (float64(stat.Totalram) - float64(stat.Freeram)) / 1024
		free = float64(stat.Freeram) / 1024
		shared = float64(stat.Sharedram) / 1024
		buffer = float64(stat.Bufferram) / 1024
		cached = calculateCachedMemory("kilobytes")
	} else {
		total = float64(stat.Totalram)
		used = float64(stat.Totalram) - float64(stat.Freeram)
		free = float64(stat.Freeram)
		shared = float64(stat.Sharedram)
		buffer = float64(stat.Bufferram)
		cached = calculateCachedMemory("bytes")
	}

	fmt.Printf("Total memory:    %f %s\n", total, unit)
	fmt.Printf("Used memory:     %f %s\n", used, unit)
	fmt.Printf("Free memory:     %f %s\n", free, unit)
	fmt.Printf("Shared memory:   %f %s\n", shared, unit)
	fmt.Printf("Buffer memory:   %f %s\n", buffer, unit)
	fmt.Printf("Cached memory:   %f %s\n", cached, unit)
}

func calculateCachedMemory(unitDivisor string) float64 {
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		fmt.Println("Failed to get memory info:", err)
		return 0
	}

	var cachedMemory float64

	if unitDivisor == "megabytes" {
		cachedMemory = (float64(info.Totalram) - float64(info.Freeram) - float64(info.Bufferram)) / 1024 / 1024
	} else if unitDivisor == "kilobytes" {
		cachedMemory = (float64(info.Totalram) - float64(info.Freeram) - float64(info.Bufferram)) / 1024
	} else {
		cachedMemory = float64(info.Totalram) - float64(info.Freeram) - float64(info.Bufferram)
	}

	return cachedMemory
}
