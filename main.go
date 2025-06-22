package main

import (
	"fmt"
	"lea/help"
	"lea/modes"
	"lea/state"
	"lea/stream"
	"lea/utils"
	"os"
	"time"

	"github.com/kopytkg/golog"
)

func main() {
	/* Initialize Golog */

	timestamp := time.Now().Unix()

	file := fmt.Sprintf("lea-%d.log", timestamp)

	err := golog.EnableLogFile(file)
	if err != nil {
		panic(err)
	}

	golog.LOGLEVEL = golog.INFO
	golog.CLILOG = golog.ENABLED

	/* --------------- */

	args := os.Args[1:]
	if len(args) < 1 {
		help.PrintHelp()
		return
	}

	// preset file paths
	argsList := utils.List{}

	// processing loop for args
	processArguments(args, &argsList)

	if argsList.Length() == 0 && state.FILEPATH != "" {
		argsList.Append("-e")
	}

	if state.KEYPATH == "" || state.SEEDPATH == "" {
		sw := ""
		if state.KEYPATH == "" {
			sw = "-ek"
		} else {
			sw = "-es"
		}

		golog.Errorf("Missing required switch (%s) run \"lea -h\"", sw)
	}

	validCommandFound, encrypted := false, false
	processCommands(argsList, &validCommandFound, &encrypted)

	if !validCommandFound {
		golog.Error("Invalid command or file path")
		help.PrintHelp()
	}
}

func processArguments(args []string, argsList *utils.List) {
	prev := ""
	state.FILEPATH = args[0]

	if len(args) == 1 {

		if state.FILEPATH == "-h" || state.FILEPATH == "--help" {
			help.PrintHelp()
			os.Exit(1)
		}

		if state.FILEPATH == "--version" {
			help.Version()
			os.Exit(1)
		}

	}

	for _, arg := range args[1:] {
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
			state.KEYPATH = arg
		case prev == "-es" || prev == "--external-seed":
			prev = ""
			state.SEEDPATH = arg

		case arg == "-r" || arg == "--recursion":
			state.RECURSION = true

		case arg == "-v" || arg == "--verbose":
			state.VERBOSE = true

		// Cypher state.CYPHERMODEs
		case arg == "--ecb":
			state.CYPHERMODE = "ecb"
		case arg == "--cbc":
			state.CYPHERMODE = "cbc"
		case arg == "--cfb":
			state.CYPHERMODE = "cfb"
		case arg == "--ofb":
			state.CYPHERMODE = "ofb"

		// Key lenght
		case arg == "--128":
			state.KEYLENGTH = 128
		case arg == "--192":
			state.KEYLENGTH = 192
		case arg == "--256":
			state.KEYLENGTH = 256

		default:
			golog.Errorf("Unknowed switch found (%s) run \"lea -h\"", arg)
			os.Exit(1)
		}
	}
}

func processCommands(argsList utils.List, validCommandFound, encrypted *bool) {

	if state.RECURSION {
		b, err := stream.IsFolder(state.FILEPATH)

		if err != nil {
			panic(err)
		}

		if !b {
			golog.Errorf("Path (%s) must be folder for recursion operation\n", state.FILEPATH)
		}
	}

	for _, arg := range argsList.Elements {
		switch arg {
		case "-e", "-d", "--encrypt", "--decrypt":
			*validCommandFound = true
			*encrypted = true
			executemode(arg)

		case "--external-key", "--external-seed", "-ek", "-es":
			*validCommandFound = true
		}
	}
}

func executemode(command string) {
	var encrypt bool = false

	byteKEY, err := stream.GetFile(state.KEYPATH)
	if err != nil {
		golog.Error(err)
	}
	state.ByteKEY = byteKEY

	byteSEED, err := stream.GetFile(state.SEEDPATH)
	if err != nil {
		golog.Error(err)
	}
	state.ByteSEED = byteSEED

	if command == "-e" || command == "--encrypt" {
		encrypt = true
	}

	if state.FILEPATH == "" {
		golog.Error("No file path provided")
		help.PrintHelp()
		os.Exit(1)
	}

	if state.CYPHERMODE == "ecb" || state.CYPHERMODE == "cbc" || state.CYPHERMODE == "cfb" || state.CYPHERMODE == "ofb" {
		modes.PerformMode(encrypt)
	} else {
		golog.Error("Invalid state.CYPHERMODE")
		help.PrintHelp()
		os.Exit(1)
	}
}
