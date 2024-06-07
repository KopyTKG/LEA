package stream

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// BinaryStream reads a binary file and returns its contents as bytes
func BinaryStream(path string) []byte {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	reader := bufio.NewReader(file)
	var chunks []byte
	for {
		n, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		chunks = append(chunks, n)
	}
	return chunks
}

// BinaryStream reads a binary file and returns its contents as a slice of uint32
func BinaryChunkStream(path string) []uint32 {
	file, err := os.Open(path)
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
	return chunks
}

// WriteBinaryStream writes a slice of uint32 to a binary file
func WriteBinaryStream(fileName string, data []uint32) {
	bytes := make([]byte, len(data)*4)
	for i, val := range data {
		binary.LittleEndian.PutUint32(bytes[i*4:(i+1)*4], val)
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(bytes)
	if err != nil {
		log.Fatalf("Failed to write bytes: %v\n", err)
		return
	}

	if err := writer.Flush(); err != nil {
		log.Fatalf("Failed to flush writer: %v\n", err)
	}

}
