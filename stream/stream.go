package stream

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// BinaryStream reads a binary file and returns its contents as a slice of uint32
func BinaryChunkStream(path string) []uint32 {
	file, err := os.Open(path)
	fmt.Println(file)
	check(err)
	defer file.Close()

	reader := bufio.NewReader(file)

	var chunks []uint32

	buf := make([]byte, 4)
	for {
		n, err := io.ReadFull(reader, buf)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				if n > 0 {
					var paddedBuf [4]byte
					copy(paddedBuf[:], buf[:n])
					chunks = append(chunks, binary.LittleEndian.Uint32(paddedBuf[:]))
				}
				break
			}
		}
		chunk := binary.LittleEndian.Uint32(buf)
		chunks = append(chunks, chunk)
	}
	for len(chunks)%4 != 0 {
		chunks = append(chunks, 0)
	}
	return chunks
}

// WriteBinaryStream appends a slice of uint32 to a binary file
func WriteBinaryStream(filePath string, data []byte) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer file.Close()

	// Convert the data to bytes
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

// WriteCMAC writes a 16-byte CMAC (as 4 uint32 words) at offset 32 in the file in little-endian order.
// It overwrites bytes 32–47.
func WriteCMAC(filePath string, data []uint32) error {
	if len(data) != 4 {
		return fmt.Errorf("expected CMAC as 4 uint32 words (16 bytes), got %d", len(data))
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer file.Close()

	// Seek to offset 32
	_, err = file.Seek(32, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek in file: %v", err)
	}

	buf := make([]byte, 16)
	for i, v := range data {
		binary.LittleEndian.PutUint32(buf[i*4:], v)
	}
	_, err = file.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write CMAC: %v", err)
	}
	return nil
}
