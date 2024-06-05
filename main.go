package main

import (
	"fmt"
	"lea/help"
	"lea/logic"
	"lea/generator"
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
		case arg == "-E" || arg == "-D" || arg == "-C" || arg == "-K":
			validCommandFound = true
			switch arg {
			case "-E":
				if filePath != "" {
					key := logic.GetInternalKey()
					logic.EncryptFile(filePath, key)
				} else {
					fmt.Println("No file path provided for encryption.")
				}
			case "-D":
				if filePath != "" {
					key := logic.GetInternalKey()
					logic.DecryptFile(filePath, key)
				} else {
					fmt.Println("No file path provided for decryption.")
				}
			case "-C":
				generator.GenerateConstants()
			case "-K":
				generator.GenerateKey()
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
