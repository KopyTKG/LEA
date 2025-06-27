package main

import (
	"bufio"
	"fmt"
	"lea/help"
	"lea/key"
	"lea/modes"
	"lea/state"
	"lea/stream"
	"lea/utils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kopytkg/golog"
)

func main() {
	/* Initialize Golog */

	timestamp := time.Now().Unix()

	tmpDir := os.TempDir()

	file := fmt.Sprintf("%s/lea-%d.log", tmpDir, timestamp)

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

	validCommandFound, encrypted := false, false
	processCommands(argsList, &validCommandFound, &encrypted)

	if !validCommandFound {
		golog.Error("Invalid command or file path")
		help.PrintHelp()
	}
}

func processArguments(args []string, argsList *utils.List) {
	state.FILEPATH = args[0]

	if len(args) == 1 {

		if state.FILEPATH == "-h" || state.FILEPATH == "--help" {
			help.PrintHelp()
			os.Exit(0)
		}

		if state.FILEPATH == "--version" {
			help.Version()
			os.Exit(0)
		}

		if state.FILEPATH == "keygen" {
			cli := bufio.NewReader(os.Stdin)
			fmt.Print("Please provide key size (128,192,256) [256]: ")
			input, _ := cli.ReadString('\n')
			input = strings.TrimSpace(input)

			size, err := strconv.Atoi(input)
			if err != nil {
				size = 256
			}
			km, err := key.KeyGen(size)

			if err != nil {
				golog.Error(err)
				os.Exit(1)
			}

			if err := key.KeyFileGen(*km, fmt.Sprintf("key%d.lea", size)); err != nil {
				golog.Error(err)
				os.Exit(1)
			}

			os.Exit(0)

		}

	}

	for _, arg := range args[1:] {
		switch {
		case strings.HasPrefix(arg, "--key="):
			armoredkm, err := key.ReadArmoredKey(strings.Split(arg, "=")[1])
			if err != nil {
				golog.Error(err)
				os.Exit(1)
			}

			km, err := key.ParseArmoredKey(armoredkm)
			if err != nil {
				golog.Error(err)
				os.Exit(1)
			}

			state.Key = *km
		case strings.HasPrefix(arg, "--mode="):
			m := strings.Split(arg, "=")[1]
			state.CYPHERMODE = m

		case strings.HasPrefix(arg, "--iter="):
			iterS := strings.Split(arg, "=")[1]

			iterI, err := strconv.Atoi(iterS)

			if err != nil {
				golog.Error(err)
				os.Exit(1)
			}

			state.Iterations = iterI

		// encrypt command must be last
		case arg == "-e" || arg == "-d" || arg == "--encrypt" || arg == "--decrypt":
			argsList.Append(arg)

		case arg == "-r" || arg == "--recursion":
			state.RECURSION = true

		case arg == "-v" || arg == "--verbose":
			state.VERBOSE = true

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
		}
	}
}

func executemode(command string) {
	var encrypt bool = false

	if command == "-e" || command == "--encrypt" {
		encrypt = true
	}

	if state.FILEPATH == "" {
		golog.Error("No file path provided")
		help.PrintHelp()
		os.Exit(1)
	}

	var m map[string]bool

	m = map[string]bool{"ecb": true, "cbc": true, "cfb": true, "ofb": true, "ctr": true}

	if m[state.CYPHERMODE] {
		modes.PerformMode(encrypt)
	} else {
		golog.Error("Invalid state.CYPHERMODE")
		help.PrintHelp()
		os.Exit(1)
	}
}
