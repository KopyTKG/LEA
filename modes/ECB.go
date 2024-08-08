package modes

import (
	"lea/core"
	"lea/stream"
	"fmt"
)

func encryptECB(filePath string, keySegments []uint32, chunks [4]uint32, size int) {
	encryptedBlock := core.SelectEncrypt(chunks, keySegments, size)
	if err := stream.WriteBinaryStream(filePath, [4]uint32(encryptedBlock)); err != nil {
		fmt.Printf("Error writing to binary stream: %v\n", err)
	}
}
func decryptECB(filePath string, keySegments []uint32, chunks [4]uint32, size int) {
	encryptedBlock := core.SelectDecrypt(chunks, keySegments, size)
	if err := stream.WriteBinaryStream(filePath, [4]uint32(encryptedBlock)); err != nil {
		fmt.Printf("Error writing to binary stream: %v\n", err)
	}
}
