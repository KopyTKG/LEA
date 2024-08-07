package modes

import (
	"fmt"
	"lea/core"
	"lea/errors"
	"lea/fingerprint"
	"lea/schedule"
	"lea/stream"
)

func PerformECB(filePath string, bKey []byte, bSeed []byte, encrypt bool, keySize int) {
	chunks := stream.BinaryChunkStream(filePath)
	
	kChunks := fingerprint.LoadSource(bKey)
	sChunks := fingerprint.LoadSource(bSeed)
	
        key := fingerprint.SelectPrint(kChunks, keySize)
        seed := fingerprint.SelectPrint(sChunks, keySize)

	rk := schedule.KeySchedule(keySize, key, seed)
	var blocks [4]uint32

	if encrypt {
		encryptECB(filePath, blocks, rk, chunks, keySize)
	} else {
		decryptECB(filePath, blocks, rk, chunks, keySize)
	}

}

func encryptECB(filePath string, blocks [4]uint32, keySegments []uint32, chunks []uint32, size int) {
	fmt.Println("Encrypting", filePath)
	var encChunks []uint32
	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		if (i+1)%4 == 0 {
			encryptedBlock := core.SelectEncrypt(blocks, keySegments, size)
			encChunks = append(encChunks, encryptedBlock[:]...)
		}
	}
	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)
}

func decryptECB(filePath string, blocks [4]uint32, keySegments []uint32, chunks []uint32, size int) {
	fmt.Println("Decrypting", filePath)
	var encChunks []uint32
	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		if (i+1)%4 == 0 {
			encryptedBlock := core.SelectDecrypt(blocks, keySegments, size)
			encChunks = append(encChunks, encryptedBlock[:]...)
		}
	}
	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)
}
