package modes

import (
	"lea/bitops"
	"lea/core"
	"lea/stream"
	"fmt"
)

func encryptCBC(filePath string, prev *[4]uint32, keySegments []uint32, chunks [4]uint32, size int) {
	after := bitops.MultiXOR32(chunks, *(prev))
	encryptedSlice := core.SelectEncrypt(after, keySegments, size) 
	
	*prev = [4]uint32(encryptedSlice)
	
	if err := stream.WriteBinaryStream(filePath, *prev); err != nil {
		fmt.Printf("Error writing to binary stream: %v\n", err)	
	}
}

func decryptCBC(filePath string, prev *[4]uint32, keySegments []uint32, chunks [4]uint32, size int) {

	encB := core.SelectDecrypt(chunks, keySegments, size)
	text := bitops.MultiXOR32([4]uint32(encB), *prev)
	
	*prev = chunks

	if err := stream.WriteBinaryStream(filePath, text); err != nil {
		fmt.Printf("Error writing to binary stream: %v\n", err)	
	}


}
