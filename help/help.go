package help

import "fmt"

func PrintHelp() {
	helpText := `LEA Encryption Tool
Usage: lea [file] ?[options]

CORE OPERATIONS:
  * -e, --encrypt                Encrypt target file/directory
  * -d, --decrypt                Decrypt target file/directory
    --key=[file]                 Use pre-generated key file (required)
    lea keygen                   Generate new LEA keyfile

KEY MANAGEMENT:
   Generated keys are saved as:
   - key128.lea (128-bit)
   - key192.lea (192-bit)
   - key256.lea (256-bit, default)

   File Format:
   -----BEGIN LEA KEY-----
   Version: 1
   Key: A1B2C3D4:E5F67890:...
   Key Digest: sha256...
   Seed: A1B2C3D4:E5F67890:...
   Seed Digest: sha256...
   -----END LEA KEY-----

ADVANCED:
    --mode=[ecb|cbc|cfb|ofb|ctr]     Encryption mode (default: cbc)
    -r, --recursive              Process directories recursively
    -v, --verbose                Show detailed progress
    --iter=N                     Overwrite N times (default: 3)

DEBUG:
    --version			 Show program version


Examples:
  Encrypt:
    $ lea file.txt -e --key=key256.lea

  Decrypt with raw LEA-generated key:
    $ lea file.txt -d --key=key.256.lea


LEGAL DISCLAIMER:
  This tool is for AUTHORIZED penetration testing only.
  By using this software, you confirm you have permission
  to test the target systems.

Report issues to: <https://github.com/kopytkg/lea/issues>
`
	fmt.Print(helpText)
}
