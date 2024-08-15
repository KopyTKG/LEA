package help

import "fmt"

func PrintHelp() {
	helpText := `Usage: lea [file] ?[options]

  * -e, --encrypt               	Encrypt the source file
  * -d, --decrypt               	Decrypt the source file
  * -ek, --external-key [file]   	Provide an external any key file. 
  * -es, --external-seed [file] 	Provide an external any seed file. 
    -h, --help                  	Display this help message
    --version               		Display the version of lea
    -v, --verbose			Display progress screen

  * marks required switch

current version supports following modes:
    --ecb                       	Electronic Codebook mode (default)
    --cbc                       	Cipher Block Chaining mode
    --cfb                       	Cipher Feedback mode
    --ofb                       	Output Feedback mod

key length:
    --128				Basic 128bit key and seed
    --192				
    --256				Recommended (default)

If no options are provided, the file will be encrypted by default. 
If nothing is provided at all lea will display help.

Any errors please report to: <https://github.com/kopytkg/lea/issues>

usage
$ lea [file] --ek [file] --es [file] -e/-d [optional mode] [optional lenght]
`
	fmt.Print(helpText)
}
