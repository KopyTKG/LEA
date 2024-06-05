package help

import "fmt"

func PrintHelp() {
    helpText := `
----------------------------------------------------
                       LEA Encryption
----------------------------------------------------

Options:
    *file*           Source file path with extension
    -E               Encrypt the source file
    -D               Decrypt the source file

Secrets options:
    By default, the program uses a preset key and seed. You can generate
    new keys and seeds with the following options:
    
    --gen-key        Generate an internal key and save it to /tmp/key
    --gen-seed       Generate an internal seed and save it to /tmp/seed


    Also, you can provide your own key and seed by adding the following flags:

    --external-key *file*.key   	Provide an external key. Needs to be 512 bytes or longer 
    					(anything longer will be truncated to 512 bytes)

    --external-seed *file*.seed 	Provide an external seed. Needs to be 256 bytes or longer 
    					(anything longer will be truncated to 256 bytes)

Usage:
    lea *file*

Default key and seed:
    lea *file* -E
    lea *file* -D

With key and seed generation:
    lea *file* --gen-key --gen-seed -E
    lea *file* -D

`
    fmt.Print(helpText)
}
