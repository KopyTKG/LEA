package help

import "fmt"

func PrintHelp() {
    helpText := `Usage: lea [file] ?[options]

    -e, --encrypt               	Encrypt the source file
    -d, --decrypt               	Decrypt the source file
    -gk, --gen-key        		Generate an internal key and save it to /tmp/key
    -gs, --gen-seed       		Generate an internal seed and save it to /tmp/seed
    -ek, --external-key [file].key   	Provide an external key. SHA3-512 will do the rest. 
    -es, --external-seed [file].seed 	Provide an external seed. SHA3-256 will do the rest. 
    -h, --help                  	Display this help message
    -v, --version               	Display the version of lea

current version supports following modes:
    --ecb                       	Electronic Codebook mode (default)
    --cbc                       	Cipher Block Chaining mode

If no options are provided, the file will be encrypted by default. 
If nothing is provided at all lea will display help.

Any errors please report to: <https://github.com/kopytkg/lea/issues>
`
    fmt.Print(helpText)
}
