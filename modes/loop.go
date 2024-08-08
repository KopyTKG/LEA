package modes

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"lea/fingerprint"
	"lea/schedule"
	"log"
	"os"
	"strings"
)

func PerformMode(mode, filePath string, bKey, bSeed []byte, encrypt bool, keySize int) {
	kChunks := fingerprint.LoadSource(bKey)
	sChunks := fingerprint.LoadSource(bSeed)
	key := fingerprint.SelectPrint(kChunks, keySize)
	seed := fingerprint.SelectPrint(sChunks, keySize)
	rk := schedule.KeySchedule(keySize, key, seed)
	tmpFilePath := filePath + ".tmp"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var chunks []uint32
	var prev [4]uint32
	var IV [4]uint32

	if mode != "ecb" {
		 cli := bufio.NewReader(os.Stdin)
		fmt.Print("Please provide an IV: ")
		input, _ := cli.ReadString('\n')
		input = strings.TrimSpace(input)

		// Convert the input to a slice of uint32
		IV = fingerprint.Fingerprint128(fingerprint.LoadSource([]byte(input)))
	}

	if encrypt || mode == "ofb" {
		prev = IV
		buf := make([]byte, 4)

		for {
			n, err := io.ReadFull(reader, buf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				if n > 0 {
					var paddedBuf [4]byte
					copy(paddedBuf[:], buf[:n])
					chunks = append(chunks, binary.LittleEndian.Uint32(paddedBuf[:]))
				}
				break
			}
			if err != nil {
				log.Fatalf("Error reading file: %v", err)
			}

			chunk := binary.LittleEndian.Uint32(buf)
			chunks = append(chunks, chunk)

			if len(chunks) == 4 {
				loop(mode, tmpFilePath, rk, [4]uint32(chunks), &prev, encrypt, false, keySize, IV)
				chunks = []uint32{}
			}
		}

		// Handle any remaining bytes
		for len(chunks)%4 != 0 {
			chunks = append(chunks, 0)
		}
		if len(chunks) > 0 {
			loop(mode, tmpFilePath, rk, [4]uint32(chunks), &prev, encrypt, true, keySize, IV)
		}
	} else {
		// Get the file size
		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatalf("Error getting file info: %v", err)
		}
		fileSize := fileInfo.Size()

		// Define chunk size (4 bytes for uint32)
		const chunkSize = 4
		var chunks []uint32

		// Read file in chunks from the end
		for offset := fileSize - chunkSize; offset >= 0; offset -= chunkSize {
			buf := make([]byte, chunkSize)
			_, err := file.ReadAt(buf, offset)
			if err != nil && err != io.EOF {
				log.Fatalf("Error reading file at offset %d: %v", offset, err)
			}

			// Convert bytes to uint32
			chunk := binary.LittleEndian.Uint32(buf)
			chunks = append([]uint32{chunk}, chunks...) // Prepend to maintain reverse order

			// Process chunks in groups of 4
			if len(chunks) == 4 {
				lastChunk := offset == 0
				loop(mode, tmpFilePath, rk, [4]uint32(chunks), &prev, encrypt, lastChunk, keySize, IV)
				chunks = []uint32{} // Clear chunks for the next group
			}
		}
	}
 	cleanup(tmpFilePath, filePath)
}

func cleanup(tmp, filePath string) {
	if err := os.Remove(filePath); err != nil {
		log.Fatalf("Error removing original file: %v", err)
	}
	if err := os.Rename(tmp, filePath); err != nil {
		log.Fatalf("Error renaming temporary file: %v", err)
	}
}

func loop(mode, filePath string, rk []uint32, chunks [4]uint32, prev *[4]uint32, encrypt, last bool, keySize int, IV [4]uint32) {
	switch mode {
	default:
		log.Fatalln("No mode selected")
		os.Exit(1)

	case "ecb":
		if encrypt {
			encryptECB(filePath, rk, chunks, keySize)
		} else {
			decryptECB(filePath, rk, chunks, keySize)
		}

	case "cbc":
		if encrypt {
			encryptCBC(filePath, prev, rk, chunks, keySize)
		} else {
			decryptCBC(filePath, prev, rk, chunks, keySize, last, IV)
		}

	case "cfb":
		if encrypt {
			encryptCFB(filePath, prev, rk, chunks, keySize)
		} else {
			decryptCFB(filePath, prev, rk, chunks, keySize, last, IV)
		}

	case "ofb":
		if encrypt {
			encryptOFB(filePath, prev, rk, chunks, keySize)
		} else {
			encryptOFB(filePath, prev, rk, chunks, keySize)
		}
	}
}
