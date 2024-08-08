package main

import (
	"fmt"
	"lea/help"
	"lea/modes"
	"lea/stream"
	"lea/utils"
	"log"
	"os"
	"strings"
)

var mode string = "ecb"
var key int = 128

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		help.PrintHelp()
		return
	}

	handleArgs(args)
}

/* Main handeling loop for arguments */
func handleArgs(args []string) {
	// preset file paths
	filePath, keyPath, seedPath := "", "", ""
	argsList := utils.List{}

	// processing loop for args
	processArguments(args, &argsList, &filePath, &keyPath, &seedPath)

	if argsList.Length() == 0 && filePath != "" {
		argsList.Append("-e")
	}

	validCommandFound, encrypted := false, false
	processCommands(argsList, filePath, keyPath, seedPath, &validCommandFound, &encrypted)

	if !validCommandFound {
		fmt.Println("Invalid command or file path")
		help.PrintHelp()
	}

	if !encrypted && filePath != "" && keyPath != "" && seedPath != "" {
		executeMode(filePath, keyPath, seedPath, mode, "-e")
	}
}

func processArguments(args []string, argsList *utils.List, filePath, keyPath, seedPath *string) {
	prev := ""
	for _, arg := range args {
		switch {
		// encrypt command must be last
		case arg == "-e" || arg == "-d" || arg == "--encrypt" || arg == "--decrypt":
			argsList.Append(arg)

		// signal for key / seed file load
		case arg == "-ek" || arg == "-es" || arg == "--external-key" || arg == "--external-seed":
			prev = arg

		// seed / key handeling
		case prev == "-ek" || prev == "--external-key":
			prev = ""
			*keyPath = arg
		case prev == "-es" || prev == "--external-seed":
			prev = ""
			*seedPath = arg
		
		// source file handeling
		case strings.Contains(arg, ".") && prev == "":
			*filePath = arg

		// console output for version and help
		case arg == "-h" || arg == "--help":
			help.PrintHelp()
			os.Exit(1)
		case arg == "-v" || arg == "--version":
			help.Version()
			os.Exit(1)

		// Cypher modes
		case arg == "--ecb":
			mode = "ecb"
		case arg == "--cbc":
			mode = "cbc"
		case arg == "--cfb":
			mode = "cfb"
		case arg == "--ofb":
			mode = "ofb"

		// Key lenght
		case arg == "--128":
			key = 128
		case arg == "--192":
			key = 192
		case arg == "--256":
			key = 256
		}
	}
}


func processCommands(argsList utils.List, filePath, keyPath, seedPath string, validCommandFound, encrypted *bool) {
	for _, arg := range argsList.Elements {
		switch arg {
		case "-e", "-d", "--encrypt", "--decrypt":
			*validCommandFound = true
			*encrypted = true
			executeMode(filePath, keyPath, seedPath, mode, arg)

		case "--external-key", "--external-seed", "-ek", "-es":
			*validCommandFound = true
		}
	}
}

func executeMode(filePath, keyPath, seedPath string, mode string, command string) {
	var encrypt bool = false
	bKey, bSeed := stream.GetFile(keyPath), stream.GetFile(seedPath)

	if command == "-e" || command == "--encrypt"  {
		encrypt = true
	}
	if filePath == "" {
		log.Fatalln("No file path provided")
		help.PrintHelp()
		os.Exit(1)
	}

	switch {
	case mode == "ecb" || mode == "cbc" || mode == "cfb":
		modes.PerformMode(mode, filePath, bKey, bSeed, encrypt, key)
	case mode == "ofb":
		modes.PerformOFB(filePath, bKey, bSeed, encrypt, key)
	default:
		log.Fatalln("Invalid mode")
		help.PrintHelp()
		os.Exit(1)
	}
}
