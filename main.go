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
	/* Setting version */

	help.VERSION = "v3.1.0"

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
	if err := processArguments(args, &argsList); err != nil {
		golog.Error(err)
		os.Exit(1)
	}

	if err := validateArguments(); err != nil {
		golog.Error(err)
		os.Exit(1)
	}

	if err := executemode(); err != nil {
		golog.Error(err)
		os.Exit(1)
	}
}

func processArguments(args []string, argsList *utils.List) error {
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
				return err
			}

			if err := key.KeyFileGen(*km, fmt.Sprintf("key%d.lea", size)); err != nil {
				return err
			}

			os.Exit(0)

		}

	}

	for _, arg := range args[1:] {
		switch {
		case strings.HasPrefix(arg, "--key="):
			armoredkm, err := key.ReadArmoredKey(strings.Split(arg, "=")[1])
			if err != nil {
				return err
			}

			km, err := key.ParseArmoredKey(armoredkm)
			if err != nil {
				return err
			}

			state.Key = *km
		case strings.HasPrefix(arg, "--mode="):
			m := strings.Split(arg, "=")[1]
			state.CYPHERMODE = m

		case strings.HasPrefix(arg, "--iter="):
			iterS := strings.Split(arg, "=")[1]

			iterI, err := strconv.Atoi(iterS)

			if err != nil {
				return err
			}

			state.Iterations = iterI

		case arg == "-e" || arg == "--encrypt":
			state.ENCRYPT = true

		case arg == "-d" || arg == "--decrypt":
			state.ENCRYPT = false

		case arg == "-r" || arg == "--recursion":
			state.RECURSION = true

		case arg == "-v" || arg == "--verbose":
			state.VERBOSE = true

		default:
			return fmt.Errorf("Unknowed switch found (%s) run \"lea -h\"", arg)
		}
	}
	return nil
}

func validateArguments() error {
	if state.RECURSION {
		b, err := stream.IsFolder(state.FILEPATH)

		if err != nil {
			return err
		}

		if !b {
			return fmt.Errorf("Path (%s) must be folder for recursion operation\n", state.FILEPATH)
		}
	}

	if state.CYPHERMODE == "" {
		return fmt.Errorf("No cipher mode provided")
	}

	if state.FILEPATH == "" {
		return fmt.Errorf("No file path provided")
	}

	return nil
}

func executemode() error {
	switch state.CYPHERMODE {
	case "ecb":
		modes.SelectedMode = &modes.ECB{}
	case "cbc":
		modes.SelectedMode = &modes.CBC{}
	case "cfb":
		modes.SelectedMode = &modes.CFB{}
	case "ofb":
		modes.SelectedMode = &modes.OFB{}
	case "ctr":
		modes.SelectedMode = &modes.CTR{}
	default:
		return fmt.Errorf("Unsupported cipher mode: %s", state.CYPHERMODE)
	}

	if state.ENCRYPT {
		modes.CryptMethod = modes.SelectedMode.Encrypt
	} else {
		modes.CryptMethod = modes.SelectedMode.Decrypt
	}

	modes.PerformMode()
	return nil
}
