# Light Encryption Algorithm (LEA)

![GitHub License](https://img.shields.io/github/license/kopytkg/LEA)


## Overview 

**LEA** (Light Encryption Algorithm) is a command-line tool written in Go for simple file encryption and decryption on Linux. This project was developed as a university assignment for the KI-ZKR course, focusing on implementing the LEA block cipher and providing a practical CLI for file operations.

## :warning: Disclaimer

"***This project is for educational purposes only.***"
* Do not use LEA for protecting sensitive or production data.
* The implementation is not audited or intended for real-world security.
* Use at your own risk.

## Features

* Encrypt and decrypt files using the LEA block cipher
* Supports multiple cipher modes: ECB (default), CBC, CFB, OFB
* Key lengths: 128, 192, or 256 bits (256-bit is default and recommended)
* Optional external key and seed files
* Recursive folder encryption
* Progress display with verbose mode
* Simple, scriptable CLI interface


## Installation

### Building from source 

```bash
git clone https://github.com/kopytkg/LEA.git
cd LEA
mkdir build
make
cd build
```

This will produce binaries for both x86_64 and ARM64 Linux in the build/ directory.

## Usage 

Run `lea --help` or `lea -h` to display the help menu:

```bash
lea --help
```

### Command-Line Options

<h3> Linux </h3>

```bash
Usage: lea [file] [options]

  * -e, --encrypt                 Encrypt the source file
  * -d, --decrypt                 Decrypt the source file
  * -ek, --external-key [file]    Provide an external key file
  * -es, --external-seed [file]   Provide an external seed file
    -h, --help                    Display this help message
    --version                     Display the version of lea
    -v, --verbose                 Display progress screen
    -r, --recursion               Use recursion for folder encryption

  * = required option

Cipher modes:
    --ecb                         Electronic Codebook mode (default)
    --cbc                         Cipher Block Chaining mode
    --cfb                         Cipher Feedback mode
    --ofb                         Output Feedback mode

Key length:
    --128                         128-bit key and seed
    --192                         192-bit key and seed
    --256                         256-bit key and seed (default)

If no options are provided, the file will be encrypted by default.
If no arguments are provided, lea will display this help message.

Report issues at: https://github.com/kopytkg/lea/issues
```

### Example Usage

Encrypt a file with a 256-bit key (default):
```bash
lea myfile.txt -e -ek key.file -es seed.file
```

Decrypt a file:
```bash
lea myfile.txt -d --ek key.file --es seed.file
```

Encrypt all files in a folder recursively:
```bash
lea folder -e -ek key.file -es seed.file -r
```

Show progress during encryption:
```bash
lea folder -e --ek key.file --es seed.file -r -v
```

## Security Notice

* This tool is for learning and demonstration only.
* Do not use for real-world or sensitive data encryption.
* The code is not security-audited.

## License

This project is licensed under the Creative Commons Zero v1.0 Universal License. See [LICENSE](LICENSE.md) for details.
