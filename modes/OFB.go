package modes

import (
	"fmt"
	"lea/bitops"
	"lea/core"
	"lea/stream"
)


func encryptOFB(filePath string, prev *[4]uint32, keySegments []uint32, chunks [4]uint32, size int) {
	*prev = [4]uint32(core.SelectEncrypt(*prev, keySegments, size))
	encB := bitops.MultiXOR32(*prev, chunks)
	if err := stream.WriteBinaryStream(filePath, encB); err != nil {
		fmt.Printf("Error writing to binary stream: %v\n", err)
	}
}
func decryptOFB(filePath string, prev *[4]uint32, keySegments []uint32, chunks [4]uint32, size int) {
	encryptOFB(filePath, prev, keySegments, chunks, size)
}

