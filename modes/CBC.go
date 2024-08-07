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
	
	if err := stream.WriteBinaryStreamv2(filePath, *prev); err != nil {
		fmt.Printf("Error writing to binary stream: %v\n", err)	
	}
}

func decryptCBC(filePath string, prev *[4]uint32, keySegments []uint32, chunks [4]uint32, size int, last bool, IV [4]uint32) {
	current := core.SelectDecrypt(chunks, keySegments, size)
	
	if *prev != [4]uint32{}{
		base := bitops.MultiXOR32(*prev, chunks)
		if err := stream.PrepWriteBinaryStream(filePath, base); err != nil {
			fmt.Printf("Error prepending to binary stream: %v\n", err)
		}
	}
	
	*prev = [4]uint32(current)

	if last {
		base := bitops.MultiXOR32(*prev, IV)
		if err := stream.PrepWriteBinaryStream(filePath, base); err != nil {
			fmt.Printf("Error prepending to binary stream: %v\n", err)
		}
	}
}
