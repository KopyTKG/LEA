package logic

import (
	"fmt"
	"os"
	"lea/stream"
)

var fallbackSeed = [8]uint32{0xC3EFE9DB, 0x44626B02, 0x79E27C8A, 0x78DF30EC, 0x715EA49E, 0xC785DA0A, 0xE04EF22A, 0xE5C40957}

func GetSeedFile(path string) [8]uint32 {
	var seed [8]uint32	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		seed = fallbackSeed
		fmt.Printf("Seed file not found: %v \nUsing fallback\n\n", path)
	} else {
		chunks := stream.BinaryStream(path)
		for i := 0; i < 8; i++ {
		 seed[i] = chunks[i]
		}
	}
	return seed
}
