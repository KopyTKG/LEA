package logic

import (
	"fmt"
	"log"
	"os"
	"lea/stream"
)

var fallback [16]uint32 = [16]uint32{0x0F, 0x1E, 0x2D, 0x3C, 0x4B, 0x5A, 0x69, 0x78, 0x87, 0x96, 0xA5, 0xB4, 0xC3, 0xD2, 0xE1, 0xF0}

func GetExternalKey(path string) [16]uint32 {
	var key [16]uint32	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Key file not found: %v\n", path)
	} else {
		// read seed file and parse it to constval
		chunks := stream.BinaryStream(path)
		for i := 0; i < 16; i++ {
		 key[i] = chunks[i]
		}
	}
	return key
}

func GetInternalKey() [16]uint32 {
	var key [16]uint32	
	if _, err := os.Stat("/tmp/key"); os.IsNotExist(err) {
		key = fallback
		fmt.Println("Key file not found, using fallback key")
	} else {
		// read seed file for /tmp/seed and parse it to constval
		chunks := stream.BinaryStream("/tmp/key")
		for i := 0; i < 16; i++ {
		 key[i] = chunks[i]
		}
	}
	return key
}
