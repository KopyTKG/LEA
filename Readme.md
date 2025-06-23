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
* Supports multiple cipher modes: ECB, CBC (default), CFB, OFB, CTR
* Key lengths: 128, 192, or 256 bits (256-bit is default and recommended)
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
LEA Encryption Tool
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
```

### Example Usage

Generate a Key file:
```bash
lea keygen

Please provide key size (128,192,256) [256]: 
17:53:57 INFO Successfully generated key file: key256.lea
``` 

Encrypt a file with a 256-bit key (default):
```bash
lea myfile.txt -e --key=key256.lea
```

Decrypt a file:
```bash
lea myfile.txt -d --key=key256.lea
```

Encrypt all files in a folder recursively:
```bash
lea folder -e -r --key=key256.lea
```

Show progress during encryption:
```bash
lea folder -e -r -v --key=key256.lea
```

## Security Notice

* This tool is for learning and demonstration only.
* Do not use for real-world or sensitive data encryption.
* The code is not security-audited.

## License

This project is licensed under the Creative Commons Zero v1.0 Universal License. See [LICENSE](LICENSE.md) for details.
