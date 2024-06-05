package main

import (
	"fmt"
	"lea/help"
	"lea/logic"
	"lea/generator"
	"lea/list"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		help.PrintHelp()
		return
	}

	filePath := ""
	validCommandFound := false
	encrypted := false
	
	argsList := list.List{}
	
	for _, arg := range args {
		if arg == "-E" || arg == "-D" {
			argsList.Append(arg)
		}
		
		if arg == "--gen-seed" || arg == "--gen-key" {
			if arg == "--gen-seed" {
			  index := argsList.IndexOf("--external-seed")
			  if index == -1 {
				argsList.Prepend(arg)
			  }
			}
			if arg == "--gen-key" {
			  index := argsList.IndexOf("--external-key")
			  if index == -1 {
				argsList.Prepend(arg)
			  }
			}
		}

		if arg == "--external-key" || arg == "--external-seed" {
			if argsList.Length() > 0 {
				if argsList.IndexOf("--gen-key") == -1 && arg == "--external-key" {
					argsList.Prepend(arg)
				}
				if argsList.IndexOf("--gen-seed") == -1 && arg == "--external-seed" {
					argsList.Prepend(arg)
				}
			} else {
				argsList.Prepend(arg)
			}
		} 
		
		if strings.Contains(arg, ".key") || strings.Contains(arg, ".seed") {
			if argsList.Length() > 0 {
				if strings.Contains(arg, ".key") {
					index := argsList.IndexOf("--external-key")
					if index != -1 {
						argsList.Insert(index+1, arg)
					}
				}
				if strings.Contains(arg, ".seed") {
					index := argsList.IndexOf("--external-seed")
					if index != -1 {
						argsList.Insert(index+1, arg)
					} 
				}
			}
		} else {
		if strings.Contains(arg, ".") {
			if argsList.Length() > 0 {
			if argsList.Get(0) == "-E" || argsList.Get(0) == "-D" {
				argsList.Prepend(arg)
			} else {
				argsList.Insert(1, arg)
			}
			} else {
				argsList.Append(arg)
			}
		}
		}
	}

	key := [16]uint32{}
	seed := [8]uint32{}

	for _, arg := range argsList.Elements {
		switch {
		case arg == "-E" || arg == "-D" || arg == "--gen-seed" || arg == "--gen-key":
			validCommandFound = true
			switch arg {
			case "-E":
				encrypted = true
				if filePath != "" {
					if seed == [8]uint32{} {
					}
					if key == [16]uint32{} {
						key = logic.GetInternalKey()
					}
					logic.EncryptFile(filePath, key)
				} else {
					fmt.Println("No file path provided for encryption.")
				}
			case "-D":
				encrypted = true
				if filePath != "" {
					if seed == [8]uint32{} {
					}
					if key == [16]uint32{} {
						key = logic.GetInternalKey()
					}
					logic.DecryptFile(filePath, key)
				} else {
					fmt.Println("No file path provided for decryption.")
				}
			case "--gen-seed":
				generator.GenerateConstants()
				fmt.Println("Seed generated and saved to /tmp/seed")

			case "--gen-key":
				generator.GenerateKey()
				fmt.Println("Key generated and saved to /tmp/key")
			}

		default:
			if strings.Contains(arg, ".key") || strings.Contains(arg, ".seed") {
				if strings.Contains(arg, ".key") {
					key = logic.GetExternalKey(arg)
					fmt.Println("Using external key")
				} else {
					/* seed = logic.GetExternalSeed(arg) */
				}
			} else if strings.Contains(arg, ".") {
				filePath = arg
				validCommandFound = true
			}
		}
	}

	if !validCommandFound {
		fmt.Println("Invalid command or file path")
		help.PrintHelp()
	}

	if !encrypted && filePath != "" {
		if key == [16]uint32{} {
			key = logic.GetInternalKey()
		}
		logic.EncryptFile(filePath, key)
	}
}
