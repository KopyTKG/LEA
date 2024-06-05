package logic

import (
	"fmt"
	"lea/encryption"
	"lea/stream"
	"math/rand"
	"time"
)

var key [16]uint32 = [16]uint32{0x0F, 0x1E, 0x2D, 0x3C, 0x4B, 0x5A, 0x69, 0x78, 0x87, 0x96, 0xA5, 0xB4, 0xC3, 0xD2, 0xE1, 0xF0}

func EncryptFile(filePath string) {
	fmt.Println("Encrypting", filePath)
	chunks := stream.BinaryStream(filePath)
	var encChunks []uint32
	keySegments := encryption.Generate(key)
	var blocks [4]uint32

	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		// Every time we fill up a block, encrypt it
		if (i+1)%4 == 0 {
			encryptedBlock := encryption.EncryptBlock(blocks, keySegments)
			encChunks = append(encChunks, encryptedBlock[:]...)
		}
	}

	if len(chunks)%4 != 0 {
		// Fill the remaining slots with some padding, if necessary, or just handle as is
		encryptedBlock := encryption.EncryptBlock(blocks, keySegments)
		encChunks = append(encChunks, encryptedBlock[:]...)
	}

	stream.WriteBinaryStream(filePath, encChunks)
}

func DecryptFile(filePath string) {
	fmt.Println("Decrypting", filePath)
	chunks := stream.BinaryStream(filePath)
	var encChunks []uint32
	keySegments := encryption.Generate(key)
	var blocks [4]uint32

	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		// Every time we fill up a block, encrypt it
		if (i+1)%4 == 0 {
			encryptedBlock := encryption.DecryptBlock(blocks, keySegments)
			encChunks = append(encChunks, encryptedBlock[:]...)
		}
	}

	if len(chunks)%4 != 0 {
		// Fill the remaining slots with some padding, if necessary, or just handle as is
		encryptedBlock := encryption.DecryptBlock(blocks, keySegments)
		encChunks = append(encChunks, encryptedBlock[:]...)
	}

	stream.WriteBinaryStream(filePath, encChunks)
}

func GenerateConstants() {
	rand.Seed(time.Now().UnixNano())
	var seed = []uint32{}

	for i := 0; i < 8; i++ {
		seed = append(seed, uint32(rand.Intn(0xFFFFFFFF)))	
	}

	stream.WriteBinaryStream("/tmp/seed", seed)

}
