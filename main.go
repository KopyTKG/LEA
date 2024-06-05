package main

import (
	"fmt"
	"lea/help"
	"lea/logic"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		help.PrintHelp()
		return
	}

	var filePath string
	validCommandFound := false

	for _, arg := range args {
		switch {
		case arg == "-E" || arg == "-D" || arg == "-C":
			validCommandFound = true
			switch arg {
			case "-E":
				if filePath != "" {
					logic.EncryptFile(filePath)
				} else {
					fmt.Println("No file path provided for encryption.")
				}
			case "-D":
				if filePath != "" {
					logic.DecryptFile(filePath)
				} else {
					fmt.Println("No file path provided for decryption.")
				}
			case "-C":
				logic.GenerateConstants()
			}
		default:
			if strings.Contains(arg, ".") {
				filePath = arg
				validCommandFound = true
			}
		}
	}

	if !validCommandFound {
		fmt.Println("Invalid command or file path")
		help.PrintHelp()
	}
}
