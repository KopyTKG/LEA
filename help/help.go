package help

import "fmt"

func PrintHelp() {
    helpText := `Usage: lea [file] ?[options]

    -e, --encrypt               	Encrypt the source file
    -d, --decrypt               	Decrypt the source file
    -gk, --gen-key        		Generate an internal key and save it to /tmp/key
    -gs, --gen-seed       		Generate an internal seed and save it to /tmp/seed
    -ek, --external-key [file].key   	Provide an external key. Needs to be 512 bytes or longer 
    					(anything longer will be truncated to 512 bytes)
    -es, --external-seed [file].seed 	Provide an external seed. Needs to be 256 bytes or longer 
    					(anything longer will be truncated to 256 bytes)
    -h, --help                  	Display this help message

If no options are provided, the file will be encrypted by default. 
If nothing is provided at all lea will display help.

Any errors please report to: <https://github.com/kopytkg/lea/issues>
`
    fmt.Print(helpText)
}
