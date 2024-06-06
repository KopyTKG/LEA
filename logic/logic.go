package logic

import (
	"fmt"
	"lea/encryption"
	"lea/stream"
)

func EncryptFile(filePath string, key [16]uint32, seed [8]uint32) {
	fmt.Println("Encrypting", filePath)
	chunks := stream.BinaryStream(filePath)
	var encChunks []uint32
	keySegments := encryption.Generate(key, seed)
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

func DecryptFile(filePath string, key [16]uint32, seed [8]uint32) {
	fmt.Println("Decrypting", filePath)
	chunks := stream.BinaryStream(filePath)
	var encChunks []uint32
	keySegments := encryption.Generate(key, seed)
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

