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

	if len(args) < 2 {
		help.PrintHelp()
		return
	}

	filePath := args[0]
	command := strings.ToUpper(args[1])

	switch command {
	case "-E":
		logic.EncryptFile(filePath)
	case "-D":
		logic.DecryptFile(filePath)
	default:
		fmt.Println("Invalid command")
		help.PrintHelp()
	}

}
