package utils

import (
	"fmt"
	"os"
	"lea/stream"
)

var fallbackKey [16]uint32 = [16]uint32{0x0F, 0x1E, 0x2D, 0x3C, 0x4B, 0x5A, 0x69, 0x78, 0x87, 0x96, 0xA5, 0xB4, 0xC3, 0xD2, 0xE1, 0xF0}

func GetKeyFile(path string) [16]uint32 {
	var key [16]uint32	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		key = fallbackKey
		fmt.Printf("Key file not found: %v \nUsing fallback\n\n", path)
	} else {
		// uint32 array of chunks
		chunks := stream.BinaryStream(path)
		key = FingerPrint512(chunks)
	}
	return key
}
