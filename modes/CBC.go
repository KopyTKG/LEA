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

func PerformCBC(filePath string, key [16]uint32, seed [8]uint32, encrypt bool) {
	chunks := stream.BinaryChunkStream(filePath)
	keySegments := encryption.Generate(key, seed)
	var blocks [4]uint32
	
	if encrypt {
		encryptCBC(filePath, blocks, keySegments, chunks)
	} else {
		decryptCBC(filePath, blocks, keySegments, chunks)
	}

}

func encryptCBC(filePath string, blocks [4]uint32, keySegments [144]uint32, chunks []uint32) {
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
			after :=  bitops.MultiXOR64(blocks, prev)
			encryptedBlock := encryption.EncryptBlock(after, keySegments)
			encChunks = append(encChunks, encryptedBlock[:]...)
			prev = encryptedBlock
		}
	}
	
	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)
}



func decryptCBC(filePath string, blocks [4]uint32, keySegments [144]uint32, chunks []uint32) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Encrypting", filePath)	
	var encChunks []uint32
	// ask user for IV through stdin
	fmt.Println("Please provide an IV: ")
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]
	
	

	// conver the input to a slice of uint32
	var IV [4]uint32 = utils.FingerPrint128([]byte(input))
	
	var stackedBlocks utils.Stack = utils.Stack{}
	var flushStack utils.Stack = utils.Stack{}

	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		if (i+1)%4 == 0 {
			stackedBlocks.Append(blocks)		
		}
	}
	for stackedBlocks.Length() > 0 {
		currentBlock := stackedBlocks.Pop()
		decryptedBlock := encryption.DecryptBlock(currentBlock, keySegments)
		next := IV
		if stackedBlocks.Length() != 0 {next= stackedBlocks.Peek()}
		base := bitops.MultiXOR64(decryptedBlock, next)
		flushStack.Append(base)
	}
	for flushStack.Length() > 0 {
		el := flushStack.Pop()
		encChunks = append(encChunks, el[:]...)
	}

	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)
}

