package modes

import (
	"fmt"
	"lea/encryption"
	"lea/stream"
)

func PerformECB(filePath string, key [16]uint32, seed [8]uint32, encrypt bool) {
	chunks := stream.BinaryChunkStream(filePath)
	keySegments := encryption.Generate(key, seed)
	var blocks [4]uint32
	
	if encrypt {
		encryptFile(filePath, blocks, keySegments, chunks)
	} else {
		decryptFile(filePath, blocks, keySegments, chunks)
	}

}


func encryptFile(filePath string, blocks [4]uint32, keySegments [144]uint32, chunks []uint32) {
	fmt.Println("Encrypting", filePath)
	var encChunks []uint32
	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		if (i+1)%4 == 0 {
			encryptedBlock := encryption.EncryptBlock(blocks, keySegments)
			encChunks = append(encChunks, encryptedBlock[:]...)
		}
	}
	if len(chunks)%4 != 0 {
		encryptedBlock := encryption.EncryptBlock(blocks, keySegments)
		encChunks = append(encChunks, encryptedBlock[:]...)
	}
	stream.WriteBinaryStream(filePath, encChunks)
}

func decryptFile(filePath string, blocks [4]uint32, keySegments [144]uint32, chunks []uint32) {
	fmt.Println("Decrypting", filePath)
	var encChunks []uint32
	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		if (i+1)%4 == 0 {
			encryptedBlock := encryption.DecryptBlock(blocks, keySegments)
			encChunks = append(encChunks, encryptedBlock[:]...)
		}
	}
	if len(chunks)%4 != 0 {
		encryptedBlock := encryption.DecryptBlock(blocks, keySegments)
		encChunks = append(encChunks, encryptedBlock[:]...)
	}
	stream.WriteBinaryStream(filePath, encChunks)
}

