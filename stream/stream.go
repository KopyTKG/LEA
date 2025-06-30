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
