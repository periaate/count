package main

import (
	"bufio"
	"os"
)

func main() {
	i := len(os.Args) - 1
	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			i++
		}
	}
	println(i)
}
