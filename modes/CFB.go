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


func PerformCFB(filePath string, key [16]uint32, seed [8]uint32, encrypt bool) {
	chunks := stream.BinaryChunkStream(filePath)
	keySegments := encryption.Generate(key, seed)
	var blocks [4]uint32
	
	if encrypt {
		encryptCFB(filePath, blocks, keySegments, chunks)
	} else {
		decryptCFB(filePath, blocks, keySegments, chunks)
	}

}


func encryptCFB(filePath string, blocks [4]uint32, keySegments [144]uint32, chunks []uint32) {
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
			prev = bitops.MultiXOR64(prev, blocks)
			encChunks = append(encChunks, prev[:]...)
		}
	}
	
	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)


}

func decryptCFB(filePath string, blocks [4]uint32, keySegments [144]uint32, chunks []uint32) {
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
		// current block X.Stext
		currentBlock := stackedBlocks.Pop()
		// decyphered X-1.Stext
		prevBlock := IV
		if stackedBlocks.Length() > 0 {
			prevBlock = stackedBlocks.Peek()
		}
		decPrevBlock := encryption.EncryptBlock(prevBlock, keySegments)
		decypheredText := bitops.MultiXOR64(currentBlock, decPrevBlock)

		flushStack.Append(decypheredText)
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
