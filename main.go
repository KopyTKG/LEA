package main

import (
	"fmt"
	"lea/generator"
	"lea/help"
	"lea/modes"
	"lea/utils"
	"os"
	"strings"
)

var mode string = "ecb"

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		help.PrintHelp()
		return
	}

	processArgs(args)
}

/* Main processing loop for arguments */
func processArgs(args []string) {
	// preset file paths
	filePath, keyPath, seedPath := "", "", ""
	argsList := utils.List{}

	//
	validCommandFound, encrypted := processArguments(args, &argsList, &filePath, &keyPath, &seedPath)

	if argsList.Length() == 0 && filePath != "" {
		argsList.Append("-e")
	}

	processCommands(argsList, filePath, keyPath, seedPath, &validCommandFound, &encrypted)

	if !validCommandFound {
		fmt.Println("Invalid command or file path")
		help.PrintHelp()
	}

	if !encrypted && filePath != "" {
		proccesMode(filePath, keyPath, seedPath, mode, "-e")
	}
}

func processArguments(args []string, argsList *utils.List, filePath, keyPath, seedPath *string) (bool, bool) {
	validCommandFound, encrypted := false, false
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
		}
	}
	return validCommandFound, encrypted
}

func generateKeyOrSeed(argsList *utils.List, arg string) {
	index := -1
	if arg == "--gen-seed" || arg == "-gs" {
		index = argsList.IndexOf("--external-seed")
		if index == -1 {
			index = argsList.IndexOf("-es")
		}
	}
	if arg == "--gen-key" || arg == "-gk" {
		index = argsList.IndexOf("--external-key")
		if index == -1 {
			index = argsList.IndexOf("-ek")
		}
	}
	if index == -1 {
		argsList.Prepend(arg)
	}
}

func handleExternalFiles(argsList *utils.List, arg string) {
	if argsList.Length() > 0 {
		if argsList.IndexOf("--gen-key") == -1 && (arg == "--external-key" || arg == "-ek") {
			argsList.Prepend(arg)
		}
		if argsList.IndexOf("--gen-seed") == -1 && (arg == "--external-seed" || arg == "-es") {
			argsList.Prepend(arg)
		}
	} else {
		argsList.Prepend(arg)
	}
}


func processCommands(argsList utils.List, filePath, keyPath, seedPath string, validCommandFound, encrypted *bool) {
	for _, arg := range argsList.Elements {
		switch arg {
		case "-e", "-d", "--encrypt", "--decrypt", "-gk", "-gs", "--gen-seed", "--gen-key":
			*validCommandFound = true
			executeCommand(arg, filePath, keyPath, seedPath, encrypted)

		case "--external-key", "--external-seed", "-ek", "-es":
			*validCommandFound = true
		}
	}
}

func executeCommand(command, filePath, keyPath, seedPath string, encrypted *bool) {
	switch command {
	case "-e", "--encrypt", "-d", "--decrypt":
		*encrypted = true
		proccesMode(filePath, keyPath, seedPath, mode, command)
	case "-gs", "--gen-seed":
		generator.GenerateConstants()
		fmt.Println("Seed generated and saved to /tmp/seed")
	case "-gk", "--gen-key":
		generator.GenerateKey()
		fmt.Println("Key generated and saved to /tmp/key")
	}
}

func proccesMode(filePath, keyPath, seedPath string, mode string, command string) {
	var encrypt bool = false
	key, seed := utils.GetKeyFile(keyPath), utils.GetSeedFile(seedPath)

	if command == "-e" || command == "--encrypt" {
		encrypt = true
	}
	if filePath == "" {
		fmt.Println("No file path provided")
		help.PrintHelp()
		os.Exit(1)
	}

	switch mode {
	case "ecb":
		modes.PerformECB(filePath, key, seed, encrypt)
	case "cbc":
		modes.PerformCBC(filePath, key, seed, encrypt)
	case "cfb":
		modes.PerformCFB(filePath, key, seed, encrypt)
	case "ofb":
		modes.PerformOFB(filePath, key, seed, encrypt)
	default:
		fmt.Println("Invalid mode")
		help.PrintHelp()
	}
}
