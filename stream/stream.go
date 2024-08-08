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

// WriteBinaryStream writes a slice of uint32 to a binary file
func WriteBinaryStream(filePath string, data [4]uint32) error {
	// Convert the data to bytes
	bytes := make([]byte, 16)
	for i, val := range data {
		binary.LittleEndian.PutUint32(bytes[i*4:(i+1)*4], val)
	}

	// Open the file for appending
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

func PrepWriteBinaryStream(filePath string, dec [4]uint32) error {
	// Check if the file exists, if not create it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if _, err := os.Create(filePath); err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
	}
	data := make([]byte, 16)
	for i, val := range dec {
		binary.LittleEndian.PutUint32(data[i*4:(i+1)*4], val)
	}
	// Create a temporary file
	tempFileName := filePath + ".tmp"
	tempFile, err := os.Create(tempFileName)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer tempFile.Close()

	// Write the new data to the temporary file
	writer := bufio.NewWriter(tempFile)
	if _, err = writer.Write(data); err != nil {
		return fmt.Errorf("failed to write data to temporary file: %v", err)
	}

	// Open the original file for reading
	originalFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open original file: %v", err)
	}
	defer originalFile.Close()

	// Copy the original content to the temporary file
	buf := make([]byte, 64*1024) // 64KB buffer
	for {
		n, err := originalFile.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read original file: %v", err)
		}
		if n == 0 {
			break
		}
		if _, err = writer.Write(buf[:n]); err != nil {
			return fmt.Errorf("failed to write original content to temporary file: %v", err)
		}
	}

	if err = writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %v", err)
	}

	// Replace the original file with the temporary file
	if err = os.Rename(tempFileName, filePath); err != nil {
		return fmt.Errorf("failed to replace original file: %v", err)
	}

	return nil
}
