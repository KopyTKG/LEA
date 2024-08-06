package help

import "fmt"

func PrintHelp() {
	helpText := `Usage: lea [file] ?[options]

  * -e, --encrypt               	Encrypt the source file
  * -d, --decrypt               	Decrypt the source file
  * -ek, --external-key [file]   	Provide an external any key file. 
  * -es, --external-seed [file] 	Provide an external any seed file. 
    -h, --help                  	Display this help message
    -v, --version               	Display the version of lea

current version supports following modes:
    --ecb                       	Electronic Codebook mode (default)
    --cbc                       	Cipher Block Chaining mode
    --cfb                       	Cipher Feedback mode
    --ofb                       	Output Feedback mod

If no options are provided, the file will be encrypted by default. 
If nothing is provided at all lea will display help.

Any errors please report to: <https://github.com/kopytkg/lea/issues>
`
	fmt.Print(helpText)
}
