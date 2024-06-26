package modes

import (
	"fmt"
	"lea/encryption"
	"lea/stream"
	"lea/utils"
	"lea/errors"
	"lea/bitops"
	"os"
	"bufio"
)

func PerformOFB(filePath string, key [16]uint32, seed [8]uint32, encrypt bool) {
	chunks := stream.BinaryChunkStream(filePath)
	keySegments := encryption.Generate(key, seed)
	var blocks [4]uint32
	
	if encrypt {
		encryptOFB(filePath, blocks, keySegments, chunks)
	} else {
		decryptOFB(filePath, blocks, keySegments, chunks)
	}

}

func encryptOFB(filePath string, blocks [4]uint32, keySegments [144]uint32, chunks []uint32) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Encrypting", filePath)	
	var encChunks []uint32
	// ask user for IV through stdin
	fmt.Println("Please provide an IV: ")
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]
	
	// conver the input to a slice of uint32
	var IV [4]uint32 = utils.FingerPrint128([]byte(input))
	var prev [4]uint32 = IV
	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		if (i+1)%4 == 0 {
			prev = encryption.EncryptBlock(prev, keySegments)
			encryptedBlock := bitops.MultiXOR64(prev, blocks)
			encChunks = append(encChunks, encryptedBlock[:]...)
		}
	}

	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)
}

func decryptOFB(filePath string, blocks [4]uint32, keySegments [144]uint32, chunks []uint32) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Encrypting", filePath)	
	var encChunks []uint32
	// ask user for IV through stdin
	fmt.Println("Please provide an IV: ")
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]
	
	

	// conver the input to a slice of uint32
	var IV [4]uint32 = utils.FingerPrint128([]byte(input))
	var prev [4]uint32 = IV	
	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		if (i+1)%4 == 0 {
			prev = encryption.EncryptBlock(prev, keySegments)
			encryptedBlock := bitops.MultiXOR64(prev, blocks)
			encChunks = append(encChunks, encryptedBlock[:]...)
		}
	}

	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)
}
