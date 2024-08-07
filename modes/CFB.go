package modes

import (
	"bufio"
	"fmt"
	"lea/bitops"
	"lea/core"
	"lea/errors"
	"lea/fingerprint"
	"lea/schedule"
	"lea/stream"
	"lea/utils"
	"os"
)


func PerformCFB(filePath string, bKey []byte, bSeed []byte, encrypt bool, keySize int) {
	chunks := stream.BinaryChunkStream(filePath)
	
	kChunks := fingerprint.LoadSource(bKey)
	sChunks := fingerprint.LoadSource(bSeed)
	
        key := fingerprint.SelectPrint(kChunks, keySize)
        seed := fingerprint.SelectPrint(sChunks, keySize)

	rk := schedule.KeySchedule(keySize, key, seed)

	var blocks [4]uint32
	
	if encrypt {
		encryptCFB(filePath, blocks, rk, chunks, keySize)
	} else {
		decryptCFB(filePath, blocks, rk, chunks, keySize)
	}

}


func encryptCFB(filePath string, blocks [4]uint32, keySegments []uint32, chunks []uint32, size int) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Encrypting", filePath)	
	var encChunks []uint32
	// ask user for IV through stdin
	fmt.Println("Please provide an IV: ")
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]
	
	// conver the input to a slice of uint32
	var IV [4]uint32 = fingerprint.Fingerprint128(fingerprint.LoadSource([]byte(input)))

	var prev [4]uint32 = IV
	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		if (i+1)%4 == 0 {
			prev = [4]uint32(core.SelectEncrypt(prev, keySegments, size))
			prev = bitops.MultiXOR32(prev, blocks)
			encChunks = append(encChunks, prev[:]...)
		}
	}
	
	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)


}

func decryptCFB(filePath string, blocks [4]uint32, keySegments []uint32, chunks []uint32, size int) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Encrypting", filePath)	
	var encChunks []uint32
	// ask user for IV through stdin
	fmt.Println("Please provide an IV: ")
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]
	
	

	// conver the input to a slice of uint32
	var IV [4]uint32 = fingerprint.Fingerprint128(fingerprint.LoadSource([]byte(input)))

	
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
		decPrevBlock := [4]uint32(core.SelectEncrypt(prevBlock, keySegments, size))
		decypheredText := bitops.MultiXOR32(currentBlock, decPrevBlock)

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
